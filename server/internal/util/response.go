package util

import "github.com/gin-gonic/gin"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Code: 0, Msg: "success", Data: data})
}

func Err(c *gin.Context, code int, msg string) {
	c.JSON(200, Response{Code: code, Msg: msg, Data: nil})
}
