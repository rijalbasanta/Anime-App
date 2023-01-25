package userlibrary

import (
	"anime-app/users"
	"anime-app/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type EntityData struct {
	Id        string `json:"id"`
	Attribute struct {
		Status string `json:"status"`
		Name   string
		Type   string
	} `json:"attributes"`
	Relationship struct {
		User  interface{} `json:"user"`
		Anime interface{} `json:"anime"`
		Manga interface{} `json:"manga"`
	} `json:"relationships"`
}

type requestData struct {
	Attribute struct {
		Status string `json:"status"`
	} `json:"attributes"`
	Relationship struct {
		Anime struct {
			Data struct {
				Id   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"anime"`
		User struct {
			Data struct {
				Id   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"user"`
	} `json:"relationships"`
	Type string `json:"type"`
}

func getEntityByTitle(title string, entityType string) (string, error) {
	path := "https://kitsu.io/api/edge/"

	request, err := http.NewRequest("GET", path+entityType, nil)
	if err != nil {
		return "", fmt.Errorf("error creating get_entity request: %w", err)
	}
	params := request.URL.Query()
	params.Add("filter[text]", title)
	params.Add("fields["+entityType+"]", "id,canonicalTitle")
	params.Add("page[limit]", "1")
	params.Add("sort", "-averageRating")
	request.URL.RawQuery = params.Encode()
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error sending the get_entity request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the get_entity response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(string(responseBody), response.Status)
		return "", fmt.Errorf("server error: %w", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling the response: %w", err)
	}

	entityId := data["data"].([]interface{})[0].(map[string]interface{})["id"].(string)
	entityTitle := data["data"].([]interface{})[0].(map[string]interface{})["attributes"].(map[string]interface{})["canonicalTitle"].(string)

	if !strings.EqualFold(entityTitle, title) {
		fmt.Println("Proceeding with the entity with largest average rating: " + entityTitle)
	}

	return entityId, err
}

func CreateLibraryEntry(entityName string, entityType string, entityStatus string, userId string) (EntityData, error) {
	var data requestData
	var entity EntityData

	status := map[string]bool{"completed": true, "current": true, "planned": true}
	class := map[string]bool{"anime": true, "manga": true}

	if _, ok := class[entityType]; !ok {
		return entity, fmt.Errorf("entity type %v not supported", entityType)
	}
	if _, ok := status[entityStatus]; !ok {
		return entity, fmt.Errorf("entity status %v not supported", entityStatus)
	}

	id, err := getEntityByTitle(entityName, entityType)
	if err != nil {
		return entity, fmt.Errorf("error getting the enitity %v: %w", entityName, err)
	}

	data.Attribute.Status = entityStatus
	data.Relationship.Anime.Data.Id = id
	data.Relationship.Anime.Data.Type = entityType
	data.Relationship.User.Data.Id = userId
	data.Relationship.User.Data.Type = "users"
	data.Type = "library-entries"

	wrapper := map[string]requestData{"data": data}
	jsonData, err := json.Marshal(wrapper)
	if err != nil {
		return entity, fmt.Errorf("error json marshalling request data: %w", err)
	}

	requestBody := bytes.NewReader(jsonData)
	request, err := http.NewRequest("POST", "https://kitsu.io/api/edge/library-entries", requestBody)
	if err != nil {
		return entity, fmt.Errorf("error creating new user request: %w", err)
	}
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")
	request.Header.Set("Authorization", "Bearer "+utils.AuthToken.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return entity, fmt.Errorf("error sending the create_library_entry request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return entity, fmt.Errorf("error reading the create_library_entry response: %w", err)
	}
	if response.StatusCode != http.StatusCreated {
		fmt.Println(string(responseBody), response.Status)
		return entity, fmt.Errorf("server error: %w", err)
	}

	entData := make(map[string]EntityData)
	err = json.Unmarshal([]byte(responseBody), &entData)
	if err != nil {
		return entity, fmt.Errorf("error unmarshalling the create_library_entry response: %w", err)
	}
	entity = entData["data"]
	entity.Attribute.Name = entityName
	entity.Attribute.Type = entityType
	return entity, err
}

func GetLibraryEntry(userid string) ([]EntityData, error) {
	var entities []EntityData
	user, err := users.GetUser(userid)
	if err != nil {
		return entities, err
	}
	path := user.Relationship.(map[string]interface{})["libraryEntries"].(map[string]interface{})["links"].(map[string]interface{})["related"].(string)

	response, err := http.Get(path)
	if err != nil {
		return entities, fmt.Errorf("error requesting at %v: %w", path, err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return entities, fmt.Errorf("error reading the get_library response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(string(responseBody), response.Status)
		return entities, fmt.Errorf("server error: %w", err)
	}

	responseData := struct {
		Entities []EntityData `json:"data"`
	}{}
	err = json.Unmarshal([]byte(responseBody), &responseData)
	if err != nil {
		return entities, fmt.Errorf("error unmarshalling the get_library response: %w", err)
	}
	entities = responseData.Entities
	return entities, err
}

func DeleteLibraryEntry() {}
