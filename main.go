package main

import (
	"anime-app/users"
	"fmt"
	"log"
	"os"
)

// type Title struct {
// 	Name string `json:"en"`
// }

type Attribute struct {
	Title string `json:"canonicalTitle"`
}

type Data struct {
}

type JSN struct {
	Data []struct {
		Attributes struct {
			Title string `json:"canonicalTitle"`
		} `json:"attributes"`

		Relationships struct {
			Episodes struct {
				Links struct {
				}
			}
		}
	} `json:"data"`
}

func main() {
	// client := &http.Client{}

	name := os.Args[1]
	user, err := users.CreateUser(name)
	// user, err := users.GetUser("1386335")
	if err != nil {
		log.Fatal("Error creating new user: ", err)
	}

	// // req, err := http.NewRequest("GET", "https://kitsu.io/api/edge/anime", nil)
	// req, email, err := users.CreateUser(name)
	// if err != nil {
	// 	log.Fatal("Error creating the request:", err)
	// }
	// response, err := client.Do(req)
	// if err != nil {
	// 	log.Fatal("Error sending the request:", err)
	// }
	// defer response.Body.Close()
	// resp_body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	log.Fatal("Error parsing the response:", err)
	// }
	// if response.Status != "201 Created" {
	// 	fmt.Println(string(resp_body), response.Status)
	// 	log.Fatal("Server Error:")
	// }
	// data := make(map[string]users.UserData)
	// err = json.Unmarshal([]byte(resp_body), &data)
	// if err != nil {
	// 	log.Fatal("Error unmarshalling the response:", err)
	// }
	// user := data["data"]
	// user.Attribute.Email = email
	// file, err := json.MarshalIndent(user, "", "\t")
	// if err != nil {
	// 	fmt.Println("Can't marshal the data: %w", err)
	// }
	// err = ioutil.WriteFile("user.json", file, 0644)
	// if err != nil {
	// 	fmt.Println("Can't write to the: %w", err)
	// }
	// fmt.Println(string(resp_body))
	fmt.Println(user.Id, user.Attribute.UserName, user.Attribute.Email, user.Relationship.(map[string]interface{})["libraryEntries"].(map[string]interface{})["links"])
}
