package handler

import (
	"fmt"
	"go-fs/common/geeweb"
	"net/http"
	"strings"
)

func VerifyToken(c *geeweb.Context) {
	if info, err := c.R.Cookie("userinfo"); err == nil {
		if token, err := c.R.Cookie("token"); err == nil {
			fmt.Println("HAHA")
			temp := strings.Split(info.Value, "-")
			if len(temp) == 2 {
				name, pwd := temp[0], temp[1]
				if createToken(name, pwd) == token.Value {
					c.Next()
					return
				}
			}
		}
	}
	c.String(http.StatusOK, "账号未登录")
}

// RegistryHandleFunc 注册路由
func RegistryHandleFunc(group *geeweb.RouterGroup) {
	user, file := group.Group("/user"), group.Group("/file")
	file.UseMiddleWare(VerifyToken)
	registryFileHandler(file)
	registryUserHandler(user)
}
