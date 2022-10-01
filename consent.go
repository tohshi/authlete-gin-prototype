package main

import (
	"encoding/json"
	"net/http"

	"github.com/authlete/authlete-go/dto"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func isUserConsented(userId string, clientId string) bool {
	return userId == "user2"
}

func consentPageHandler(ctx *gin.Context) {
	authzSess := sessions.DefaultMany(ctx, AUTHORIZATION_SESSION)
	if authzSess.Get(TICKET) == nil {
		ctx.String(http.StatusBadRequest, "Bad Request")
		return
	}

	scopes := []dto.Scope{}
	json.Unmarshal(authzSess.Get(SCOPES).([]byte), &scopes)

	ctx.HTML(http.StatusOK, CONSENT_TEMPLATE, gin.H{"action": CONSENT_ENDPOINT, "scopes": scopes})
}

func consentAttemptHandler(ctx *gin.Context) {
	var consent bool
	switch ctx.Request.FormValue("consent") {
	case "agree":
		consent = true
	case "reject":
		consent = false
	default:
		ctx.String(http.StatusBadRequest, "Bad Request")
		return
	}

	if consent {
		authorizationIssueCaller(ctx)
	} else {
		authorizationFailCaller(ctx, dto.AuthorizationFailReason_DENIED)
	}
}
