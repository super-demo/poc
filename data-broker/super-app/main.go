// main.go
package main

import (
	"log"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		// Enable JSON parsing of requests with any content type
		// (Similar to Gin's behavior)
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Add logger middleware to log requests
	app.Use(logger.New())

	// Initialize our handlers
	NewSuperHandler(app)

	// Start the server
	log.Println("🚀 Server starting on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

// Mini-App Registry
var miniAppRegistry = make(map[string][]string)

func NewSuperHandler(app *fiber.App) {
	// Configure routes to match what the SDK expects
	app.Get("/v1/super/list", getlistMiniApp)
	app.Post("/v1/super/register", registerMiniApp)
	app.Post("/v1/super/call-function", callMiniAppFunction)
}

func registerMiniApp(c *fiber.Ctx) error {
	var data struct {
		AppName   string   `json:"appName"`
		Functions []string `json:"functions"`
	}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid registration data"})
	}
	miniAppRegistry[data.AppName] = data.Functions
	log.Printf("📌 Registered Mini-App: %s\n", data.AppName)
	log.Printf("📌 Functions: %v\n", miniAppRegistry)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Mini-App registered successfully!"})
}

func getlistMiniApp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"miniApps": miniAppRegistry,
	})
}

func callMiniAppFunction(c *fiber.Ctx) error {
	var req struct {
		Url          string         `json:"url"`
		Caller       string         `json:"caller"`
		TargetApp    string         `json:"targetApp"`
		FunctionName string         `json:"functionName"`
		Payload      map[string]any `json:"payload"`
	}

	// Log the raw request body
	requestData := c.Body()
	log.Printf("📥 Raw Request Body: %s\n", string(requestData))

	if err := c.BodyParser(&req); err != nil {
		log.Printf("❌ JSON Parsing Error: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	log.Printf("✅ Parsed Request: %+v\n", req)

	functions, exists := miniAppRegistry[req.TargetApp]
	if !exists || !contains(functions, req.FunctionName) {
		log.Printf("❌ Function not found: %s.%s\n", req.TargetApp, req.FunctionName)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Function not found"})
	}

	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		log.Printf("❌ Error encoding payload: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error encoding payload"})
	}

	// Try multiple possible hostnames to connect to mini-app-b
	urls := []string{
		// use this URL if you are running the server locally
		fmt.Sprintf("%s/%s", req.Url, req.FunctionName),

		// use this URL if you are running the server in a Docker container
		fmt.Sprintf("http://localhost:3001/%s", req.FunctionName),
		fmt.Sprintf("http://host.docker.internal:3001/%s", req.FunctionName),
	}

	var responseBody []byte
	var responseErr error
	var successful bool

	// Try each URL until one works
	for _, url := range urls {
		log.Printf("🔄 Trying to forward request to: %s with payload: %s\n", url, string(payloadBytes))

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("❌ Error creating request: %v\n", err)
			continue
		}

		request.Header.Set("Content-Type", "application/json")

		response, err := client.Do(request)
		if err != nil {
			log.Printf("❌ Error forwarding to %s: %v\n", url, err)
			responseErr = err
			continue
		}

		// Read the response body
		responseBody, err = ioutil.ReadAll(response.Body)
		response.Body.Close()

		if err != nil {
			log.Printf("❌ Error reading response from %s: %v\n", url, err)
			responseErr = err
			continue
		}

		log.Printf("✅ Successfully received response from %s: %s\n", url, string(responseBody))
		successful = true
		break
	}

	if !successful {
		log.Printf("❌ All connection attempts failed. Last error: %v\n", responseErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error forwarding request: %v", responseErr)})
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		log.Printf("❌ Error parsing response: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing response"})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
