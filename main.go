package main

import (
	"net/http"
	"os"

	"github.com/authlete/authlete-go/api"
	"github.com/authlete/authlete-go/conf"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var authleteApi api.AuthleteApi
var logger *zap.Logger

func init() {
	llogger, _ := zap.NewProduction()
	logger = llogger

	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file")
	}

	conf := new(conf.AuthleteEnvConfiguration)
	authleteApi = api.New(conf)
}

func main() {
	r := gin.New()

	r.LoadHTMLGlob("templates/*.tmpl")

	r.Use(ginzapMiddleware())
	r.Use(ginzap.RecoveryWithZap(logger, true))

	store, _ := redis.NewStore(10, os.Getenv("REDIS_NETWORK"), os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"), []byte(os.Getenv("REDIS_KEY")))
	r.Use(sessions.SessionsMany([]string{AUTHORIZATION_SESSION, AUTHENTICATION_SESSION}, store))

	r.GET(PING_ENDPOINT, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	r.GET(AUTHORIZATION_ENDPINT, authorizationHandler)
	r.POST(AUTHORIZATION_ENDPINT, authorizationHandler)
	r.POST(TOKEN_ENDPOINT, tokenHandler)
	r.GET(LOGIN_ENDPOINT, loginPageHandler)
	r.POST(LOGIN_ENDPOINT, loginAttemptHandler)
	r.GET(CONSENT_ENDPOINT, consentPageHandler)
	r.POST(CONSENT_ENDPOINT, consentAttemptHandler)

	r.Run(":" + os.Getenv("PORT"))
}
