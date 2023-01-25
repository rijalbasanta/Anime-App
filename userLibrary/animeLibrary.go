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

type EntryData struct {
	Id        string `json:"id"`
	Attribute struct {
		Status string `json:"status"`
		Name   string
		Type   string
		User   string
	} `json:"attributes"`
	Relationship struct {
		User  interface{} `json:"user"`
		Anime interface{} `json:"anime"`
		Manga interface{} `json:"manga"`
	} `json:"relationships"`
}

type animeRequestData struct {
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

type mangaRequestData struct {
	Attribute struct {
		Status string `json:"status"`
	} `json:"attributes"`
	Relationship struct {
		Manga struct {
			Data struct {
				Id   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"manga"`
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
		return "", fmt.Errorf("error creating get_entity request: %v", err)
	}
	params := request.URL.Query()
	params.Add("filter[text]", title)
	params.Add("fields["+entityType+"]", "id,canonicalTitle")
	params.Add("page[limit]", "1")
	request.URL.RawQuery = params.Encode()
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error sending the get_entity request: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the get_entity response: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(string(responseBody), response.Status)
		return "", fmt.Errorf("server error: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling the response: %v", err)
	}

	entityId := data["data"].([]interface{})[0].(map[string]interface{})["id"].(string)
	entityTitle := data["data"].([]interface{})[0].(map[string]interface{})["attributes"].(map[string]interface{})["canonicalTitle"].(string)

	if !strings.EqualFold(entityTitle, title) {
		fmt.Println("Proceeding with the most relevent entity: " + entityTitle)
	}

	return entityId, err
}

func getDataByEntry(entity EntryData) (userName, entryName, entryType string, err error) {
	userPath := entity.Relationship.User.(map[string]interface{})["links"].(map[string]interface{})["related"].(string)
	animePath := entity.Relationship.Anime.(map[string]interface{})["links"].(map[string]interface{})["related"].(string)
	mangaPath := entity.Relationship.Manga.(map[string]interface{})["links"].(map[string]interface{})["related"].(string)

	response, err := http.Get(userPath)
	if err != nil {
		return userName, entryName, entryType, fmt.Errorf("error getting user: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return userName, entryName, entryType, fmt.Errorf("error reading the get_user response: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(string(responseBody), response.Status)
		return userName, entryName, entryType, fmt.Errorf("server error: %v", err)
	}
	user := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &user)
	if err != nil {
		return userName, entryName, entryType, fmt.Errorf("error unmarshalling the get_user response: %v", err)
	}

	entryType = "anime"
	response, err = http.Get(animePath)
	if err != nil {
		return userName, entryName, entryType, fmt.Errorf("error getting anime: %v", err)
	}
	defer response.Body.Close()
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return userName, entryName, entryType, fmt.Errorf("error reading the get_anime response: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(string(responseBody), response.Status)
		return userName, entryName, entryType, fmt.Errorf("server error: %v", err)
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return userName, entryName, entryType, fmt.Errorf("error unmarshalling the get_anime response: %v", err)
	}
	if (data["data"]) == nil {
		entryType = "manga"
		response, err = http.Get(mangaPath)
		if err != nil {
			return userName, entryName, entryType, fmt.Errorf("error getting manga: %v", err)
		}
		defer response.Body.Close()
		responseBody, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return userName, entryName, entryType, fmt.Errorf("error reading the get_manga response: %v", err)
		}
		if response.StatusCode != http.StatusOK {
			fmt.Println(string(responseBody), response.Status)
			return userName, entryName, entryType, fmt.Errorf("server error: %v", err)
		}
		data = make(map[string]interface{})
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			return userName, entryName, entryType, fmt.Errorf("error unmarshalling the get_manga response: %v", err)
		}

	}

	userName = user["data"].(map[string]interface{})["attributes"].(map[string]interface{})["name"].(string)
	entryName = data["data"].(map[string]interface{})["attributes"].(map[string]interface{})["canonicalTitle"].(string)

	return userName, entryName, entryType, nil
}

func CreateLibraryEntry(entityName string, entityType string, entityStatus string, userId string) (EntryData, error) {
	var animeData animeRequestData
	var mangaData mangaRequestData
	var entity EntryData

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
		return entity, fmt.Errorf("error getting the enitity %v: %v", entityName, err)
	}

	var jsonData []byte
	if entityType == "anime" {
		animeData.Attribute.Status = entityStatus
		animeData.Relationship.Anime.Data.Id = id
		animeData.Relationship.Anime.Data.Type = entityType
		animeData.Relationship.User.Data.Id = userId
		animeData.Relationship.User.Data.Type = "users"
		animeData.Type = "library-entries"

		wrapper := map[string]animeRequestData{"data": animeData}
		jsonData, err = json.Marshal(wrapper)
		if err != nil {
			return entity, fmt.Errorf("error json marshalling request data: %v", err)
		}

	} else {
		mangaData.Attribute.Status = entityStatus
		mangaData.Relationship.Manga.Data.Id = id
		mangaData.Relationship.Manga.Data.Type = entityType
		mangaData.Relationship.User.Data.Id = userId
		mangaData.Relationship.User.Data.Type = "users"
		mangaData.Type = "library-entries"

		wrapper := map[string]mangaRequestData{"data": mangaData}
		jsonData, err = json.Marshal(wrapper)
		if err != nil {
			return entity, fmt.Errorf("error json marshalling request data: %v", err)
		}
	}

	requestBody := bytes.NewReader(jsonData)
	request, err := http.NewRequest("POST", "https://kitsu.io/api/edge/library-entries", requestBody)
	if err != nil {
		return entity, fmt.Errorf("error creating new user request: %v", err)
	}
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")
	request.Header.Set("Authorization", "Bearer "+utils.AuthToken.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return entity, fmt.Errorf("error sending the create_library_entry request: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return entity, fmt.Errorf("error reading the create_library_entry response: %v", err)
	}
	if response.StatusCode != http.StatusCreated {
		fmt.Println(string(responseBody), response.Status)
		return entity, fmt.Errorf("server error: %v", err)
	}

	entData := make(map[string]EntryData)
	err = json.Unmarshal([]byte(responseBody), &entData)
	if err != nil {
		return entity, fmt.Errorf("error unmarshalling the create_library_entry response: %v", err)
	}
	entity = entData["data"]

	userName, entryName, entryType, err := getDataByEntry(entity)
	if err != nil {
		fmt.Printf("Error getting the entity details: %v\n", err)
	}
	entity.Attribute.Name = entryName
	entity.Attribute.Type = entryType
	entity.Attribute.User = userName

	return entity, err
}

func GetLibraryEntry(userid string) ([]EntryData, error) {
	var entities []EntryData
	user, err := users.GetUser(userid)
	if err != nil {
		return entities, err
	}
	path := user.Relationship.(map[string]interface{})["libraryEntries"].(map[string]interface{})["links"].(map[string]interface{})["related"].(string)

	response, err := http.Get(path)
	if err != nil {
		return entities, fmt.Errorf("error requesting at %v: %v", path, err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return entities, fmt.Errorf("error reading the get_library response: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(string(responseBody), response.Status)
		return entities, fmt.Errorf("server error: %v", err)
	}

	responseData := struct {
		Entities []EntryData `json:"data"`
	}{}
	err = json.Unmarshal([]byte(responseBody), &responseData)
	if err != nil {
		return entities, fmt.Errorf("error unmarshalling the get_library response: %v", err)
	}
	entities = responseData.Entities

	for i, entity := range entities {
		userName, entryName, entryType, err := getDataByEntry(entity)
		if err != nil {
			fmt.Printf("Error getting the entity details: %v\n", err)
		}
		entities[i].Attribute.Name = entryName
		entities[i].Attribute.Type = entryType
		entities[i].Attribute.User = userName
	}

	return entities, err
}

func DeleteLibraryEntry() {}
