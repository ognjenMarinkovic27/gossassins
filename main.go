package main

import (
	"fmt"
	"mognjen/gossassins/db"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

func main() {
	r := gin.Default()
	client, err := supabase.NewClient(db.API_URL, db.API_KEY, &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}

	registerRoutes(r, client)

	r.Run()
}
