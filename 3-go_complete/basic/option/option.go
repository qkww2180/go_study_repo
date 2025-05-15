package option

import (
	"fmt"
)

type User struct {
	Name string
	Age  int
	Tags map[string]string
}

func NewUser(opts ...UserOption) *User {
	user := new(User)
	for _, opt := range opts {
		opt(user)
	}
	return user
}

// 第一种方式：直接把Option定义为一个函数
type UserOption func(*User)

var tomOption = func(u *User) { u.Name = "tom" }
var jimOption = func(u *User) { u.Name = "jim" }

func WithName(name string) UserOption {
	return func(u *User) {
		u.Name = name
	}
}

func WithAge(age int) UserOption {
	return func(u *User) {
		u.Age = age
	}
}

func WithTag(k, v string) UserOption {
	return func(u *User) {
		if u.Tags == nil {
			u.Tags = make(map[string]string)
		}
		u.Tags[k] = v
	}
}

// 第二种方式：把Option定义一个接口。比第一种写法复杂了很多，但是接口里可以包含多个函数
type UserOption2 interface {
	Apply(*User)
}

type UserName struct {
	Name string
}

func (un *UserName) Apply(user *User) { //UserName实现了接口UserOption2
	user.Name = un.Name
}

func NewUserName(name string) *UserName {
	return &UserName{name}
}

type UserAge func(*User)

func (ua UserAge) Apply(user *User) { //UserAge实现了接口UserOption2
	ua(user)
}

func NewUser2(opts ...UserOption2) *User {
	user := new(User)
	for _, opt := range opts {
		opt.Apply(user)
	}
	return user
}

func constructUser() {
	user1 := NewUser(
		WithName("大乔乔"),
		WithTag("address", "双榆树"),
		WithTag("capital", "100万"),
	)
	fmt.Printf("%+v\n", *user1)
	ageOpt := func(u *User) {
		u.Age = 18
	}
	user2 := NewUser2(NewUserName("大乔乔"), UserAge(ageOpt)) // h_http.HandlerFunc()就是这种模式，把
	fmt.Printf("%+v\n", *user2)
}

// 第二种方式：把Option定义一个接口
type FilterOption interface {
	Judge(*User) bool
}

type NameEqual struct {
	Name string
}

func (ne *NameEqual) Judge(user *User) bool {
	if user.Name == ne.Name {
		return true
	} else {
		return false
	}
}

type AgeBetween struct {
	FromAge int
	ToAge   int
}

func (ab *AgeBetween) Judge(user *User) bool {
	if user.Age >= ab.FromAge && user.Age <= ab.ToAge {
		return true
	} else {
		return false
	}
}

func UserQuery(filters ...FilterOption) []*User {
	users := []*User{} //模拟从数据里检索出了一批User
	filterdUsers := make([]*User, 0, len(users))
L:
	for _, user := range users {
		for _, opt := range filters {
			if !opt.Judge(user) { //但几有一个条件不满足
				continue L //检查下一个User
			}
		}
		filterdUsers = append(filterdUsers, user) //所有条件都满足
	}
	return filterdUsers
}

func getUser() {
	UserQuery()
	UserQuery(&NameEqual{"大乔乔"})
	UserQuery(&AgeBetween{18, 28})
	UserQuery(&NameEqual{"大乔乔"}, &AgeBetween{18, 28})
}

/**
**********  GRPC  **********
func NewServer(opt ...ServerOption) *Server
type ServerOption interface {
	apply(*serverOptions)
}
func MaxRecvMsgSize(m int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.maxReceiveMessageSize = m
	})
}
type funcServerOption struct {
	f func(*serverOptions)
}
func (fdo *funcServerOption) apply(do *serverOptions) {
	fdo.f(do)
}
func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}


**********  ETCD  **********
Put(ctx context.Context, key, val string, opts ...OpOption) (*PutResponse, error)
type OpOption func(*Op)
// WithLease attaches a lease ID to a key in 'Put' request.
func WithLease(leaseID LeaseID) OpOption {
	return func(op *Op) { op.leaseID = leaseID }
}


**********  Jaeger  **********
StartSpan(operationName string, opts ...StartSpanOption) Span
type StartSpanOption interface {
	Apply(*StartSpanOptions)
}
type StartSpanOptions struct {
	References []SpanReference
	StartTime time.Time
	Tags map[string]interface{}
}
type Tag struct {
	Key   string
	Value interface{}
}
func (t Tag) Apply(o *StartSpanOptions) {
	if o.Tags == nil {
		o.Tags = make(map[string]interface{})
	}
	o.Tags[t.Key] = t.Value
}


**********  gorm  **********
Open(dialector Dialector, opts ...Option) (db *DB, err error)
type Option interface {
	Apply(*Config) error
	AfterInitialize(*DB) error
}
**/
