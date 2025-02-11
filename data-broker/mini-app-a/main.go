package main

import (
	"fmt"
	"log"
	sdk "mini-app-a/super-app-sdk"
)

func main() {
	sdk := sdk.NewSuperAppSDK("super-secret-key")

	// ✅ Register Mini-App A
	sdk.Register("mini-app-a", []string{"getProfile", "getBalance"})

	// ✅ Call Mini-App B's getUser function
	result, err := sdk.CallFunction("mini-app-a", "mini-app-b", "getUser", map[string]interface{}{"userId": 123})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("👤 User Data from Mini-App B:", result)

	result2, err := sdk.CallFunction("mini-app-a", "mini-app-b", "getSettings", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("⚙️ Settings from Mini-App B:", result2)
}
