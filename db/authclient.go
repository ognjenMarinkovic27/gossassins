package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/supabase-community/supabase-go"
)

type AuthSupabaseClient struct {
	*supabase.Client
}

func (c *AuthSupabaseClient) AuthWithMagicLink(email string, createUser bool) error {
	url := fmt.Sprintf("%s/auth/v1/otp", API_URL)

	// Create the JSON payload for the request
	payload := map[string]interface{}{
		"email":       email,
		"create_user": createUser,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Send POST request to Supabase
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+API_KEY)
	req.Header.Set("apikey", API_KEY)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed to send magic link, status: " + resp.Status)
	}

	fmt.Println("Magic link sent successfully!")
	return nil
}
