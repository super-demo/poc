package main

import (
	"encoding/json"
	"log"
	sdk "mini-app-b/super-app-sdk"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Start the server first to ensure we have something listening on port 3001
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	// Log all requests
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("📥 Incoming request to: %s %s\n", c.Method(), c.Path())
		log.Printf("📄 Request body: %s\n", string(c.Body()))
		return c.Next()
	})

	// ✅ Function: Get User
	app.Post("/getUser", func(c *fiber.Ctx) error {
		log.Println("📥 getUser function called")

		var req map[string]interface{}

		// Print Raw Request Body
		log.Printf("📥 Raw Request Body: %s\n", string(c.Body()))

		if err := json.Unmarshal(c.Body(), &req); err != nil {
			log.Printf("❌ JSON Parsing Error: %v\n", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		log.Printf("✅ Parsed Request JSON: %v\n", req)

		// Extract userId with flexible type handling
		var userID float64

		switch v := req["userId"].(type) {
		case float64:
			userID = v
		case int:
			userID = float64(v)
		case json.Number:
			userID, _ = v.Float64()
		default:
			log.Printf("❌ Invalid userId type: %T\n", req["userId"])
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing or invalid userId"})
		}

		response := fiber.Map{"id": int(userID), "name": "John Doe", "email": "john@example.com"}
		log.Printf("📤 Sending response: %v\n", response)
		return c.JSON(response)
	})

	// ✅ Function: Get Settings
	app.Post("/getSettings", func(c *fiber.Ctx) error {
		log.Println("📥 getSettings function called")
		response := fiber.Map{"theme": "dark", "notifications": true}
		log.Printf("📤 Sending response: %v\n", response)
		return c.JSON(response)
	})

	// Start server in a goroutine
	go func() {
		log.Println("📦 Mini-App B running at http://localhost:3001")
		log.Println("✅ Ready to accept connections")
		if err := app.Listen(":3001"); err != nil {
			log.Fatalf("❌ Server error: %v\n", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)

	// Now register with the super app
	sdk := sdk.NewSuperAppSDK("super-secret-key")

	// Try registration multiple times if needed
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to register with Super App (attempt %d/%d)\n", i+1, maxRetries)

		err := sdk.Register("mini-app-b", []string{"getUser", "getSettings"})
		if err == nil {
			log.Println("✅ Mini-App B registered successfully")
			break
		}

		log.Printf("❌ Registration attempt %d failed: %v\n", i+1, err)

		if i < maxRetries-1 {
			log.Println("Waiting before retry...")
			time.Sleep(2 * time.Second)
		} else {
			log.Println("⚠️ All registration attempts failed, but continuing...")
		}
	}

	// Keep the server running
	select {}
}
