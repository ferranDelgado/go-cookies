package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getInt(value interface{}) int {
	if value == nil {
		return 0
	} else {
		return value.(int)
	}
}

func getString(value interface{}) string {
	if value == nil {
		return ""
	} else {
		return value.(string)
	}
}

var AuthorizeCookie gin.HandlerFunc = func(c *gin.Context) {
	session := sessions.Default(c)
	userId := getInt(session.Get("user.Id"))
	userName := getString(session.Get("user.Name"))
	if userId > 0 {
		c.Set("User.Id", userId)
		c.Set("User.Name", userName)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
