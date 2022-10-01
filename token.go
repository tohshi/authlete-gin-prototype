package main

import (
	"net/http"

	"github.com/authlete/authlete-go/dto"
	"github.com/gin-gonic/gin"
)

func tokenHandler(ctx *gin.Context) {

	username, password, _ := ctx.Request.BasicAuth()
	params := getParamsFromContext(ctx)

	tokenRes, authleteErr := authleteApi.Token(&dto.TokenRequest{Parameters: params, ClientId: username, ClientSecret: password})
	if authleteErr != nil {
		logger.Error(authleteErr.Error())
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	switch tokenRes.Action {
	case dto.TokenAction_OK:
		ctx.String(http.StatusOK, tokenRes.ResponseContent)
	case dto.TokenAction_INVALID_CLIENT:
		ctx.String(http.StatusUnauthorized, tokenRes.ResponseContent)
	case dto.TokenAction_INTERNAL_SERVER_ERROR:
		ctx.String(http.StatusInternalServerError, tokenRes.ResponseContent)
	case dto.TokenAction_BAD_REQUEST:
		ctx.String(http.StatusBadRequest, tokenRes.ResponseContent)
	default:
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}
}
