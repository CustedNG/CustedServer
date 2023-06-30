package consts

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// web log 格式
const LogFormat = "${time_rfc3339} ${method} ${status} ${uri} \nLatency: ${latency_human}  ${error}\n"

// log跳过static文件夹
func StaticLogSkipper(context echo.Context) bool {
	url := context.Request().URL.Path
	if strings.Contains(url, "/static") || url == "/favicon.ico" || strings.HasPrefix(url, "/res") {
		return true
	}
	return false
}

// cookie的键
const CookieSignKey = "s"
const CookieNameKey = "n"

// 默认并发数
const CoroutinesNum = 10

// 文件权限
const Permission = 0755

// 路径配置
const LogDir = ".log/"

const UsersEnableSchedulePushKey = "users_enable_schedule_push"
const SuperUserListKey = "sus"
const TimeOfSemesterStartKey = "semester_start"

// 推送相关
const APNSPushTopic = "com.tusi.app"
const XiaomipushUrl = "https://api.xmpush.xiaomi.com/v3/message/regid"

const (
	WeatherAPIUrl = "https://cust.app/app/weather"
)

// 最多向后搜索两次，跳过周末
const NextScheduleSearchPassDays = 3
