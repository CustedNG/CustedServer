package consts

const (
	// 避免0（golang中的零值）
	CodeSuccess      = -1
	CodeUnknownError = iota
	CodeNoPermission
	CodeInvalidParam
	CodeInvalidBody
	CodeInvalidJson
	CodeDBOperationFailed
	CodeServerHTTPRequestFailed
	CodePushFailed
	CodeInternal
)
