package api

import (
	"github.com/LollipopKit/custed-server/push"
	"github.com/labstack/echo/v4"
)

// 推送消息接口
func PushMsg(c echo.Context) error {
	_, userName := accountVerify(c)
	if !isSuperUser(userName) {
		return noPermission(c)
	}

	user := c.FormValue("user")
	title := c.FormValue("title")
	if title == "" {
		return invalidParam(c, "title", "为空")
	}
	content := c.FormValue("content")
	if content == "" {
		return invalidParam(c, "content", "为空")
	}
	production := c.FormValue("env") == "production"

	var users []string
	err := json.UnmarshalFromString(user, &users)
	if err != nil {
		return invalidParam(c, "user", "非法json list")
	}

	err = push.SendAllPlatform(title, content, production, users)
	if err != nil {
		return pushFailed(c, err.Error())
	}

	return ok(c)
}

// 推送消息接口（iOS）
func PushMsgiOS(c echo.Context) error {
	_, userName := accountVerify(c)
	if !isSuperUser(userName) {
		return noPermission(c)
	}

	token := c.FormValue("token")
	title := c.FormValue("title")
	if title == "" {
		return invalidParam(c, "title", "为空")
	}
	content := c.FormValue("content")
	if content == "" {
		return invalidParam(c, "content", "为空")
	}
	production := c.FormValue("env") == "production"

	var tokens []string
	err := json.UnmarshalFromString(token, &tokens)
	if err != nil {
		return invalidParam(c, "token", "非法json list")
	}

	err = push.SendiOS(title, content, production, tokens)
	if err != nil {
		return pushFailed(c, err.Error())
	}

	return ok(c)
}

// 推送消息接口（Android）
func PushAndroidMsg(c echo.Context) error {
	_, userName := accountVerify(c)
	if !isSuperUser(userName) {
		return noPermission(c)
	}

	id := c.FormValue("token")
	if id == "" {
		return invalidParam(c, "token", "为空")
	}

	var ids []string
	err := json.UnmarshalFromString(id, &ids)
	if err != nil {
		return invalidParam(c, "token", "非法json list")
	}

	content := c.FormValue("content")
	title := c.FormValue("title")

	if content == "" {
		return invalidParam(c, "content", "为空")
	}
	if title == "" {
		return invalidParam(c, "title", "为空")
	}

	err = push.SendAndroid(content, title, ids)
	if err != nil {
		var msg any
		err = json.UnmarshalFromString(err.Error(), &msg)
		if err != nil {
			return invalidJson(c, err.Error())
		}
		return pushFailed(c, msg)
	}

	return ok(c)
}
