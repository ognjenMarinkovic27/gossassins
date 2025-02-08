package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

func main() {
	r := gin.Default()
	client, err := supabase.NewClient(os.Getenv("API_URL"), os.Getenv("API_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fmt.Println("jwt secret is empty")
		return
	}

	registerRoutes(r, client, jwtSecret)

	r.Run()
}
