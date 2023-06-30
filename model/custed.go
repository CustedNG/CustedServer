package model

import (
	"strconv"
	"strings"
)

type CustedConfig struct {
	Notify []struct {
		Version int    `json:"version"`
		Content string `json:"content"`
	} `json:"notify"`
	JwURL  []string `json:"jw_url"`
	Update struct {
		Version struct {
			Android int `json:"android"`
			Ios     int `json:"ios"`
		} `json:"version"`
		Priority struct {
			Android int `json:"android"`
			Ios     int `json:"ios"`
		} `json:"priority"`
		URL struct {
			Android string `json:"android"`
			Ios     string `json:"ios"`
		} `json:"url"`
		Changelog struct {
			Android string `json:"android"`
			Ios     string `json:"ios"`
		} `json:"changelog"`
	} `json:"update"`
	NotShowRealUI []int    `json:"not_show_real_ui"`
	UseKbpro      []string `json:"use_kbpro"`
	TesterList    []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"tester_list"`
	Banner struct {
		ImgURL    string `json:"img_url"`
		ActionURL string `json:"action_url"`
	} `json:"banner"`
	HaveExam       bool `json:"have_exam"`
	SchoolCalendar []struct {
		Term       string `json:"term"`
		PicURL     string `json:"pic_url"`
		StrSummary string `json:"str_summary"`
	} `json:"school_calendar"`
	SemesterStart []int `json:"semester_start"`
}

type TokenItem struct {
	Tokens     []Token  `json:"tokens"`
	CreateTime string   `bson:"ctime"`
	LastTime   string   `bson:"ltime"`
	IPs        []string `bson:"ips"`
	Id         string   `bson:"id"`
}

type Token struct {
	ID string `bson:"id"`
	// Android 1, iOS 2, Windows 3, MacOS 4
	Platform int `bson:"platform"`
}

type User struct {
	Name       string `bson:"name"`
	Pwd        string `bson:"pwd"`
	CreateTime string `bson:"create_time"`
	Email      string `bson:"email"`
}

type NextLesson struct {
	Name      string
	Position  string
	Teacher   string
	StartTime string
	Weeks     WeeksList
}

type WeeksList []int

func (a WeeksList) Len() int      { return len(a) }
func (a WeeksList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a WeeksList) Less(i, j int) bool {
	return a[i] < a[j]
}

type LessonList []NextLesson

func (a LessonList) Len() int      { return len(a) }
func (a LessonList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a LessonList) Less(i, j int) bool {
	str1, _ := strconv.Atoi(strings.Replace(a[i].StartTime, ":", "", 1))
	str2, _ := strconv.Atoi(strings.Replace(a[j].StartTime, ":", "", 1))
	return str1 < str2
}

type HourMinute struct {
	Hour   int
	Minute int
}

type YearMonthDay struct {
	Year  int
	Month int
	Day   int
}

// 节假日跳过上课推送
// 如果今天是From的年月日，将会使用虚假数据（To属性的年月日）
// 如果To为nil，则不使用虚假数据
type HolidaySkip struct {
	From *YearMonthDay
	To   *YearMonthDay
}
