package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
)

type UserRole int64

func storeSession(ctx *gin.Context, userId int, userName string) error {
	session := sessions.Default(ctx)
	session.Set("user.Id", userId)
	session.Set("user.Name", userName)
	return session.Save()
}
func findUserId(dto loginDto) int {
	if dto.User == "root_user" && dto.Password == "root_password" {
		return rand.Int() + 1
	} else {
		return 0
	}
}

type loginDto struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

var LoginHandler gin.HandlerFunc = func(ctx *gin.Context) {
	var dto loginDto
	err := ctx.ShouldBind(&dto)
	if err != nil {
		log.Println(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	id := findUserId(dto)
	if id > 0 {
		err := storeSession(ctx, id, dto.User)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Ok",
			})
		}
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}
