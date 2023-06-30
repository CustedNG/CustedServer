package api

import (
	"fmt"
	"net/http"

	"github.com/LollipopKit/custed-server/consts"
	"github.com/labstack/echo/v4"
)

func ok(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"code":    consts.CodeSuccess,
		"message": "ok",
	})
}

func okWithData(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, map[string]any{
		"code":    consts.CodeSuccess,
		"message": "ok",
		"data":    data,
	})
}

func noPermission(c echo.Context) error {
	return c.JSON(http.StatusForbidden, map[string]any{
		"code":    consts.CodeNoPermission,
		"message": "没有权限",
	})
}

func invalidParam(c echo.Context, paramName, err string) error {
	return c.JSON(http.StatusBadRequest, map[string]any{
		"code":    consts.CodeInvalidParam,
		"message": fmt.Sprintf("参数[%s]错误：%s", paramName, err),
	})
}

func invalidBody(c echo.Context, err string) error {
	return c.JSON(http.StatusBadRequest, map[string]any{
		"code":    consts.CodeInvalidBody,
		"message": "HTTP Body错误：" + err,
	})
}

func dbOpFailed(c echo.Context, err string) error {
	return c.JSON(http.StatusInternalServerError, map[string]any{
		"code":    consts.CodeDBOperationFailed,
		"message": "数据库操作失败：" + err,
	})
}

func invalidJson(c echo.Context, err string) error {
	return c.JSON(http.StatusBadRequest, map[string]any{
		"code":    consts.CodeInvalidJson,
		"message": "JSON错误：" + err,
	})
}

func serverHTTPRequestFailed(c echo.Context, err string) error {
	return c.JSON(http.StatusInternalServerError, map[string]any{
		"code":    consts.CodeServerHTTPRequestFailed,
		"message": "服务器HTTP请求失败：" + err,
	})
}

func pushFailed(c echo.Context, err any) error {
	return c.JSON(http.StatusInternalServerError, map[string]any{
		"code":    consts.CodePushFailed,
		"message": err,
	})
}

func internalError(c echo.Context, err any) error {
	return c.JSON(http.StatusInternalServerError, map[string]any{
		"code":    consts.CodeInternal,
		"message": err,
	})
}
