package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SuperAppSDK struct {
	APIKey  string
	BaseURL string
}

func NewSuperAppSDK(apiKey string) *SuperAppSDK {
	return &SuperAppSDK{
		APIKey:  apiKey,
		BaseURL: "http://localhost:3000/api",
	}
}

// ✅ Register Mini-App
func (sdk *SuperAppSDK) Register(appName string, functions []string) error {
	payload, _ := json.Marshal(map[string]interface{}{
		"appName":   appName,
		"functions": functions,
	})

	resp, err := http.Post(sdk.BaseURL+"/register", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ✅ Call Another Mini-App's Function (🔴 Fixes JSON Encoding Issue)
func (sdk *SuperAppSDK) CallFunction(caller, targetApp, functionName string, payload map[string]interface{}) (map[string]interface{}, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"caller":       caller,
		"targetApp":    targetApp,
		"functionName": functionName,
		"payload":      payload,
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding request JSON: %v", err)
	}

	resp, err := http.Post(sdk.BaseURL+"/call-function", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error calling function: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response JSON: %v", err)
	}
	return result, nil
}
