package raft

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"log/slog"

	"github.com/bytedance/sonic"
)

const (
	SERVER_STOP_ERROR = "server stop"
)

type RPCResponse struct {
	Response interface{}
	Error    error
}

type RPC struct {
	Command  interface{}      //请求、命令
	RespChan chan RPCResponse //响应
}

func (r *RPC) Respond(resp interface{}, err error) {
	r.RespChan <- RPCResponse{resp, err}
}

// 负责通过网络发送request和response。可以基于TCP、UDP、TLS、http、grpc等来实现。
type Transporter interface {
	Start(port int, server *Server)
	AppendEntries(peer *Peer, req *AppendEntriesRequest) (*AppendEntriesResponse, error)
	RequestVote(peer *Peer, req *VoteRequest) (*VoteResponse, error)
}

type HttpTransporter struct {
	appendEntriesPath string
	votePath          string
	httpClient        http.Client
}

func NewHttpTransporter(prefix string, timeout time.Duration) *HttpTransporter {
	t := &HttpTransporter{
		appendEntriesPath: joinUrlPath(prefix, "/ae"), //简单处理可以用prefix+ "/ae"
		votePath:          joinUrlPath(prefix, "/vote"),
	}
	t.httpClient.Timeout = timeout //设置请求超时
	return t
}

// --------------------------------------
// 发送请求
// --------------------------------------
func (t *HttpTransporter) AppendEntries(peer *Peer, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	var resp AppendEntriesResponse
	err := post(t.httpClient, peer, t.appendEntriesPath, req, &resp)
	if err != nil {
		return nil, err
	} else {
		return &resp, nil
	}
}

func (t *HttpTransporter) RequestVote(peer *Peer, req *VoteRequest) (*VoteResponse, error) {
	var resp VoteResponse
	err := post(t.httpClient, peer, t.votePath, req, &resp)
	if err != nil {
		return nil, err
	} else {
		return &resp, nil
	}
}

type RequestInterface interface {
	*AppendEntriesRequest | *VoteRequest
}

type ResponseInterface interface {
	*AppendEntriesResponse | *VoteResponse
}

// 不支持泛型方法，只支持泛型函数
func post[E RequestInterface, S ResponseInterface](httpClient http.Client, peer *Peer, path string, req E, resp S) error {
	slog.Debug("send http request", "type", reflect.TypeOf(req).Elem().Name(), "to server", peer.Id)
	bs, err := sonic.Marshal(req)
	if err != nil {
		slog.Error("marchal request failed", "request", reflect.TypeOf(req).Elem().Name(), "error", err)
		return err
	}
	url := joinUrlPath(peer.ConnectionString, path)
	httpResp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		slog.Error("post request failed", "request", reflect.TypeOf(req).Elem().Name(), "error", err)
		return err
	}
	defer httpResp.Body.Close()
	bs, err = io.ReadAll(httpResp.Body)
	if err != nil {
		slog.Error("read response failed", "response", reflect.TypeOf(resp).Elem().Name(), "error", err)
		return err
	}
	respBody := string(bs)
	var methodName string
	funcName, _, _, ok := runtime.Caller(1)
	if ok {
		methodName = runtime.FuncForPC(funcName).Name()
	}
	if httpResp.StatusCode != http.StatusOK {
		if respBody != SERVER_STOP_ERROR+"\n" {
			slog.Error(methodName+" got abnormal code", "status", httpResp.Status, "msg", respBody)
		}
		return fmt.Errorf("%s got abnormal code status %s msg %s", methodName, httpResp.Status, respBody)
	} else {
		slog.Debug("receive http response", "type", reflect.TypeOf(resp).Elem().Name(), "from server", peer.Id)
		err := sonic.Unmarshal(bs, resp)
		if err != nil {
			slog.Error("unmarshal response failed", "response", reflect.TypeOf(resp).Elem().Name(), "error", err)
			return err
		}
		return nil
	}
}

// --------------------------------------
// 处理响应
// --------------------------------------

func handler[E RequestInterface](server *Server, req E) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if server.GetState() == Stopped {
			http.Error(w, SERVER_STOP_ERROR, http.StatusInternalServerError)
			return
		}

		slog.Debug("receive http request", "type", reflect.TypeOf(req).Elem().Name(), "server", server.Id)
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid json request", http.StatusBadRequest)
			return
		}

		err = sonic.Unmarshal(bs, req)
		if err != nil {
			http.Error(w, "Unmarshal json failed: "+string(bs)+". "+err.Error(), http.StatusInternalServerError)
			return
		}

		rpc := RPC{
			Command:  req,
			RespChan: make(chan RPCResponse),
		}

		if server.GetState() == Stopped { //再检查一次，避免写server.rpcCh时报错：send on closed channel
			http.Error(w, SERVER_STOP_ERROR, http.StatusInternalServerError)
			return
		}

		server.rpcCh <- rpc //把请求放进去

		rpcResp := <-rpc.RespChan //把结果取出来（要阻塞）
		resp, err := rpcResp.Response, rpcResp.Error
		if err != nil || resp == nil {
			http.Error(w, "server failed", http.StatusInternalServerError)
			return
		}

		bs, err = sonic.Marshal(resp)
		if err != nil {
			http.Error(w, "marshal json failed", http.StatusInternalServerError)
			return
		}
		slog.Debug("send http response", "type", reflect.TypeOf(resp).Elem().Name(), "server", server.Id)
		w.Write(bs)
	}
}

func (t *HttpTransporter) AppendEntriesHandler(server *Server) http.HandlerFunc {
	return handler(server, &AppendEntriesRequest{})
}

func (t *HttpTransporter) VoteHandler(server *Server) http.HandlerFunc {
	return handler(server, &VoteRequest{})
}

// 启动一个Http Server。注意会一直阻塞
func (t *HttpTransporter) Start(port int, server *Server) {
	mux := http.NewServeMux()
	mux.Handle(t.appendEntriesPath, t.AppendEntriesHandler(server))
	mux.Handle(t.votePath, t.VoteHandler(server))
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), mux); err != nil {
		panic(err)
	}
}
