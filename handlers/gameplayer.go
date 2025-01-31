package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type GamePlayerHandler struct {
	db *supabase.Client
}

func (h *GamePlayerHandler) GetGamePlayers(context *gin.Context) {

}

func (h *GamePlayerHandler) CreateGamePlayer(context *gin.Context) {

}

func (h *GamePlayerHandler) DeleteGamePlayer(context *gin.Context) {

}
