package config

import "time"

var (
	// 调试模式
	Debug = false
	// 是否开启上课推送通知
	RunSchedulePush = false
	// 开学时间时间
	SemesterStartTime time.Time
	// 数据库链接
	DBUrl = "http://localhost:3777/"
)

// 计算开学到现在位于第几周
func CalculateWeeksOfSemester(now time.Time) int {
	return int(now.Sub(SemesterStartTime).Hours()/24/7) + 1
}
