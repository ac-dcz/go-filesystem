package handler

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"go-fs/common/geeweb"
	"go-fs/internal/repo"
	"go-fs/internal/repo/user"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrUserExists    = errors.New("user is exists already")
	ErrUserNotExists = errors.New("user not exists")
	ErrPwdError      = errors.New("password invaild")
)

const _pwdSalt = "beffffeb"

func encodePwd(pwd string) string {
	pwd += _pwdSalt
	return hex.EncodeToString([]byte(pwd))
}

// 用户登录时生成token
func createToken(name, pwd string) string {
	return hex.EncodeToString([]byte(fmt.Sprintf("[%s@%s]-%s", name, pwd, _pwdSalt)))
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
	name, pwd := c.PostForm("username"), c.PostForm("password")
	pwd = encodePwd(pwd)
	info, err := repo.SelectUser(name)
	if err == sql.ErrNoRows {
		c.String(http.StatusOK, "[%s] %v", name, ErrUserNotExists.Error())
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if info.PassWD != pwd {
		c.String(http.StatusInternalServerError, ErrPwdError.Error())
		return
	}
	info.LastActiveTs = time.Now().Unix()
	if err := repo.UpdateUser(info); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	token := createToken(name, pwd)
	info_cookie := http.Cookie{Name: "userinfo", Value: fmt.Sprintf("%s-%s", name, pwd), Expires: time.Now().Add(time.Second * 60)}
	token_cookie := http.Cookie{Name: "token", Value: token, Expires: time.Now().Add(time.Second * 60)}
	http.SetCookie(c.W, &token_cookie)
	http.SetCookie(c.W, &info_cookie)
	c.String(http.StatusOK, "Successfully SignIn")
}

// SignOutHandle: post /user/signout 用户退出
func SignOutHandle(c *geeweb.Context) {
	c.String(http.StatusOK, "Out")
}

// ModifyPwdHandle: post /user/modify/pwd 修改密码
func ModifyPwdHandle(c *geeweb.Context) {
	name, pwd, new_pwd := c.PostForm("username"), c.PostForm("password"), c.PostForm("new_password")
	pwd = encodePwd(pwd)
	info, err := repo.SelectUser(name)
	if err == sql.ErrNoRows {
		c.String(http.StatusOK, "[%s] %v", name, ErrUserNotExists.Error())
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	if info.PassWD != pwd {
		c.String(http.StatusInternalServerError, ErrPwdError.Error())
		return
	}
	info.PassWD = encodePwd(new_pwd)
	if err := repo.UpdateUser(info); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.String(http.StatusOK, "Successfully Update")
}

// ModifyStatusHandle: post /user/modify/status
func ModifyStatusHandle(c *geeweb.Context) {
	name, pwd, status := c.PostForm("username"), c.PostForm("password"), c.PostForm("status")
	pwd = encodePwd(pwd)
	info, err := repo.SelectUser(name)
	if err == sql.ErrNoRows {
		c.String(http.StatusOK, "[%s] %v", name, ErrUserNotExists.Error())
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	if info.PassWD != pwd {
		c.String(http.StatusInternalServerError, ErrPwdError.Error())
		return
	}
	if num, err := strconv.Atoi(status); err == nil {
		info.Status = num
	} else {
		c.String(http.StatusBadRequest, "the value type of status is not int %v", err)
		return
	}
	if err := repo.UpdateUser(info); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.String(http.StatusOK, "Successfully Update")
}

func registryUserHandler(group *geeweb.RouterGroup) {
	group.POST("/signup", SignUpHandle)
	group.POST("/signin", SignInHandle)
	group.POST("/signout", SignOutHandle)
	group.POST("/modify/pwd", ModifyPwdHandle)
	group.POST("/modify/status", ModifyStatusHandle)
}
