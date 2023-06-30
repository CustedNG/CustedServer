# config.json 字段说明

```json
{
    // app 通知
    // 客户端需选择版本大于等于`version`的通知
    // `content`需以'。n月m日'结尾
    // `url`为通知详情页面的URL
    "notify": [
        {
            "version": 695,
            "content": "example",
            "url": "https://example.com"
        },
        {
            "version": 580,
            "content": "请更新至最新版本。5月1日"
            "url": "https://example.com"
        }
    ],
    // 教务URL，客户端需从中随机选择一个
    "jw_url": [
        "https://jwgls0.cust.edu.cn",
        "https://jwgls1.cust.edu.cn",
        "https://jwgls2.cust.edu.cn",
        "https://jwgls3.cust.edu.cn"
    ],
    // 客户端最新版本信息
    "update": {
        "version": {
            "android": 695,
            "ios": 695
        },
        // 升级优先级。
        // 0：不弹窗、不Snackbar
        // 1：Snackbar显示
        // 2：弹窗
        "priority": {
            "android": 1,
            "ios": 1
        },
        "url": {
            "android": "",
            "ios": ""
        },
        "changelog": {
            "android": "新UI",
            "ios": "修复xxx"
        }
    },
    // int list
    // 显示fake data的版本，用来通过AppStore审核（登录需要验证码）
    "not_show_real_ui": [
        695
    ],
    // str list
    // 当教务无法使用时，使用kbpro数据源
    // 可以用`695`单个数字，或者`695-696`（上下限都包含）
    "use_kbpro": [
        "695", "695-696"
    ],
    // 客户端内需显示为可点击的文字
    // name为显示内容，url为跳转链接
    // 当url为空时，点击无反应
    "tester_list": [
        {
            "name": "LollipopKit",
            "url": "https://github.com/LollipopKit"
        }
    ],
    // img_url为banner的图片链接
    // action_url为跳转链接，为空时点击无反应
    "banner": {
        "img_url": "https://cdn3.cust.app/banner/example.png",
        "action_url": ""
    },
    // 是否有考试/是否显示首页的考试卡片
    "have_exam": true,
    // // 长春天气
    // "weather": {
    //     // 晴、雨、雾等
    //     "condition": "晴",
    //     // 全天温度
    //     "day_temp": "6-15",
    //     // 现在的温度
    //     "now_temp": "9"
    // },
    // 校历
    "school_calendar": []
}
```