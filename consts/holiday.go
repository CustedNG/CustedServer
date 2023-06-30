package consts

import "github.com/LollipopKit/custed-server/model"

var (
	// TODO: 2023更新
	// 以下仅为template
	Holidays = []model.HolidaySkip{
		{
			From: &model.YearMonthDay{2022, 1, 1},
			To:   &model.YearMonthDay{2022, 1, 1},
		},
		{
			From: &model.YearMonthDay{2022, 10, 1},
		},
	}
)
