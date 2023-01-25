package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Authentication struct {
	Token string `json:"access_token"`
}

var AuthToken Authentication

func Authenticate(email string) error {
	jsonData, err := json.Marshal(struct {
		Type     string `json:"grant_type"`
		Email    string `json:"username"`
		Password string `json:"password"`
	}{
		Type:     "password",
		Email:    email,
		Password: "tempPassword",
	})
	if err != nil {
		return fmt.Errorf("error marshalling email and password: %w", err)
	}

	requestBody := bytes.NewReader(jsonData)
	request, err := http.NewRequest("POST", "https://kitsu.io/api/oauth/token", requestBody)
	if err != nil {
		return fmt.Errorf("error creating authentication request: %w", err)
	}
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending the authentication request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading the authentication response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(string(responseBody), response.Status)
		return fmt.Errorf("server error: %w", err)
	}

	err = json.Unmarshal([]byte(responseBody), &AuthToken)
	if err != nil {
		return fmt.Errorf("error unmarshalling the authentication response: %w", err)
	}

	return err
}
