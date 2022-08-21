package e

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	ErrorAuthCheckTokenFail    = 20001 // token 错误
	ErrorAuthCheckTokenTimeout = 20002 // token 过期
)
