package model

import "github.com/sideshow/apns2"

type PushResponseWithID struct {
	Resp *apns2.Response
	Id   string
}

// type JPushSend2AllPayload struct {
// 	Platform     string       `json:"platform"`
// 	Audience     string       `json:"audience"`
// 	Notification struct {
// 		Android struct {
// 			Alert string `json:"alert"`
// 			Title string `json:"title"`
// 		} `json:"android"`
// 	} `json:"notification"`
// }

// type JPushPayload struct {
// 	Platform     string         `json:"platform"`
// 	Audience     RegisrationIds `json:"audience"`
// 	Notification struct {
// 		Android struct {
// 			Alert string `json:"alert"`
// 			Title string `json:"title"`
// 		} `json:"android"`
// 	}   `json:"notification"`
// }

// type RegisrationIds struct {
// 	RegisrationIdList []string `json:"registration_id"`
// }
