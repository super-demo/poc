package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Mini-App Registry
var miniAppRegistry = make(map[string][]string)

func main() {
	app := fiber.New()

	// ✅ Mini-App Registration
	app.Post("/api/register", func(c *fiber.Ctx) error {
		var data struct {
			AppName   string   `json:"appName"`
			Functions []string `json:"functions"`
		}

		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid registration data"})
		}

		miniAppRegistry[data.AppName] = data.Functions
		log.Printf("📌 Registered Mini-App: %s\n", data.AppName)
		return c.JSON(fiber.Map{"message": "Mini-App registered successfully!"})
	})

	// ✅ Mini-App A calls Mini-App B's function
	app.Post("/api/call-function", func(c *fiber.Ctx) error {
		var req struct {
			Caller       string                 `json:"caller"`
			TargetApp    string                 `json:"targetApp"`
			FunctionName string                 `json:"functionName"`
			Payload      map[string]interface{} `json:"payload"`
		}

		// ✅ Parse request body correctly
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
		}

		// ✅ Ensure that the target function exists
		functions, exists := miniAppRegistry[req.TargetApp]
		if !exists || !contains(functions, req.FunctionName) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Function not found"})
		}

		// ✅ Convert payload to JSON format
		payloadBytes, err := json.Marshal(req.Payload)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error encoding payload"})
		}

		// ✅ Forward the request to the target mini-app
		url := fmt.Sprintf("http://localhost:3001/%s/%s", req.TargetApp, req.FunctionName)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error forwarding request"})
		}
		defer response.Body.Close()

		// ✅ Return the response from Mini-App B
		var result map[string]interface{}
		json.NewDecoder(response.Body).Decode(&result)
		return c.JSON(result)
	})

	log.Println("🏰 Super-App API Gateway running at http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

// Helper function to check if a slice contains a value
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
