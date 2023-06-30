package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/LollipopKit/custed-server/config"
	"github.com/LollipopKit/custed-server/db"
	"github.com/LollipopKit/custed-server/logger"
	"github.com/LollipopKit/custed-server/model"
	"github.com/LollipopKit/custed-server/push"
	"github.com/LollipopKit/custed-server/utils"
)

var (
	timeShouldCheck = []model.HourMinute{
		{7, 45},
		{9, 45},
		{13, 15},
		{15, 15},
		{17, 45},
		{19, 20},
	}
	checkInterval = 5
)

func init() {
	go func() {
		time.Sleep(time.Second * 3)
		if config.RunSchedulePush {
			go sendSchedulePush()
		}
	}()
}

func sendSchedulePush() {
	ticker := time.NewTicker(time.Minute * time.Duration(checkInterval))
	logger.I("上课推送已开启")

	for range ticker.C {
		now := time.Now()

		for _, ti := range timeShouldCheck {
			diffMinute := ti.Minute - now.Minute()
			if now.Hour() == ti.Hour && diffMinute < checkInterval+1 && diffMinute > 0 {

				logger.I("午时已到，开始检测推送")

				users, err := db.GetUsersEnableSchedulePush()
				if err != nil {
					logger.E("service.push GetUsersEnableSchedulePush err: %s", err)
					// https://blog.lolli.tech/go-cheat-sheet/#ticker
					// ticker中continue不会直接进入下一个循环，而是会等待ticker到达预定时间才运行
					continue
				}

				tokenItems, err := db.GetTokenItems(false)
				if err != nil {
					logger.E("service.push GetTokenItems err: %s", err)
					continue
				}

				for _, user := range users {
					lesson, err := db.GetNextLesson(user, false)
					if err != nil {
						continue
					}

					startTime := string(lesson.StartTime)
					// 包含空格代表非今天的课程，如“明天 08:00”
					notToday := strings.Contains(startTime, " ")
					if notToday {
						continue
					}

					startHour, err := utils.Str2int(strings.Split(startTime, ":")[0])

					if err != nil {
						logger.E("service.push Str2int err: %s", err)
						continue
					}

					if startHour-ti.Hour > 1 {
						continue
					}

					for _, tokenItem := range tokenItems {
						if tokenItem.Id == user {
							for _, token := range tokenItem.Tokens {
								switch token.Platform {
								// case 1:
								// 	err = sendAndroidPush(
								// 		fmt.Sprintf("%s %s %s", lesson.Position, lesson.StartTime, lesson.Teacher),
								// 		lesson.Name,
								// 		[]string{token.ID},
								// 	)
								case 2:
									err = push.SendiOS(
										lesson.Name,
										fmt.Sprintf("%s %s %s", lesson.StartTime, lesson.Position, lesson.Teacher),
										true,
										[]string{token.ID},
									)
								}

								if err != nil {
									logger.E("service.push SendPush err: %s", err)
								}
							}
						}
					}
				}
			}
		}
	}
}
