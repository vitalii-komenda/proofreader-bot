package main

import (
	"encoding/json"
	"fmt"

	// "io"
	"net/http"
	"net/url"
	"sync"
)

var userTokens = sync.Map{}

func startOAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("OAuth start request received\n")

	config := parseConfig()
	authURL := "https://slack.com/oauth/v2/authorize"
	params := url.Values{}
	params.Add("client_id", config.SlackClientID)
	params.Add("scope", "chat:write")
	params.Add("redirect_uri", config.SlackRedirectURL)

	http.Redirect(w, r, fmt.Sprintf("%s?%s", authURL, params.Encode()), http.StatusFound)
}

func handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	accessToken, userID, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}

	userTokens.Store(userID, accessToken)
	t, _ := userTokens.Load(userID)

	fmt.Fprintf(w, "Access Token: %+v, user: %v", t, userID)
}

func exchangeCodeForToken(code string) (string, string, error) {
	config := parseConfig()

	resp, err := http.PostForm("https://slack.com/api/oauth.v2.access", url.Values{
		"client_id":     {config.SlackClientID},
		"client_secret": {config.SlackClientSecret},
		"code":          {code},
		"redirect_uri":  {config.SlackRedirectURL},
	})
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", "", err
	// }
	//
	// // Print the response body as text
	// fmt.Printf("Response body: %s\n", string(body))

	var response struct {
		AccessToken string `json:"access_token"`
		AuthedUser  struct {
			AccessToken string `json:"access_token"`
			UserID      string `json:"id"`
		} `json:"authed_user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", err
	}

	return response.AuthedUser.AccessToken, response.AuthedUser.UserID, nil
}
