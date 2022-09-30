package pkg

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sandbox.go/cookies/internal/handler"
	"sandbox.go/cookies/internal/middleware"
	"syscall"
	"time"
)

func server(port int) *http.Server {
	router := gin.Default()
	router.Use(createSessionFunc())
	router.POST("/login", handler.LoginHandler)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	api := router.Group("/v1")
	api.Use(middleware.AuthorizeCookie)
	api.GET("/auth", func(ctx *gin.Context) {
		name, _ := ctx.Get("User.Name")
		id, idExists := ctx.Get("User.Id")
		if idExists {
			data := gin.H{"id": id, "name": name}
			ctx.JSON(http.StatusOK, data)
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

}

func createSessionFunc() gin.HandlerFunc {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/go_cookies?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	store := gormsessions.NewStore(db, true, []byte("secret"))
	return sessions.Sessions("name", store)
}

func StartApp() {
	port := 8080
	log.Printf("Starting server at Port %v", port)
	srv := server(port)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv.RegisterOnShutdown(func() {
		log.Println("Shutting down server")
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %srv\n", err)
		}
	}()
	log.Printf("Server Started at %d", port)
	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
