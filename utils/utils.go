package utils

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/LollipopKit/custed-server/consts"
	"github.com/LollipopKit/custed-server/logger"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func EncodeBase64(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(name))
}

func DecodeBase64(b64 string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		logger.W("api.decodeBase64 err: %s", err)
	}
	return string(b), err
}

func ReadFile(name string) []byte {
	f, err := ioutil.ReadFile(name)
	if err != nil {
		logger.E("utils.ReadFile err: %s", err)
	}
	return f
}

func WriteFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, consts.Permission)
}

// 如果不存在则创建
// [mode]: 0 -> mkdir, 1 -> new file
func DoIfNotExists(path, data string, mode int) {
	if !IsExist(path) {
		if mode == 0 {
			err := os.Mkdir(path, consts.Permission)
			if err != nil {
				logger.E("utils.DoIfNotExists err: %s", err)
			}
		} else {
			_, err := os.Create(path)
			if err != nil {
				logger.E("utils.DoIfNotExists err: %s", err)
			}
			err = WriteFile(path, []byte(data))
			if err != nil {
				logger.E("utils.DoIfNotExists err: %s", err)
			}
		}
	}
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsFileOutdate(path string, hours int) bool {
	info, err := os.Stat(path)
	if err != nil {
		logger.E("utils.IsFileOutdate err: %s", err)
		return false
	}
	return info.ModTime().Add(time.Hour * time.Duration(hours)).Before(time.Now())
}

func Str2int(str string) (int, error) {
	i, err := strconv.Atoi(str)
	return i, err
}

func Int2str(i int) string {
	return strconv.Itoa(i)
}

func GetTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 根据课表节数，转换为时间
func GetStartTimeBySection(section string) string {
	switch section {
	case "01", "1":
		return "8:00"
	case "03", "3":
		return "10:05"
	case "05", "5":
		return "13:30"
	case "07", "7":
		return "15:35"
	case "09", "9":
		return "18:00"
	case "11":
		return "19:45"
	default:
		return "00:00"
	}
}

func Contains[T int | float64 | int64 | string](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func GetTerm() string {
	now := time.Now()
	year := Int2str(now.Year())
	month := now.Month()
	if month >= 8 {
		return year + "2"
	}
	return year + "1"
}
