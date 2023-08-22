package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LollipopKit/custed-server/config"
	"github.com/LollipopKit/custed-server/db"
	"github.com/LollipopKit/custed-server/logger"
	"github.com/LollipopKit/custed-server/model"
	"github.com/LollipopKit/custed-server/utils"
	"github.com/labstack/echo/v4"
)

var (
	appConfig model.CustedConfig
)

func init() {
	go func() {
		for {
			data, err := os.ReadFile("res/config.json")
			if err != nil {
				logger.E("[api.init] read config.json failed: %s", err.Error())
				goto SLEEP
			}
			err = json.Unmarshal(data, &appConfig)
			if err != nil {
				logger.E("[api.init] unmarshal config.json failed: %s", err.Error())
				goto SLEEP
			}
			if len(appConfig.SemesterStart) == 3 {
				year := appConfig.SemesterStart[0]
				month := appConfig.SemesterStart[1]
				day := appConfig.SemesterStart[2]
				config.SemesterStartTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
			} else {
				logger.E("[api.init] config.json semester_start format err: %s", appConfig.SemesterStart)
			}
		SLEEP:
			time.Sleep(time.Minute * 5)
		}
	}()
}

func GetConfig(c echo.Context) error {
	return okWithData(c, appConfig)
}

func GetNewestApk(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, appConfig.Update.URL.Android)
}

// 更新成绩接口
func UpdateGrade(c echo.Context) error {
	_, user := accountVerify(c)
	if !isSuperUser(user) {
		logger.W("[api.UpdateGrade] %s has no permission", user)
		return noPermission(c)
	}

	grade := new(model.JwGrade)
	err := c.Bind(grade)
	if err != nil {
		logger.E("[api.UpdateGrade] bind err: %s", err.Error())
		return invalidBody(c, err.Error())
	}
	if grade.State != 0 {
		logger.W("[api.UpdateGrade] state code is %d", grade.State)
		return invalidBody(c, "state code is not 0")
	}

	err = db.UpdateGrade(user, *grade)
	if err != nil {
		logger.E("[api.UpdateGrade] update grade failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

	return ok(c)
}

// 获取成绩接口
func GetGrade(c echo.Context) error {
	_, user := accountVerify(c)
	if user == "CustedFakeUser" {
		goto JUMP
	}
	if !isSuperUser(user) {
		logger.W("[api.GetGrade] %s has no permission", user)
		return noPermission(c)
	}

JUMP:
	grade, err := db.GetGrade(user)
	if err != nil {
		logger.E("[api.GetGrade] get grade failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

	return okWithData(c, grade)
}

// 更新课表接口
func UpdateSchedule(c echo.Context) error {
	loggedIn, user := accountVerify(c)
	if !loggedIn {
		return noPermission(c)
	}

	scheduleModel := new(model.JwSchedule)
	err := c.Bind(scheduleModel)
	if err != nil {
		logger.E("[api.UpdateSchedule] bind err: %s", err.Error())
		return invalidBody(c, err.Error())
	}
	if scheduleModel.State != 0 {
		logger.W("[api.UpdateSchedule] state code is %d", scheduleModel.State)
		return invalidBody(c, "state code is not 0")
	}

	user = fmt.Sprintf("%s.%s", user, utils.GetTerm())
	err = db.UpdateSchedule(*scheduleModel, user)
	if err != nil {
		logger.E("[api.UpdateSchedule] update course schedule failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

	return ok(c)
}

// 获取课表接口
func GetSchedule(c echo.Context) error {
	loggedIn, user := accountVerify(c)
	if !loggedIn {
		return noPermission(c)
	}

	id := c.FormValue("id")
	if isSuperUser(user) {
		user = id
	}

	// 增加学期，eg：2019003373 => 2019003373-20222
	if !strings.Contains(user, ".") {
		user = fmt.Sprintf("%s.%s", user, utils.GetTerm())
	}

	schedule, err := db.GetSchedule(user)
	if err != nil {
		logger.E("[api.GetSchedule] get schedule failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

	return okWithData(c, schedule)
}

// 获取下一节课（桌面课表）接口
func GetNextLesson(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return invalidParam(c, "id", "为空")
	}
	id = fmt.Sprintf("%s.%s", id, utils.GetTerm())

	lesson, err := db.GetNextLesson(id, false)
	if err == db.ErrTodayNoMoreLesson {
		return okWithData(c, err.Error())
	}
	if err != nil {
		return dbOpFailed(c, err.Error())
	}

	return okWithData(c, lesson)
}

// 更新考试表接口
func UpdateExam(c echo.Context) error {
	loggedIn, user := accountVerify(c)
	if !loggedIn {
		return noPermission(c)
	}

	var exam model.JwExam
	err := c.Bind(&exam)
	if err != nil {
		logger.E("[api.UpdateExam] bind err: %s", err.Error())
		return invalidBody(c, err.Error())
	}

	if exam.State != 0 {
		logger.W("[api.UpdateExam] state code is %d", exam.State)
		return invalidBody(c, "state code is not 0")
	}

	user = fmt.Sprintf("%s.%s", user, utils.GetTerm())
	err = db.UpdateExam(user, exam)
	if err != nil {
		logger.E("[api.UpdateExam] update exam failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

	return ok(c)
}

// 获取考试表接口
func GetExam(c echo.Context) error {
	loggedIn, user := accountVerify(c)
	if !loggedIn {
		return noPermission(c)
	}

	if isSuperUser(user) {
		user = c.FormValue("id")
	}

	if !strings.Contains(user, ".") {
		user = fmt.Sprintf("%s.%s", user, utils.GetTerm())
	}

	exam, err := db.GetExam(user)
	if err != nil {
		logger.E("[api.GetExam] get exam failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

	return okWithData(c, exam)
}

// 切换开启课表推送接口
func SwitchSchedulePush(c echo.Context) error {
	loggedIn, id := accountVerify(c)
	if !loggedIn {
		return noPermission(c)
	}

	value := c.Param("value")
	on := value == "on"

	users, err := db.GetUsersEnableSchedulePush()
	if err != nil {
		logger.E("[api.SwitchSchedulePush] get users enable schedule push failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

	if on {
		for _, user := range users {
			if user == id {
				goto END
			}
		}
		users = append(users, id)
	} else {
		for idx, user := range users {
			if user == id {
				users = append(users[:idx], users[idx+1:]...)
				break
			}
		}
	}

	err = db.UpdateUsersEnableSchedulePush(users)
	if err != nil {
		logger.E("[api.SwitchSchedulePush] update users enable schedule push failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

END:
	return ok(c)
}

// 验证用户（等同于登陆到后端）
func VerifyUser(c echo.Context) error {
	cookie := c.FormValue("cookie")
	id := c.FormValue("id")
	url := c.FormValue("url")

	if id == "FakeUser" {
		setCookie(c, "CustedFakeUser", false, false)
		return ok(c)
	}
	if id == "null" {
		return invalidParam(c, "id", "为null")
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(
		"GET",
		url+"/welcome",
		nil,
	)
	if err != nil {
		if err != http.ErrUseLastResponse {
			logger.E("[api.VerifyUser] new request failed: %s", err.Error())
			return serverHTTPRequestFailed(c, "new request failed")
		}
	}

	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		logger.E("[api.VerifyUser] do request failed: %s", err.Error())
		return serverHTTPRequestFailed(c, "do request failed")
	}

	defer resp.Body.Close()

	if resp.StatusCode < 400 {
		setCookie(c, id, false, false)
		return ok(c)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.E("[api.VerifyUser] read body failed: %s", err.Error())
		return serverHTTPRequestFailed(c, "read body failed")
	}

	logger.W("[api.VerifyUser] %s failed:\n%s", id, string(body))
	return invalidParam(c, "cookie", "非法/无效")
}

// 更新kbpro接口
func UpdateKBPro(c echo.Context) error {
	loggedIn, user := accountVerify(c)
	if !loggedIn {
		return noPermission(c)
	}

	lessons := new(model.JwKBPro)
	err := c.Bind(lessons)
	if err != nil {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logger.E("[api.UpdateKBPro] read body failed: %s", err.Error())
			return internalError(c, err.Error())
		}
		logger.E("[api.UpdateKBPro] bind err: \n%s", string(body))
		return invalidJson(c, err.Error())
	}

	user = fmt.Sprintf("%s.%s", user, utils.GetTerm())
	err = db.UpdateKBPro(user, *lessons)
	if err != nil {
		logger.E("[api.UpdateKBPro] update kbpro failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

	return ok(c)
}

// 更新kbpro数据接口
func GetKBPro(c echo.Context) error {
	loggedIn, user := accountVerify(c)
	if !loggedIn {
		return noPermission(c)
	}

	if isSuperUser(user) {
		user = c.FormValue("id")
	}
	if !strings.Contains(user, ".") {
		user = fmt.Sprintf("%s.%s", user, utils.GetTerm())
	}

	schedule, err := db.GetKBPro(user)
	if err != nil {
		logger.E("[api.GetKBPro] get kbpro failed: %s", err.Error())
		return dbOpFailed(c, err.Error())
	}

	return okWithData(c, schedule)
}

func GetSemesterStart(c echo.Context) error {
	return okWithData(c, config.SemesterStartTime)
}

func GetWeeksOfSemester(c echo.Context) error {
	return okWithData(c, config.CalculateWeeksOfSemester(time.Now()))
}
