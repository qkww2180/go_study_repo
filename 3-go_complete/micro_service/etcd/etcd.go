package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var (
	wg = sync.WaitGroup{}
)

func getPut(ctx context.Context, client *etcdv3.Client) {
	if _, err := client.Put(ctx, "1", "张三"); err != nil { //写入一对<key,value>
		log.Printf("写入key value失败：%v", err)
		return
	}
	if _, err := client.Put(ctx, "2", "李四"); err != nil { //写入一对<key,value>
		log.Printf("写入key value失败：%v", err)
		return
	}
	if _, err := client.Put(ctx, "3", "王五"); err != nil { //写入一对<key,value>
		log.Printf("写入key value失败：%v", err)
		return
	}

	if response, err := client.Get(ctx, "1", etcdv3.WithRange("3")); err != nil { //etcdv3.WithRange(endKey string)是前闭后开区间，即不包含endKey。范围查询还有WithFromKey, WithRev, WithSort, WithPrefix等
		log.Printf("获取key对应的value失败：%v", err)
		return
	} else {
		if len(response.Kvs) == 0 {
			log.Println("key不存在")
		} else {
			for _, kv := range response.Kvs {
				log.Printf("%s = %s\n", kv.Key, kv.Value)
			}
		}
	}

	if response, err := client.Get(ctx, "1"); err != nil {
		log.Printf("获取key对应的value失败：%v", err)
		return
	} else {
		if len(response.Kvs) == 0 {
			log.Println("key不存在")
		} else {
			for _, kv := range response.Kvs {
				log.Printf("%s = %s\n", kv.Key, kv.Value)
			}
		}
	}
}

func keepAlive(ctx context.Context, client *etcdv3.Client) {
	const TTL = 2
	if lease, err := client.Grant(ctx, TTL); err != nil { //第二个参数ttl的单位是秒
		log.Printf("创建租约失败：%v", err)
		return
	} else {
		if _, err = client.Put(ctx, "name", "张三", etcdv3.WithLease(lease.ID)); err != nil { //写入一对<key,value>，带上过期时间
			log.Printf("写入key value失败：%v", err)
			return
		} else {
			time.Sleep(TTL / 2 * time.Second)
			if response, err := client.Get(ctx, "name"); err != nil {
				log.Printf("获取key对应的value失败：%v", err)
				return
			} else {
				if len(response.Kvs) == 0 {
					log.Println("key不存在")
				} else {
					for _, kv := range response.Kvs {
						log.Printf("%s = %s\n", kv.Key, kv.Value)
					}
				}
			}
			//本来租约有效期只有2秒，现在通过KeepAlive使租约永久生效
			if _, err = client.KeepAlive(ctx, lease.ID); err != nil { //必须在到期之前续，否则找不到对应的租约
				log.Printf("keep租约alive失败：%v", err)
				return
			} else {
				time.Sleep(2 * TTL * time.Second)
				if response, err := client.Get(ctx, "name"); err != nil {
					log.Printf("获取key对应的value失败：%v", err)
					return
				} else {
					if len(response.Kvs) == 0 {
						log.Println("key不存在")
					} else {
						for _, kv := range response.Kvs {
							log.Printf("%s = %s\n", kv.Key, kv.Value)
						}
					}
				}
			}
		}
	}
}

func keepAliveOnce(ctx context.Context, client *etcdv3.Client) {
	const TTL = 2
	if lease, err := client.Grant(ctx, TTL); err != nil { //第二个参数ttl的单位是秒
		log.Printf("创建租约失败：%v", err)
		return
	} else {
		if _, err = client.Put(ctx, "name", "张三", etcdv3.WithLease(lease.ID)); err != nil {
			log.Printf("写入key value失败：%v", err)
			return
		} else {
			if response, err := client.Get(ctx, "name"); err != nil {
				log.Printf("获取key对应的value失败：%v", err)
				return
			} else {
				if len(response.Kvs) == 0 {
					log.Println("key不存在")
				} else {
					for _, kv := range response.Kvs {
						log.Printf("%s = %s\n", kv.Key, kv.Value)
					}
				}
			}
			//本来租约有效期只有2秒，现在通过KeepAliveOnce再续2秒(之前剩余的寿命会清0)
			if _, err = client.KeepAliveOnce(ctx, lease.ID); err != nil { //必须在到期之前续，否则找不到对应的租约
				log.Printf("keep租约alive失败：%v", err)
				return
			} else {
				time.Sleep(TTL/2*time.Second + 200*time.Millisecond)
				if response, err := client.Get(ctx, "name"); err != nil {
					log.Printf("获取key对应的value失败：%v", err)
					return
				} else {
					if len(response.Kvs) == 0 {
						log.Println("key不存在")
					} else {
						for _, kv := range response.Kvs {
							log.Printf("%s = %s\n", kv.Key, kv.Value)
						}
					}
				}
				time.Sleep(TTL/2*time.Second + 200*time.Millisecond) //并不是绝对精确地TTL之后过期，所以这里稍微地拖延一下
				if response, err := client.Get(ctx, "name"); err != nil {
					log.Printf("获取key对应的value失败：%v", err)
					return
				} else {
					if len(response.Kvs) == 0 {
						log.Println("key不存在")
					} else {
						for _, kv := range response.Kvs {
							log.Printf("%s = %s\n", kv.Key, kv.Value)
						}
					}
				}
			}
		}
	}
}

// 尝试加锁，获取不到时会立即返回。返回的bool为false说明没有获得到锁
func tryLock(ctx context.Context, client *etcdv3.Client, lockName string) (*concurrency.Session, *concurrency.Mutex, bool) {
	if session, err := concurrency.NewSession(client); err != nil {
		log.Printf("协程%d 创建会话失败:%v\n", ctx.Value("id").(int), err)
		return nil, nil, false
	} else {
		mutex := concurrency.NewMutex(session, lockName)
		if err := mutex.TryLock(ctx); err != nil { //锁被其他协程(或进程)持有时立即返回concurrency.ErrLocked
			if err != concurrency.ErrLocked { //如果是没获得锁，不打印error信息
				log.Printf("协程%d TryLock异常:%#v\n", ctx.Value("id").(int), err)
			}
			return session, mutex, false
		} else {
			return session, mutex, true
		}
	}
}

func testTryLock(ctx context.Context, client *etcdv3.Client, lockName string, routineID int) {
	defer wg.Done()
	//go官方认为协程id不应该暴露给应用层，官方建议通过context关联上下文
	ctx = context.WithValue(ctx, "id", routineID)
	session, mutex, success := tryLock(ctx, client, lockName)
	if success == true {
		log.Printf("协程%d 获得锁\n", routineID)
		time.Sleep(time.Second) //执行业务需求
		mutex.Unlock(ctx)       //释放锁
		log.Printf("协程%d 释放锁\n", routineID)
		session.Close() //关闭会话
	} else {
		log.Printf("协程%d 锁被其他会话持有\n", routineID)
	}
}

// 尝试加锁，获取不到时会一直阻塞，直到获得锁
func lockWithoutTimeout(ctx context.Context, client *etcdv3.Client, lockName string) (*concurrency.Session, *concurrency.Mutex) {
	if session, err := concurrency.NewSession(client); err != nil {
		log.Printf("协程%d 创建会话失败:%v\n", ctx.Value("id").(int), err)
		return nil, nil
	} else {
		mutex := concurrency.NewMutex(session, lockName)
		if err := mutex.Lock(ctx); err != nil { //获得不到锁时，就一直阻塞
			log.Printf("协程%d Lock异常:%#v\n", ctx.Value("id").(int), err)
			return session, nil
		} else {
			return session, mutex
		}
	}
}

func testLockWithoutTimeout(ctx context.Context, client *etcdv3.Client, lockName string, routineID int) {
	defer wg.Done()
	//go官方认为协程id不应该暴露给应用层，官方建议通过context关联上下文
	ctx = context.WithValue(ctx, "id", routineID)
	session, mutex := lockWithoutTimeout(ctx, client, lockName)
	if session != nil {
		defer session.Close()
		if mutex != nil {
			log.Printf("协程%d 获得锁\n", routineID)
			time.Sleep(time.Second)
			mutex.Unlock(ctx) //释放锁
			log.Printf("协程%d 释放锁\n", routineID)
		}
	}
}

// 尝试加锁，获取不到时会一直阻塞，直到超时。返回的Mutex为nil说明没有获得到锁
func lockWithTimeout(ctx context.Context, client *etcdv3.Client, lockName string) (*concurrency.Session, *concurrency.Mutex) {
	toctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond) //设定超时100ms
	defer cancel()
	if session, err := concurrency.NewSession(client); err != nil {
		log.Printf("协程%d 创建会话失败:%v\n", ctx.Value("id").(int), err)
		return nil, nil
	} else {
		mutex := concurrency.NewMutex(session, lockName)
		if err := mutex.Lock(toctx); err != nil { //toctx是一个带超时的context。如果一直获得不到锁，超时后就返回error
			if err != context.DeadlineExceeded { //如果是超时，不打印erro信息
				log.Printf("协程%d Lock异常:%#v\n", ctx.Value("id").(int), err)
			} else {
				log.Printf("协程%d 指定时间内未获得到锁，放弃\n", ctx.Value("id").(int))
			}
			return session, nil
		} else {
			return session, mutex
		}
	}
}

func testLockWithTimeout(ctx context.Context, client *etcdv3.Client, lockName string, routineID int) {
	defer wg.Done()
	//go官方认为协程id不应该暴露给应用层，官方建议通过context关联上下文
	ctx = context.WithValue(ctx, "id", routineID)
	session, mutex := lockWithTimeout(ctx, client, lockName)
	if session != nil {
		defer session.Close()
		if mutex != nil {
			log.Printf("协程%d 获得锁\n", routineID)
			time.Sleep(time.Second)
			mutex.Unlock(ctx) //释放锁
			log.Printf("协程%d 释放锁\n", routineID)
		}
	}
}

func main1() {
	ctx := context.Background()
	if client, err := etcdv3.New(
		etcdv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 3 * time.Second,
			//通过etcdctl命令可创建多个账号(Username和Password)，一个账号可以拥有某种角色，一个角色可以对某些etcd目录持有读/写权限。详情参见https://juejin.cn/post/6844903678269210632
			Username: "",
			Password: "",
			//etcd证书分2种：client请求证书；etcd集群模式下peer节点之间通信证书。
			//etcd的client证书和我们浏览网站时https证书有一点不一样，浏览HTTPS网站时证书是由网站服务器提供，浏览器校验证书合法性。而etcd client端访问etcd server，则是client端请求时携带证书，由etcd server校验client证书的合法性。
			//详情参见https://zhuanlan.zhihu.com/p/596605552
			TLS: nil,
		},
	); err != nil {
		log.Printf("连接不上etcd服务器: %v", err)
		return
	} else {
		defer client.Close()

		// getPut(ctx, client)
		// log.Println(strings.Repeat("-", 50))

		// keepAlive(ctx, client)
		// log.Println(strings.Repeat("-", 50))

		// keepAliveOnce(ctx, client)
		// log.Println(strings.Repeat("-", 50))

		const LOCK_NAME = "LN"
		wg.Add(3)
		go testTryLock(ctx, client, LOCK_NAME, 1)
		go testTryLock(ctx, client, LOCK_NAME, 2)
		go testTryLock(ctx, client, LOCK_NAME, 3)
		wg.Wait()
		log.Println(strings.Repeat("-", 50))

		wg.Add(3)
		go testLockWithoutTimeout(ctx, client, LOCK_NAME, 1)
		go testLockWithoutTimeout(ctx, client, LOCK_NAME, 2)
		go testLockWithoutTimeout(ctx, client, LOCK_NAME, 3)
		wg.Wait()
		log.Println(strings.Repeat("-", 50))

		wg.Add(3)
		go testLockWithTimeout(ctx, client, LOCK_NAME, 1)
		go testLockWithTimeout(ctx, client, LOCK_NAME, 2)
		go testLockWithTimeout(ctx, client, LOCK_NAME, 3)
		wg.Wait()
		log.Println(strings.Repeat("-", 50))
	}
}

// go run .\micro_service\etcd\
