package user

import "time"

type UserInfo struct {
	Name         string
	PassWD       string
	Email        string
	Tel          string
	EmailChecked int
	TelChecked   int
	SignUpTs     int64
	LastActiveTs int64
	Profile      string
	Status       int //用户状态0/1/2/3(启用/禁用/标记/删除)
}

func NewInfo(name, pwd, email, tel string) *UserInfo {
	return &UserInfo{
		Name:         name,
		PassWD:       pwd,
		Email:        email,
		Tel:          tel,
		EmailChecked: 0,
		TelChecked:   0,
		SignUpTs:     time.Now().Unix(),
		LastActiveTs: time.Now().Unix(),
		Profile:      "",
		Status:       0,
	}
}
