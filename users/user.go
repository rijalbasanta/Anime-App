package users

import (
	"anime-app/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Attributes struct {
	UserName string `json:"name"`
	Email    string
}

type UserData struct {
	Attribute    Attributes  `json:"attributes"`
	Relationship interface{} `json:"relationships"`
	Id           string      `json:"id"`
}

type requestData struct {
	Attribute struct {
		UserName string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"attributes"`
	Type string `json:"type"`
}

// getTempMail gives a temporary email using the Guerrillamail Apis
func getTempMail() (string, error) {
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

// fileData is a template for writing user data, accessable only during creation, to a file
type fileData struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// writeUser writes the user data after creation to a JSON file for future references
func writeUser(user UserData) error {
	var data []map[string]fileData
	userFile, err := ioutil.ReadFile("./users/users.json")
	if userFile != nil {
		if err != nil {
			return fmt.Errorf("can't open user file: %w", err)
		}
		err = json.Unmarshal(userFile, &data)
		if err != nil {
			return fmt.Errorf("can't read user file: %w", err)
		}
	}
	n := strconv.Itoa(len(data))
	data = append(data, map[string]fileData{"user" + n: {
		Id:    user.Id,
		Name:  user.Attribute.UserName,
		Email: user.Attribute.Email,
	}})

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("can't marshal the data to write to the file: %w", err)
	}
	err = ioutil.WriteFile("./users/users.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("can't write to the file: %w", err)
	}
	return err
}

// CreateUser function creates a new user at Kitsu Backend and returns it
func CreateUser(name string) (UserData, error) {
	var err error
	var data requestData
	var user UserData
	data.Attribute.UserName = name
	data.Attribute.Email, err = getTempMail()
	if err != nil {
		return user, fmt.Errorf("error getting email: %w", err)
	}
	// data.Attribute.Email = "nenever249@nevyxus.com"
	data.Attribute.Password = "tempPassword" //please change
	data.Type = "users"

	wrapper := map[string]requestData{"data": data}
	jsonData, err := json.Marshal(wrapper)
	if err != nil {
		return user, fmt.Errorf("error json marshalling: %w", err)
	}

	requestBody := bytes.NewReader(jsonData)
	request, err := http.NewRequest("POST", "https://kitsu.io/api/edge/users", requestBody)
	if err != nil {
		return user, fmt.Errorf("error creating new user request: %w", err)
	}
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return user, fmt.Errorf("error sending the create_user request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, fmt.Errorf("error reading the create_user response: %w", err)
	}
	if response.StatusCode != http.StatusCreated {
		fmt.Println(string(responseBody), response.Status)
		return user, fmt.Errorf("server error: %w", err)
	}

	usrData := make(map[string]UserData)
	err = json.Unmarshal([]byte(responseBody), &usrData)
	if err != nil {
		return user, fmt.Errorf("error unmarshalling the create_user response: %w", err)
	}
	user = usrData["data"]
	user.Attribute.Email = data.Attribute.Email

	err = writeUser(user)
	if err != nil {
		fmt.Println(err)
	}

	return user, err
}

// GetUser gives the user with a particular id
func GetUser(id string) (UserData, error) {
	var user UserData

	request, err := http.NewRequest("GET", "https://kitsu.io/api/edge/users/"+id, nil)
	if err != nil {
		return user, fmt.Errorf("error creating get_user request: %w", err)
	}
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return user, fmt.Errorf("error sending the get_user request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, fmt.Errorf("error reading the get_user response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(string(responseBody), response.Status)
		return user, fmt.Errorf("server error: %w", err)
	}

	usrData := make(map[string]UserData)
	err = json.Unmarshal([]byte(responseBody), &usrData)
	if err != nil {
		return user, fmt.Errorf("error unmarshalling the create_user response: %w", err)
	}
	user = usrData["data"]

	var data []map[string]fileData
	userFile, err := ioutil.ReadFile("./users/users.json")
	if userFile != nil {
		if err != nil {
			return user, fmt.Errorf("can't open user file: %w", err)
		}
		err = json.Unmarshal(userFile, &data)
		if err != nil {
			return user, fmt.Errorf("can't read user file: %w", err)
		}

		for i, value := range data {
			if value["user"+strconv.Itoa(i)].Id == user.Id {
				user.Attribute.Email = value["user"+strconv.Itoa(i)].Email
			}
		}
	}
	return user, err
}

// TODO: Function to delete the user
func DeleteUser(id string) error {
	request, err := http.NewRequest("DELETE", "https://kitsu.io/api/edge/users/"+id, nil)
	if err != nil {
		return fmt.Errorf("error creating delete_user request: %w", err)
	}
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")
	request.Header.Set("Authorization", "Bearer "+utils.AuthToken.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending the delete_user request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusAccepted {
		return nil
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading the get_user response: %w", err)
	}
	fmt.Println(string(responseBody), response.Status)
	return fmt.Errorf("server error: %w", err)

}

// TODO: Function to update the user
func UpdateUser() error {
	return nil
}
