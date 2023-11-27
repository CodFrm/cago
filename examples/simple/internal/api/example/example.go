package example

import "github.com/codfrm/cago/server/mux"

// LoginRequest 需要登录的接口
type LoginRequest struct {
	mux.Meta `path:"/example" method:"POST"`
}

// LoginResponse 需要登录的接口
type LoginResponse struct {
	Username string `json:"username"`
}

// PingRequest ping
type PingRequest struct {
	mux.Meta `path:"/example/ping" method:"GET"`
}

// PingResponse ping
type PingResponse struct {
	Pong string `json:"pong"`
}

// GinFunRequest gin function
type GinFunRequest struct {
	mux.Meta `path:"/example/gin-fun" method:"GET"`
}

// GinFunResponse gin function
type GinFunResponse struct {
}
