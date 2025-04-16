package main

import (
	"database/sql"
	"log"
	"net/http"
	"pvz/internal/delivery"
	"pvz/internal/delivery/middlewares"
	"pvz/internal/storage"
	"pvz/internal/storage/usecase"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:6432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	auth := usecase.NewAuthService("secret")
	userRepo := storage.NewUsersStorage(db)
	userUsecase := usecase.NewUserUsecase(userRepo, auth)
	loginHandler := delivery.NewLoginHandler(userUsecase)
	registerHandler := delivery.NewRegisterHandler(userUsecase)

	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/register", registerHandler.Register)
	r.POST("/login", loginHandler.Login)

	protected := r.Group("")
	protected.Use(middlewares.JWTAuthMiddleware(auth))
	{
		protected.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(200, "123")
		})
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
