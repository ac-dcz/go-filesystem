package handler

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"go-fs/common/geeweb"
	"go-fs/internal/repo"
	"go-fs/internal/repo/user"
	"net/http"
)

var (
	ErrUserExists = errors.New("user is exists already")
)

const _pwdSalt = "beffffeb"

func encodePwd(pwd string) string {
	pwd += _pwdSalt
	return hex.EncodeToString([]byte(pwd))
}

func createUser(name, pwd, email, tel string) (*user.UserInfo, error) {
	pwd = encodePwd(pwd)
	info := user.NewInfo(name, pwd, email, tel)
	_, err := repo.SelectUser(name)
	if err == sql.ErrNoRows {
		return info, nil
	} else if err != nil {
		return nil, err
	}
	return nil, ErrUserExists
}

// SignUpHandle: post /user/signup 用户注册
func SignUpHandle(c *geeweb.Context) {
	name, pwd := c.PostForm("username"), c.PostForm("password")
	email := c.PostForm("email")
	tel := c.PostForm("tel")
	user, err := createUser(name, pwd, email, tel)
	if err == ErrUserExists {
		c.String(http.StatusOK, err.Error())
	} else if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		if err := repo.InsertUser(user); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "Successfully SignUp!!!")
		}
	}
}

// SignInHandle: post /user/signin 用户登录
func SignInHandle(c *geeweb.Context) {

}

// SignOutHandle: post /user/signout 用户退出
func SignOutHandle(c *geeweb.Context) {

}

// ModifyPwdHandle: post /user/modify/pwd 修改密码
func ModifyPwdHandle(c *geeweb.Context) {

}

// ModifyStatusHandle: post /user/modify/status
func ModifyStatusHandle(c *geeweb.Context) {

}

func registryUserHandler(group *geeweb.RouterGroup) {
	group.POST("/user/signup", SignUpHandle)
	group.POST("/user/signin", SignInHandle)
	group.POST("/user/signout", SignOutHandle)
	group.POST("/user/modify/pwd", ModifyPwdHandle)
	group.POST("/user/modify/status", ModifyStatusHandle)
}
