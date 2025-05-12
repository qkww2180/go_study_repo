package main

/*
Option模式没有固定的应用场景，通常是作为不定长参数(所以才叫option)。
通过type把Option定义为函数或接口。
*/

type User struct {
	Name string
	Age  int
}

func NewUser(options ...UserOption) *User {
	user := new(User)
	for _, opt := range options {
		opt(user)
	}
	return user
}

type UserOption func(*User)

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

func main() {
	user := NewUser(WithAge(18), WithName("大乔乔"))
	_ = user

	users := []*User{user}
	QueryUser(users, Where{Name: "张三"}, Limit{Count: 3})
}

type QueryOption interface {
	Apply([]*User) []*User
}

type Where struct {
	Name    string
	FromAge int
	ToAge   int
}

func (w Where) Apply(users []*User) []*User {
	rect := make([]*User, 0, len(users))
	for _, user := range users {
		if user.Name == w.Name && user.Age >= w.FromAge && user.Age <= w.ToAge {
			rect = append(rect, user)
		}
	}
	return rect
}

type Limit struct {
	Offset int
	Count  int
}

func (l Limit) Apply(users []*User) []*User {
	if l.Offset >= len(users) {
		return nil
	}
	if l.Offset+l.Count >= len(users) {
		return users[l.Offset:]
	}
	return users[l.Offset : l.Offset+l.Count]
}

func QueryUser(users []*User, options ...QueryOption) []*User {
	rect := users
	for _, opt := range options {
		rect = opt.Apply(rect)
	}
	return rect
}
