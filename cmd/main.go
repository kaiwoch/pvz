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
	userRepo := storage.NewUsersStorage(db)
	pvzRepo := storage.NewPVZPostgresStorage(db)
	receptionRepo := storage.NewReceptionPostgresStorage(db)
	productRepo := storage.NewProductPostgresStorage(db)

	auth := usecase.NewAuthService("secret")
	receptionUsecase := usecase.NewReceptionUsecase(receptionRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, auth)
	pvzUsecase := usecase.NewPVZUsecase(pvzRepo, receptionRepo, productRepo)
	productUsecase := usecase.NewProductUsecase(productRepo, receptionRepo)

	loginHandler := delivery.NewLoginHandler(userUsecase)
	registerHandler := delivery.NewRegisterHandler(userUsecase)
	dummyLoginHandler := delivery.NewDummyLoginHandler(userUsecase)
	postPVZHandler := delivery.NewPVZHandler(pvzUsecase)
	receptionHandler := delivery.NewReceptionHandler(receptionUsecase)
	productHandler := delivery.NewProductHandler(productUsecase)

	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/register", registerHandler.Register)
	r.POST("/login", loginHandler.Login)
	r.POST("/dummyLogin", dummyLoginHandler.DummyLogin)

	protected := r.Group("")
	protected.Use(middlewares.JWTAuthMiddleware(auth))
	{
		protected.POST("/pvz", postPVZHandler.PostPVZ)
		protected.POST("/receptions", receptionHandler.Reception)
		protected.POST("/products", productHandler.Reception)
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
