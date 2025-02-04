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

		registerGameApprovalRoutes(gameGroup, client)
		registerGameActionRoutes(gameGroup, gameRepo, client)
	}
}

func registerGameApprovalRoutes(gameGroup *gin.RouterGroup, client *supabase.Client) {
	approvalRepo := repos.NewGameApprovalRepo(client)
	approvalHandler := handlers.NewGameApprovalHandler(approvalRepo)
	approvalGroup := gameGroup.Group("/approvals/:game_id")
	{
		approvalGroup.GET("/", approvalHandler.GetAllByGameId)
		approvalGroup.POST("/", approvalHandler.Create)
		approvalGroup.PATCH("/:user_id", approvalHandler.Patch)
	}
}

func registerGameActionRoutes(gameGroup *gin.RouterGroup, gameRepo *repos.GameRepo, client *supabase.Client) {
	actionService := services.NewGameActionService(gameRepo, client)
	actionHandler := handlers.NewGameActionHandler(actionService)
	actionGroup := gameGroup.Group("/actions")
	{
		actionGroup.POST("/start", actionHandler.Start)

		/*
		   TODO: this could be a REST style approval resource
		   and current approvals could be renamed to approval
		   requests
		*/
		actionGroup.POST("/approve", actionHandler.Approve)
		actionGroup.POST("/unapprove", actionHandler.Unapprove)

		actionGroup.POST("/kill", actionHandler.Kill)
	}
}

func registerGamePlayerRoutes(r *gin.Engine, client *supabase.Client) {
	playerRepo := repos.NewGamePlayerRepo(client)
	playerHandler := handlers.NewGamePlayerHandler(client, playerRepo)
	playerGroup := r.Group("/game-player/:game_id")
	{
		playerGroup.GET("/", playerHandler.GetAllByGameId)
		playerGroup.DELETE("/:player_id", playerHandler.Delete)
	}
}
