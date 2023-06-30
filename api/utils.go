package api

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/LollipopKit/custed-server/consts"
	"github.com/LollipopKit/custed-server/db"
	"github.com/LollipopKit/custed-server/logger"
	"github.com/LollipopKit/custed-server/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

var (
	suCache = *new(SuCache)
	json    = jsoniter.ConfigCompatibleWithStandardLibrary
)

// SuperUsers缓存
type SuCache struct {
	sus  []string
	time time.Time
}

// 校验是否为su
func isSuperUser(userName string) bool {
	if suCache.time.Before(time.Now().Add(-time.Minute)) {
		sus, err := db.GetSuperUserNames()
		if err == nil {
			suCache.sus = sus
			suCache.time = time.Now()
		} else {
			logger.E("api.isSuperUser err: %s", err)
		}
	}

	return utils.Contains(suCache.sus, userName)
}

func GenerateCookieMd5(name string) string {
	cookie := utils.EncodeBase64(consts.CookieSalt1 + name + consts.CookieSalt2)
	md5Hex := md5.Sum([]byte(cookie))
	cookie = hex.EncodeToString(md5Hex[:])
	cookie = reverseString(cookie)
	return cookie[11:13] + cookie + cookie[4:7]
}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func generateCookie(key, value string, remove, setExpire bool) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	if setExpire {
		cookie.Expires = time.Now().Add(time.Hour * 24 * 365)
	}
	cookie.Path = "/"
	if remove {
		cookie.MaxAge = -1
	}
	return cookie
}

// 校验账户用户名
func accountVerify(c echo.Context) (bool, string) {
	ip := c.RealIP()
	cookieSign, errSign := c.Cookie(consts.CookieSignKey)
	cookieName, errName := c.Cookie(consts.CookieNameKey)
	if errSign != nil || errName != nil {
		return false, ip
	}
	userName, err := utils.DecodeBase64(cookieName.Value)
	if err != nil {
		return false, ip
	}
	if cookieSign.Value == GenerateCookieMd5(userName) {
		return true, userName
	}
	return false, ip
}

func setCookie(c echo.Context, name string, remove, setExpire bool) {
	c.SetCookie(generateCookie(consts.CookieSignKey, GenerateCookieMd5(name), remove, setExpire))
	c.SetCookie(generateCookie(consts.CookieNameKey, utils.EncodeBase64(name), remove, setExpire))
}
