package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// getTempMail gives a temporary email using the Guerrillamail Apis
func GetTempMail() (string, error) {
	mailResponse, err := http.Get("http://api.guerrillamail.com/ajax.php?f=get_email_address")
	if err != nil {
		return "", fmt.Errorf("error getting temp-mail response: %w", err)
	}
	defer mailResponse.Body.Close()
	mailBody, err := ioutil.ReadAll(mailResponse.Body)
	if err != nil {
		return "", fmt.Errorf("error reading temp-mail response: %w", err)
	}

	// Using hardcoded indexes to extract the email address only
	return string(mailBody[15:46]), err
}
