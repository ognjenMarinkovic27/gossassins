package main

import (
	"mognjen/gossassins/handlers"
	"mognjen/gossassins/handlers/middleware"
	"mognjen/gossassins/repos"
	"mognjen/gossassins/services"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

// TODO: this whole thing is kinda disgusting
func registerRoutes(r *gin.Engine, client *supabase.Client, jwtSecret string) {
	userRepo := repos.NewUserRepo(client)
	registerAuthRoutes(r, client, userRepo)

	r.Use(middleware.AuthMiddleware(jwtSecret))
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
	userGroup := r.Group("/users")
	{
		userGroup.GET("/:user_id", userHandler.GetById)
	}
}

func registerGameRoutes(r *gin.Engine, client *supabase.Client) {
	gameRepo := repos.NewGameRepo(client)
	gameHandler := handlers.NewGameHandler(gameRepo)
	gameGroup := r.Group("/games")
	{
		gameGroup.GET("/", gameHandler.GetAll)
		gameGroup.GET("/:game_id", gameHandler.GetById)
		gameGroup.POST("/", gameHandler.Create)

		ownerOnlyGroup := gameGroup.Group("/")
		ownerOnlyGroup.Use(middleware.IsGameOwnerMiddleware(gameRepo))
		ownerOnlyGroup.PATCH("/:game_id", gameHandler.Patch)
		ownerOnlyGroup.DELETE("/:game_id", gameHandler.Delete)

		registerGameActionRoutes(gameGroup, gameRepo, client)
		registerGamePlayerRoutes(gameGroup, gameRepo, client)
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

func registerGamePlayerRoutes(gameGroup *gin.RouterGroup, gameRepo *repos.GameRepo, client *supabase.Client) {
	playerRepo := repos.NewGamePlayerRepo(client)
	playerHandler := handlers.NewGamePlayerHandler(client, playerRepo)
	playerGroup := gameGroup.Group("/players/:game_id")
	playerGroup.Use(middleware.IsGameOwnerMiddleware(gameRepo))
	{
		playerGroup.GET("/", playerHandler.GetAllByGameId)
		playerGroup.POST("/", playerHandler.Create)
		playerGroup.GET("/:user_id", playerHandler.GetByGameIdUserId)
		playerGroup.PATCH("/:user_id", playerHandler.Patch)
	}
}
