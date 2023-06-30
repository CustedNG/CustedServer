package model

type Weather struct {
	Ok   bool `json:"ok"`
	Data struct {
		City       string `json:"city"`
		Updatetime string `json:"updatetime"`
		Wendu      string `json:"wendu"`
		Fengli     string `json:"fengli"`
		Shidu      string `json:"shidu"`
		Fengxiang  string `json:"fengxiang"`
		Sunrise1   string `json:"sunrise_1"`
		Sunset1    string `json:"sunset_1"`
		Sunrise2   struct {
		} `json:"sunrise_2"`
		Sunset2 struct {
		} `json:"sunset_2"`
		Yesterday struct {
			Date1 string `json:"date_1"`
			High1 string `json:"high_1"`
			Low1  string `json:"low_1"`
			Day1  struct {
				Type1 string `json:"type_1"`
				Fx1   string `json:"fx_1"`
				Fl1   string `json:"fl_1"`
			} `json:"day_1"`
			Night1 struct {
				Type1 string `json:"type_1"`
				Fx1   string `json:"fx_1"`
				Fl1   string `json:"fl_1"`
			} `json:"night_1"`
		} `json:"yesterday"`
		Forecast struct {
			Weather []struct {
				Date string `json:"date"`
				High string `json:"high"`
				Low  string `json:"low"`
				Day  struct {
					Type      string `json:"type"`
					Fengxiang string `json:"fengxiang"`
					Fengli    string `json:"fengli"`
				} `json:"day"`
				Night struct {
					Type      string `json:"type"`
					Fengxiang string `json:"fengxiang"`
					Fengli    string `json:"fengli"`
				} `json:"night"`
			} `json:"weather"`
		} `json:"forecast"`
		Zhishus struct {
			Zhishu []struct {
				Name   string `json:"name"`
				Value  string `json:"value"`
				Detail string `json:"detail"`
			} `json:"zhishu"`
		} `json:"zhishus"`
	} `json:"data"`
}
