package users

import (
	"anime-app/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type httpError struct {
	Code     int
	Response io.Reader
}

func (h *httpError) Error() string {

	var resp interface{}

	if err := json.NewDecoder(h.Response).Decode(&resp); err != nil {
		return err.Error()
	}

	data := map[string]interface{}{
		"code":     h.Code,
		"response": resp,
	}

	bts, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(bts)
}

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
			return fmt.Errorf("can't open user file: %v", err)
		}
		if len(userFile) > 0 {
			err = json.Unmarshal(userFile, &data)
			if err != nil {
				return fmt.Errorf("can't read user file: %v", err)
			}
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
		return fmt.Errorf("can't marshal the data to write to the file: %v", err)
	}
	err = ioutil.WriteFile("./users/users.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("can't write to the file: %v", err)
	}
	return err
}

// CreateUser function creates a new user at Kitsu Backend and returns it
func CreateUser(name string) (UserData, error) {
	var err error
	var data requestData
	var user UserData
	data.Attribute.UserName = name
	// data.Attribute.Email, err = utils.GetTempMail()
	// if err != nil {
	// 	return user, fmt.Errorf("error getting email: %v", err)
	// }
	data.Attribute.Email = "nenever249@nevyxus.com"
	data.Attribute.Password = "tempPassword" //please change
	data.Type = "users"

	wrapper := map[string]requestData{"data": data}
	jsonData, err := json.Marshal(wrapper)
	if err != nil {
		return user, fmt.Errorf("error json marshalling: %v", err)
	}

	requestBody := bytes.NewReader(jsonData)
	request, err := http.NewRequest("POST", "https://kitsu.io/api/edge/users", requestBody)
	if err != nil {
		return user, fmt.Errorf("error creating new user request: %v", err)
	}
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return user, fmt.Errorf("error sending the create_user request: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, fmt.Errorf("error reading the create_user response: %v", err)
	}
	if response.StatusCode != http.StatusCreated {
		fmt.Println(string(responseBody), response.Status)
		err := &httpError{
			Code:     response.StatusCode,
			Response: bytes.NewReader(responseBody),
		}
		return user, err
	}

	usrData := make(map[string]UserData)
	err = json.Unmarshal([]byte(responseBody), &usrData)
	if err != nil {
		return user, fmt.Errorf("error unmarshalling the create_user response: %v", err)
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
		return user, fmt.Errorf("error creating get_user request: %v", err)
	}
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return user, fmt.Errorf("error sending the get_user request: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, fmt.Errorf("error reading the get_user response: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(string(responseBody), response.Status)
		return user, fmt.Errorf("server error: %v", err)
	}

	usrData := make(map[string]UserData)
	err = json.Unmarshal([]byte(responseBody), &usrData)
	if err != nil {
		return user, fmt.Errorf("error unmarshalling the create_user response: %v", err)
	}
	user = usrData["data"]

	var data []map[string]fileData
	userFile, err := ioutil.ReadFile("./users/users.json")
	if userFile != nil {
		if err != nil {
			return user, fmt.Errorf("can't open user file: %v", err)
		}
		err = json.Unmarshal(userFile, &data)
		if err != nil {
			return user, fmt.Errorf("can't read user file: %v", err)
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
		return fmt.Errorf("error creating delete_user request: %v", err)
	}
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")
	request.Header.Set("Authorization", "Bearer "+utils.AuthToken.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending the delete_user request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusAccepted {
		return nil
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading the get_user response: %v", err)
	}
	fmt.Println(string(responseBody), response.Status)
	return fmt.Errorf("server error: %v", err)

}

// TODO: Function to update the user
func UpdateUser() error {
	return nil
}
