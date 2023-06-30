<h1 align="center">Custed Server</h1>

<p align="center">
    <img alt="Build Status" src="https://badgen.net/badge/Go/Echo/cyan?icon=github" href="https://github.com/labstack/echo">
</p>

## 开发
1. 安装Go环境1.18+[^1]
   - （可选）若无梯子，请使用goproxy[^2]
2. 克隆至本地
3. 运行项目
   - `go build`，会生成一个叫 `custed-server` 的二进制文件
   - 执行`./custed-server`，就会在**1327**端口运行custed的后端服务
4. 使用测试服务器
   - 找到CustedNG内后端域名：[constants.dart](https://github.com/CustedNG/CustedNG/blob/material/lib/constants.dart#L19) 内的 `backendUrl`
   - 从 `https://api.backend.cust.team` 改为 `http://你运行二进制的机器的IP:1327` 进行真机调试
5. 二进制具体用法，可以执行 `./custed-server --help` 来查看帮助

## 功能
- 接受、发送后端缓存的课表、考表
- 作为推送服务器，收集token，推送通知
- 首页提醒、app更新、有无考试、校历等....

## TODO
- [x] 完善接口信息，标准化api.txt
- [x] 完善项目内的注释
- [ ] 标准化课表结构（教务课表、kbpro）


[^1]:https://golang.google.cn
[^2]:https://goproxy.io/zh/