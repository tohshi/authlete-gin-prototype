package main

import (
	"io"

	"github.com/gin-gonic/gin"
)

func getParamsFromContext(ctx *gin.Context) string {
	var params string

	switch ctx.Request.Method {
	case "GET":
		params = ctx.Request.URL.RawQuery
	case "POST":
		bytes, _ := io.ReadAll(ctx.Request.Body)
		params = string(bytes)
	}

	return params
}
