package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func loginPageHandler(ctx *gin.Context) {
	authzSess := sessions.DefaultMany(ctx, AUTHORIZATION_SESSION)
	if authzSess.Get(TICKET) == nil {
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ctx.HTML(http.StatusOK, LOGIN_TEMPLATE, gin.H{"action": LOGIN_ENDPOINT})
}

func loginAttemptHandler(ctx *gin.Context) {
	authzSess := sessions.DefaultMany(ctx, AUTHORIZATION_SESSION)
	authnSess := sessions.DefaultMany(ctx, AUTHENTICATION_SESSION)

	userId := ctx.Request.FormValue("user_id")
	password := ctx.Request.FormValue("password")
	user := findUser(userId)

	// Login Failed
	if !(user.Id == userId && user.Password == password) {
		ctx.HTML(http.StatusOK, LOGIN_TEMPLATE, gin.H{"action": LOGIN_ENDPOINT, "message": "Login Failed"})
		return
	}

	authnSess.Set(USER_ID, user.Id)
	authnSess.Options(sessions.Options{MaxAge: (60 * 60 * 24) * 30}) // One Month
	if err := authnSess.Save(); err != nil {
		logger.Error(err.Error())
	}

	if isUserConsented(userId, authzSess.Get(CLIENT_ID).(string)) {
		authorizationIssueCaller(ctx)
		return
	} else {
		ctx.Redirect(http.StatusFound, CONSENT_ENDPOINT)
		return
	}
}
