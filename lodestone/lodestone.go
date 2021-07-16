package lodestone

import (
	"fmt"
	"net/http"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

// Gets a session token from Lodestone to read private data
func GetSessionToken() {
	// Some crap I wrote to see if we could get our own tokens
	// Maybe I'll work more on this later?

	// idPrompt := promptui.Prompt{
	// 	Label: "Square Enix ID",
	// }
	// id, err := idPrompt.Run()
	// if err != nil {
	// 	fmt.Printf("Prompt failed %v\n", err)
	// 	return ""
	// }
	// if id == "" {
	// 	return ""
	// }

	// passwordPrompt := promptui.Prompt{
	// 	Label: "Square Enix Password",
	// 	Mask:  '*',
	// }
	// password, err := passwordPrompt.Run()
	// if err != nil {
	// 	fmt.Printf("Prompt failed %v\n", err)
	// 	return ""
	// }
	// if password == "" {
	// 	return ""
	// }

	// fmt.Printf("Square Enix ID: %s\n", id)
	// fmt.Printf("Square Enix Password: %s\n", password)
	// return ""
	tokenPrompt := promptui.Prompt{
		Label: "Lodestone Session Token",
	}
	lodestone_session_token, err := tokenPrompt.Run()
	if err != nil {
		panic(err)
	}
	viper.Set("lodestone_session_token", lodestone_session_token)
	viper.WriteConfig()
}

// Sets up a request
func SetupRequest(endpoint string, character_id string) *http.Request {
	lodestone_session_token := viper.Get("lodestone_session_token")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://na.finalfantasyxiv.com/lodestone/character/%s/%s", character_id, endpoint), nil)
	if err != nil {
		panic(err)
	}

	if lodestone_session_token != "" {
		req.Header.Set("Cookie", fmt.Sprintf("ldst_sess=%s;", lodestone_session_token))
	}

	return req
}

// Sets up a request
func SetupMobileRequest(endpoint string, character_id string) *http.Request {
	USER_AGENT := "Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.76 Mobile Safari/537.36"

	lodestone_session_token := viper.Get("lodestone_session_token")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://na.finalfantasyxiv.com/lodestone/character/%s/%s", character_id, endpoint), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", USER_AGENT)

	if lodestone_session_token != "" {
		req.Header.Set("Cookie", fmt.Sprintf("ldst_sess=%s;", lodestone_session_token))
	}

	return req
}
