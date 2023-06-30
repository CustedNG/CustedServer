package api

import (
	"fmt"
	"time"

	"github.com/LollipopKit/custed-server/db"
	"github.com/labstack/echo/v4"
)

// 接受iOS token push请求接口
func AddiOSToken(c echo.Context) error {
	loggedIn, userId := accountVerify(c)
	if !loggedIn {
		return noPermission(c)
	}

	token := c.FormValue("token")
	if len(token) < 64 {
		return invalidParam(c, "token", "长度不足")
	}

	err := db.UpdateToken(token, c.RealIP(), userId, 2)
	if err == nil {
		return ok(c)
	}
	return dbOpFailed(c, err.Error())
}

// 接受Android token push请求接口
func AddAndroidToken(c echo.Context) error {
	loggedIn, userId := accountVerify(c)
	if !loggedIn {
		return noPermission(c)
	}

	token := c.FormValue("token")
	err := db.UpdateToken(token, c.RealIP(), userId, 1)
	if err == nil {
		return ok(c)
	}
	return dbOpFailed(c, err.Error())
}

// 获取token详细信息接口
func TokenDetail(c echo.Context) error {
	id := c.QueryParam("id")
	if len(id) < 10 {
		return invalidParam(c, "id", "长度不足")
	}
	logged, user := accountVerify(c)
	if !logged || !isSuperUser(user) {
		return noPermission(c)
	}
	token, err := db.GetToken(id, true)
	if err != nil {
		return dbOpFailed(c, err.Error())
	}
	return okWithData(c, *token)
}

// 获取数据库速度测试接口
func TokenDBSpeedTest(c echo.Context) error {
	_, userName := accountVerify(c)
	if !isSuperUser(userName) {
		return noPermission(c)
	}

	t1 := time.Now()
	items, err := db.GetTokenItems(true)
	t2 := time.Now()
	if err != nil {
		return dbOpFailed(c, err.Error())
	}

	resultFmt := "Start: %s\nEnd: %s\nSpent(UnixNano): %d\nTotal tokens: %d items"
	return okWithData(c, fmt.Sprintf(resultFmt, t1, t2, int(t2.UnixNano()-t1.UnixNano()), len(items)))
}
