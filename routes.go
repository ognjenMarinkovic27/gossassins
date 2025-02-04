package main

import (
	"mognjen/gossassins/handlers"
	"mognjen/gossassins/repos"
	"mognjen/gossassins/services"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

func registerRoutes(r *gin.Engine, client *supabase.Client) {
	userRepo := repos.NewUserRepo(client)
	registerAuthRoutes(r, client, userRepo)
	registerUserRoutes(r, userRepo)
	registerGameRoutes(r, client)
}

func registerAuthRoutes(r *gin.Engine, client *supabase.Client, userRepo handlers.UserRepoForAuth) {
	authHandler := handlers.NewAuthHandler(client, userRepo)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/signup", authHandler.Signup)
	}
}

func registerUserRoutes(r *gin.Engine, userRepo handlers.UserRepo) {
	userHandler := handlers.NewUserHandler(userRepo)
	r.GET("/users/:id", userHandler.GetById)
}

func registerGameRoutes(r *gin.Engine, client *supabase.Client) {
	gameRepo := repos.NewGameRepo(client)
	gameHandler := handlers.NewGameHandler(gameRepo)
	gameGroup := r.Group("/games")
	{
		gameGroup.GET("/", gameHandler.GetAll)
		gameGroup.GET("/:id", gameHandler.GetById)
		gameGroup.POST("/", gameHandler.Create)
		gameGroup.PATCH("/:id", gameHandler.Patch)
		gameGroup.DELETE("/:id", gameHandler.Delete)

		registerJoinRequestRoutes(gameGroup, gameRepo, client)
		registerGameActionRoutes(gameGroup, gameRepo, client)
		registerGamePlayerRoutes(gameGroup, client)
	}
}

func registerJoinRequestRoutes(gameGroup *gin.RouterGroup, gameRepo *repos.GameRepo, client *supabase.Client) {
	joinRequestRepo := repos.NewJoinRequestRepo(client)
	joinRequestService := services.NewJoinRequestService(gameRepo, joinRequestRepo, client)
	joinRequestHandler := handlers.NewJoinRequestHandler(joinRequestService)
	joinRequestGroup := gameGroup.Group("/join-requests/:game_id")
	{
		joinRequestGroup.GET("/", joinRequestHandler.GetAllByGameId)
		joinRequestGroup.POST("/", joinRequestHandler.Create)
		joinRequestGroup.POST("/:user_id/approval", joinRequestHandler.Approve)
		joinRequestGroup.DELETE("/:user_id/approval", joinRequestHandler.Unapprove)
	}
}

func registerGameActionRoutes(gameGroup *gin.RouterGroup, gameRepo *repos.GameRepo, client *supabase.Client) {
	actionService := services.NewGameActionService(gameRepo, client)
	actionHandler := handlers.NewGameActionHandler(actionService)
	actionGroup := gameGroup.Group("/actions")
	{
		actionGroup.POST("/start", actionHandler.Start)
		actionGroup.POST("/kill", actionHandler.Kill)
	}
}

func registerGamePlayerRoutes(gameGroup *gin.RouterGroup, client *supabase.Client) {
	playerRepo := repos.NewGamePlayerRepo(client)
	playerHandler := handlers.NewGamePlayerHandler(client, playerRepo)
	playerGroup := gameGroup.Group("/players/:game_id")
	{
		playerGroup.GET("/", playerHandler.GetAllByGameId)
		playerGroup.GET("/:user_id", playerHandler.GetByGameIdUserId)
	}
}
