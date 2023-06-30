package push

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/LollipopKit/custed-server/consts"
	"github.com/LollipopKit/custed-server/db"
	"github.com/LollipopKit/custed-server/logger"
	"github.com/LollipopKit/custed-server/model"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

var (
	cert, errLoadPushCert = certificate.FromP12File(consts.APNSFile, consts.APNSPassword)
)

func init() {
	// 检查苹果推送证书加载是否存在问题
	if errLoadPushCert != nil {
		logger.E("[api.push] load push cert err: ", errLoadPushCert)
		panic("[api.push] loadCert err: " + errLoadPushCert.Error())
	}
}

// 发送所有平台消息接口
func SendAllPlatform(title, content string, production bool, users []string) error {
	tokenItems, err := db.GetTokenItems(false)

	if err != nil {
		logger.E("[api.SendAllPlatform] get token items err: ", err)
		return err
	}

	var androidTokenIDs []string
	var iosTokenIDs []string
	for _, user := range users {
		for _, tokenItem := range tokenItems {
			if tokenItem.Id != user {
				continue
			}

			for _, token := range tokenItem.Tokens {
				if token.Platform == 1 {
					androidTokenIDs = append(androidTokenIDs, token.ID)
				} else {
					iosTokenIDs = append(iosTokenIDs, token.ID)
				}
			}
		}
	}

	// 因为可能很慢，就不等待返回error
	// 错误会输出到日志中
	go SendiOS(title, content, production, iosTokenIDs)
	androidErr := SendAndroid(content, title, androidTokenIDs)

	if androidErr == nil {
		return nil
	}

	return fmt.Errorf("androir err: %w", androidErr)
}

// 发送iOS推送的具体实现
func SendiOS(title, content string, production bool, receivers []string) error {
	var client *apns2.Client
	// 环境配置，是否为生产环境（发布于App Store即为生产环境）
	if production {
		client = apns2.NewClient(cert).Production()
	} else {
		client = apns2.NewClient(cert).Development()
	}

	payload := payload.NewPayload().AlertTitle(title).AlertBody(content)

	if len(receivers) == 1 && receivers[0] == "all" {
		tokenItems, err := db.GetTokenItems(false)
		if err != nil {
			return err
		}

		// 清除 receivers内的 "all"
		receivers = []string{}

		for _, tokenItem := range tokenItems {
			for _, token := range tokenItem.Tokens {
				if token.Platform != 2 {
					continue
				}
				receivers = append(receivers, token.ID)
			}
		}
	}

	notifications := make(chan *apns2.Notification, len(receivers))
	responses := make(chan *model.PushResponseWithID, len(receivers))

	for _, token := range receivers {
		n := &apns2.Notification{
			DeviceToken: token,
			Topic:       consts.APNSPushTopic,
			Payload:     payload,
		}
		notifications <- n
	}

	workerCount := len(receivers) / 50
	if workerCount == 0 {
		workerCount = 1
	}

	for i := 0; i < workerCount; i++ {
		go iOSPushWorker(client, notifications, responses)
	}

	for i := 0; i < len(receivers); i++ {
		<-responses
		// res := <-responses
		// logger.I("%v %v %v\n", res.Resp.Sent(), res.Id, res.Resp.Reason)
	}
	return nil
}

// iOS推送发送的具体线程
func iOSPushWorker(client *apns2.Client, notis chan *apns2.Notification, resps chan *model.PushResponseWithID) {
	for noti := range notis {
		resp, err := client.Push(noti)
		if err != nil {
			logger.E("api.iOSPushWorker send err: ", err)
		}
		resps <- &model.PushResponseWithID{
			Resp: resp,
			Id:   noti.DeviceToken,
		}
	}
}

// 发送Android推送的具体实现
func SendAndroid(content, title string, ids []string) error {
	client := &http.Client{}
	var payload url.Values
	var err error

	if len(ids) == 1 && ids[0] == "all" {
		tokenItems, err := db.GetTokenItems(false)
		if err != nil {
			return err
		}

		ids = []string{}

		for _, tokenItem := range tokenItems {
			for _, token := range tokenItem.Tokens {
				if token.Platform != 1 {
					continue
				}
				ids = append(ids, token.ID)
			}
		}
	}

	payload = url.Values{
		"payload":                 {content},
		"title":                   {title},
		"restricted_package_name": {"cc.xuty.custed2"},
		"description":             {content},
		"pass_through":            {"0"},
		"notify_type":             {"-1"},
		"registration_id":         ids,
	}

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", consts.XiaomipushUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", consts.XiaomiPushAuthorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request push api failed: %s", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !bytes.Contains(body, []byte(`"result":"ok"`)) {
		return errors.New(string(body))
	}

	return nil
}
