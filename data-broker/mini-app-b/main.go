package main

import (
	"encoding/json"
	"fmt"
	"log"
	sdk "mini-app-b/super-app-sdk"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	sdk := sdk.NewSuperAppSDK("super-secret-key")

	// ✅ Register Mini-App B
	sdk.Register("mini-app-b", []string{"getUser", "getSettings"})

	app := fiber.New()

	app.Post("/mini-app-b/getUser", func(c *fiber.Ctx) error {
		var req map[string]interface{}

		// 🛠️ Debugging: Print Raw Request Body
		fmt.Println("📥 Raw Request Body:", string(c.Body()))

		// 🛠️ Ensure the JSON format is correct
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			fmt.Println("❌ JSON Parsing Error:", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		// 🛠️ Debugging: Print Parsed JSON
		fmt.Println("✅ Parsed Request JSON:", req)

		// Extract userId
		userID, exists := req["userId"].(float64) // JSON numbers are float64 by default
		if !exists {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing or invalid userId"})
		}

		// Return user data
		return c.JSON(fiber.Map{"id": int(userID), "name": "John Doe", "email": "john@example.com"})
	})

	// ✅ Function: Get Settings
	app.Post("/mini-app-b/getSettings", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"theme": "dark", "notifications": true})
	})

	log.Println("📦 Mini-App B running at http://localhost:3001")
	log.Fatal(app.Listen(":3001"))
}
