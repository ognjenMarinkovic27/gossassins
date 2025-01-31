package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type GameActionHandler struct {
	db *supabase.Client
}

func (h *GameActionHandler) GameActionKill(context *gin.Context) {

}
