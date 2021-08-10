package lodestone

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func Login() {
	url, _ := url.Parse("https://finalfantasyxiv.com")
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	stored, csrfToken := getStoredAndCsrfToken(client)
	usernameAndPassword(client, stored, csrfToken)
	cis_sessid := oneTimePassword(client, stored, csrfToken)
	characterMap := finishLogin(client, cis_sessid, csrfToken)
	selectCharacter(client, characterMap, csrfToken)

	for _, cookie := range jar.Cookies(url) {
		if cookie.Name == "ldst_sess" {
			viper.Set("lodestone_session_token", cookie.Value)
			viper.WriteConfig()
			fmt.Println("Logged in to Lodestone succesfully!")
			return
		}
	}

	fmt.Printf("Error logging in to Lodestone!")
}

func selectCharacter(client *http.Client, characterMap map[string]string, csrfToken string) {
	var characterNames []string

	for characterName := range characterMap {
		characterNames = append(characterNames, characterName)
	}

	prompt := promptui.Select{
		Label: "Select A Character",
		Items: characterNames,
	}

	_, result, _ := prompt.Run()

	characterId := characterMap[result]
	timestamp := time.Now().UnixNano() / 1000000

	url := "https://na.finalfantasyxiv.com/lodestone/api/account/select_character/"
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf("csrf_token=%s&character_id=%s&__timestamp=%d", csrfToken, characterId, timestamp))

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = client.Do(req)
}

func finishLogin(client *http.Client, cis_sessid string, csrfToken string) map[string]string {
	characterMap := make(map[string]string)

	url := fmt.Sprintf("https://login.finalfantasyxiv.com/lodestone/account/login_back_reload?rgn=na&csrf_token=%s", csrfToken)
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf("cis_sessid=%s", cis_sessid))

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	characterElements := doc.Find(".entry")

	characterElements.Each(func(_ int, characterElement *goquery.Selection) {
		name := characterElement.Find(".entry__name").Text()
		id, _ := characterElement.Find(".bt_character_login").Attr("data-character_id")
		characterMap[name] = id
	})

	return characterMap
}

func oneTimePassword(client *http.Client, stored string, csrfToken string) string {
	oneTimePasswordPrompt := promptui.Prompt{
		Label: "Lodestone OTP",
	}
	oneTimePassword, err := oneTimePasswordPrompt.Run()
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("https://secure.square-enix.com/oauth/oa/oauthlogin.sendOtp?response_type=code&lang=en-us&redirect_uri=https://login.finalfantasyxiv.com/lodestone/account/login_back/?rgn=na&csrf_token=%s&client_id=ff14lodestone&alar=1", csrfToken)
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf("_STORED_=%s&otppw=%s&wfp=1", stored, oneTimePassword))

	req, err := http.NewRequest(method, url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	cis_sessid, _ := doc.Find("input[name='cis_sessid']").Attr(("value"))

	return cis_sessid
}

func usernameAndPassword(client *http.Client, stored string, csrfToken string) {
	usernamePrompt := promptui.Prompt{
		Label: "Lodestone Username",
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		panic(err)
	}

	passwordPrompt := promptui.Prompt{
		Label: "Lodestone Password",
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		panic(err)
	}

	url2 := fmt.Sprintf("https://secure.square-enix.com/oauth/oa/oauthlogin.send?response_type=code&lang=en-us&redirect_uri=https://login.finalfantasyxiv.com/lodestone/account/login_back/?rgn=na&csrf_token=%s&client_id=ff14lodestone&alar=1", csrfToken)
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf("_STORED_=%s&sqexid=%s&password=%s&wfp=1", stored, url.QueryEscape(username), url.QueryEscape(password)))

	req, err := http.NewRequest(method, url2, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req)

	if err != nil {
		panic(err)
	}
}

func getStoredAndCsrfToken(client *http.Client) (string, string) {
	req, err := http.NewRequest("GET", "https://na.finalfantasyxiv.com/lodestone/account/login/", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	stored, _ := doc.Find("input[name='_STORED_']").Attr(("value"))
	action, _ := doc.Find("form[name='mainForm']").Attr("action")
	csrf_token := strings.Split(strings.Split(action, "csrf_token%3D")[1], "&client_id")[0]

	return stored, csrf_token
}

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
