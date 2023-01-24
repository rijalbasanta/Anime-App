package main

import (
	"anime-app/users"
	"anime-app/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
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
	args := os.Args[1:]
	switch args[0] {

	// cleanup
	case "flush":
		var data []map[string]struct {
			Id    string `json:"id"`
			Email string `json:"email"`
		}
		fileData, err := ioutil.ReadFile("./users/users.json")
		if err != nil {
			fmt.Printf("Can't open the users file: %v\n", err)
		}
		err = json.Unmarshal(fileData, &data)
		if err != nil {
			fmt.Printf("Can't read the users file data: %v\n", err)
		}
		for i, user := range data {
			key := "user" + strconv.Itoa(i)
			err1 := utils.Authenticate(user[key].Email)
			if err1 != nil {
				fmt.Printf("Error authenticating the user %v: %v\n", user[key].Id, err1)
			}
			err1 = users.DeleteUser(user[key].Id)
			if err1 != nil {
				fmt.Printf("Error deleting the user %v: %v\n", user[key].Id, err1)
			} else {
				fmt.Printf("User %v deleted.\n", user[key].Id)
			}

			time.Sleep(3 * time.Second)
		}
		err = os.Remove("./users/users.json")
		if err != nil {
			fmt.Printf("Error deleting the users file: %v\n", err)
		}

	// add entities
	case "add":

		// add account
		if args[1] == "account" {
			for _, name := range args[2:] {
				user, err := users.CreateUser(name)
				if err != nil {
					fmt.Printf("Error creating the user %v: %v\n", name, err)
				}
				fmt.Printf("Id: %v\t\t Name: %v\t\t Email: %v\n", user.Id, user.Attribute.UserName, user.Attribute.Email)

				time.Sleep(3 * time.Second)
			}
			fmt.Println("Completed")
			return
		}

	// read entities
	case "get":

		// read account
		if args[1] == "account" {
			// var accounts []users.UserData
			for _, id := range args[2:] {
				user, err := users.GetUser(id)
				if err != nil {
					fmt.Printf("Error getting the user %v: %v\n", id, err)
				}
				// accounts = append(accounts, user)
				fmt.Printf("Id: %v\t\t Name: %v\t\t Email: %v\n", user.Id, user.Attribute.UserName, user.Attribute.Email)

				time.Sleep(3 * time.Second)
			}
			fmt.Println("Completed")
			return
		}

	// delete entities
	case "delete":

		// delete account
		if args[1] == "account" {
			for _, id := range args[2:] {
				user, err := users.GetUser(id)
				if err != nil {
					fmt.Printf("Error getting the user %v for authentication: %v\n", id, err)
				}
				err = utils.Authenticate(user.Attribute.Email)
				if err != nil {
					fmt.Printf("Error authenticating the user %v: %v\n", id, err)
				}
				err = users.DeleteUser(id)
				if err != nil {
					fmt.Printf("Error deleting the user %v: %v\n", id, err)
				}

				time.Sleep(3 * time.Second)
			}
			fmt.Println("Completed")
			return
		}
	}
}

// func main() {
// 	// name := os.Args[1]
// 	// user, err := users.CreateUser(name)
// 	user, err := users.GetUser("1386365")
// 	if err != nil {
// 		fmt.Println("Error getting user: ", err)
// 	}
// 	fmt.Println(user.Id, user.Attribute.UserName, user.Attribute.Email, user.Relationship.(map[string]interface{})["libraryEntries"].(map[string]interface{})["links"])

// 	err = utils.Authenticate(user.Attribute.Email)
// 	if err != nil {
// 		fmt.Println("Error authenticating the user: %w", err)
// 	}
// 	fmt.Println("Authenticated")
// 	err = users.DeleteUser(user.Id)
// 	if err != nil {
// 		fmt.Println("Error deleting the user: %w", err)
// 	}
// 	fmt.Println("Deleted")
// 	user, err = users.GetUser("1386365")
// 	if err != nil {
// 		fmt.Println("Error getting user: ", err)
// 	}

// }
