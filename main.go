package main

import (
	"flag"

	"github.com/LollipopKit/custed-server/api"
	"github.com/LollipopKit/custed-server/config"
	"github.com/LollipopKit/custed-server/consts"
	"github.com/LollipopKit/custed-server/logger"
	_ "github.com/LollipopKit/custed-server/service"
	"github.com/LollipopKit/custed-server/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 命令行解析
	addr := flag.String("a", "0.0.0.0:1327", "部署的地址")
	debug := flag.Bool("d", false, "是否开启调试模式")
	username := flag.String("u", "", "用户名")
	enableSchedulePush := flag.Bool("push", false, "是否开启上课推送通知")
	flag.Parse()

	config.RunSchedulePush = *enableSchedulePush
	config.Debug = *debug

	// 如果传入了用户名，则输出cookie后退出，方便获取cookie
	if *username != "" {
		logger.I("n=%s; s=%s", utils.EncodeBase64(*username), api.GenerateCookieMd5(*username))
		return
	}

	// 启动web服务
	startWeb(*addr)
}

func startWeb(addr string) {
	e := echo.New()

	// 恢复进程，panic不会终止
	e.Use(middleware.Recover())

	// 设置 web log
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:  consts.LogFormat,
		Skipper: consts.StaticLogSkipper,
	}))

	// Routes
	e.GET("/", api.HomeView)
	e.HEAD("/", api.HeadHome)
	// App 配置
	e.GET("/config", api.GetConfig)
	e.GET("/config/android/newest", api.GetNewestApk)
	// Token
	e.POST("/token/ios", api.AddiOSToken)
	e.POST("/token/android", api.AddAndroidToken)
	e.GET("/token/detail", api.TokenDetail)
	e.GET("/token/stat", api.TokenDBSpeedTest)
	// 推送
	e.POST("/push/token/ios", api.PushMsgiOS)
	e.POST("/push/token/android", api.PushAndroidMsg)
	e.POST("/push/user", api.PushMsg)
	// 课表
	e.GET("/schedule", api.GetSchedule)
	e.POST("/schedule", api.UpdateSchedule)
	e.GET("/scheduleKBPro", api.GetKBPro)
	e.POST("/scheduleKBPro", api.UpdateKBPro)
	e.GET("/kbpro", api.GetKBPro)
	e.POST("/kbpro", api.UpdateKBPro)
	e.GET("/schedule/push/:value", api.SwitchSchedulePush)
	e.GET("/schedule/next/:id", api.GetNextLesson)
	// 成绩
	e.GET("/grade", api.GetGrade)
	e.POST("/grade", api.UpdateGrade)
	// 考试
	e.GET("/exam", api.GetExam)
	e.POST("/exam", api.UpdateExam)
	// Custed用户登录后端
	e.POST("/verify", api.VerifyUser)
	// 日志
	e.GET("/log/:name", api.LogView)
	// 天气
	e.GET("/weather", api.WeatherApi)
	// 学期开始时间
	e.GET("/semester", api.GetSemesterStart)
	e.GET("/weeks", api.GetWeeksOfSemester)

	// Static
	e.Static("/res", "res")

	// 后端服务 Start！
	e.HideBanner = true
	e.Logger.Fatal(e.Start(addr))
}
