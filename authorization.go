package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/authlete/authlete-go/dto"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func authorizationHandler(ctx *gin.Context) {
	authzRes, authleteErr := authleteApi.Authorization(&dto.AuthorizationRequest{Parameters: getParamsFromContext(ctx)})
	if authleteErr != nil {
		logger.Error(authleteErr.Error())
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	switch authzRes.Action {
	case dto.AuthorizationAction_INTERACTION:
		authorizationInteractionHandler(ctx, authzRes)
	case dto.AuthorizationAction_BAD_REQUEST:
		ctx.String(http.StatusBadRequest, authzRes.ResponseContent)
	case dto.AuthorizationAction_LOCATION:
		ctx.Redirect(http.StatusFound, authzRes.ResponseContent)
	case dto.AuthorizationAction_INTERNAL_SERVER_ERROR:
		ctx.String(http.StatusInternalServerError, authzRes.ResponseContent)
	default:
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}
}

func authorizationInteractionHandler(ctx *gin.Context, authzRes *dto.AuthorizationResponse) {
	authzSess := sessions.DefaultMany(ctx, AUTHORIZATION_SESSION)
	authzSess.Set(TICKET, authzRes.Ticket)
	clientId := strconv.FormatUint(authzRes.Client.ClientId, 10)
	if authzRes.ClientIdAliasUsed {
		clientId = authzRes.Client.ClientIdAlias
	}
	authzSess.Set(CLIENT_ID, clientId)
	scopes, _ := json.Marshal(authzRes.Scopes)
	authzSess.Set(SCOPES, scopes)

	if err := authzSess.Save(); err != nil {
		logger.Error(err.Error())
	}

	authnSess := sessions.DefaultMany(ctx, AUTHENTICATION_SESSION)
	if authnSess.Get(USER_ID) == nil {
		ctx.Redirect(http.StatusFound, LOGIN_ENDPOINT)
		return
	}

	if !isUserConsented(authnSess.Get(USER_ID).(string), clientId) {
		ctx.Redirect(http.StatusFound, CONSENT_ENDPOINT)
		return
	}

	authorizationIssueCaller(ctx)
}

func authorizationIssueCaller(ctx *gin.Context) {
	authzSess := sessions.DefaultMany(ctx, AUTHORIZATION_SESSION)
	ticket := authzSess.Get(TICKET).(string)
	authzSess.Options(sessions.Options{MaxAge: -1})
	if err := authzSess.Save(); err != nil {
		logger.Error(err.Error())
	}

	authnSess := sessions.DefaultMany(ctx, AUTHENTICATION_SESSION)
	subject := authnSess.Get(USER_ID).(string)
	authzIssueRes, authleteErr := authleteApi.AuthorizationIssue(&dto.AuthorizationIssueRequest{Ticket: ticket, Subject: subject})
	if authleteErr != nil {
		logger.Error(authleteErr.Error())
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	switch authzIssueRes.Action {
	case dto.AuthorizationIssueAction_BAD_REQUEST:
		ctx.String(http.StatusBadRequest, authzIssueRes.ResponseContent)
	case dto.AuthorizationIssueAction_LOCATION:
		ctx.Redirect(http.StatusFound, authzIssueRes.ResponseContent)
	case dto.AuthorizationIssueAction_INTERNAL_SERVER_ERROR:
		ctx.String(http.StatusInternalServerError, authzIssueRes.ResponseContent)
	default:
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}
}

func authorizationFailCaller(ctx *gin.Context, reason dto.AuthorizationFailReason) {
	authzSess := sessions.DefaultMany(ctx, AUTHORIZATION_SESSION)
	ticket := authzSess.Get(TICKET).(string)
	authzSess.Options(sessions.Options{MaxAge: -1})
	if err := authzSess.Save(); err != nil {
		logger.Error(err.Error())
	}

	authzFailRes, authleteErr := authleteApi.AuthorizationFail(&dto.AuthorizationFailRequest{Ticket: ticket, Reason: reason})
	if authleteErr != nil {
		logger.Error(authleteErr.Error())
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	switch authzFailRes.Action {
	case dto.AuthorizationFailAction_INTERNAL_SERVER_ERROR:
		ctx.String(http.StatusInternalServerError, authzFailRes.ResponseContent)
	case dto.AuthorizationFailAction_BAD_REQUEST:
		ctx.String(http.StatusBadRequest, authzFailRes.ResponseContent)
	case dto.AuthorizationFailAction_LOCATION:
		ctx.Redirect(http.StatusFound, authzFailRes.ResponseContent)
	default:
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}
}
