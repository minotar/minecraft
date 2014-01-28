package minecraft

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Id   string
	Name string
}

type ProfileResponse struct {
	Size     int
	Profiles []struct {
		User
	}
}

func GetUser(username string) (User, error) {
	postBody := []byte(`{"agent":"Minecraft","name":"` + username + `"}`)
	body := bytes.NewBuffer(postBody)

	r, httpErr := http.Post("https://api.mojang.com/profiles/page/1", "application/json", body)
	if httpErr != nil {
		log.Fatalln(httpErr)
	}
	defer r.Body.Close()

	response, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatalln(readErr)
	}

	proResponse := ProfileResponse{}
	if err := json.Unmarshal(response, &proResponse); err != nil {
		log.Fatalln(err)
	}

	if len(proResponse.Profiles) == 0 {
		return User{}, errors.New("User not found.")
	}

	return proResponse.Profiles[0].User, nil
}
