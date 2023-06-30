package api

import (
	"github.com/LollipopKit/custed-server/consts"
	"github.com/LollipopKit/custed-server/utils"
	"github.com/labstack/echo/v4"
)

func HeadHome(c echo.Context) error {
	return c.NoContent(200)
}

func HomeView(c echo.Context) error {
	_, userName := accountVerify(c)
	if isSuperUser(userName) {
		return okWithData(c, "Welcome, SU.")
	}
	return okWithData(c, "Welcome Custed Backend, "+userName)
}

func LogView(c echo.Context) error {
	_, user := accountVerify(c)
	if !isSuperUser(user) {
		return noPermission(c)
	}

	name := c.Param("name")
	path := consts.LogDir + name

	if utils.IsExist(path) {
		return c.File(path)
	}
	return invalidParam(c, "name", "不存在此名称的log文件")
}
