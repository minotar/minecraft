package minecraft

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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
		panic(httpErr)
	}
	response, _ := ioutil.ReadAll(r.Body)

	proResponse := ProfileResponse{}
	if err := json.Unmarshal(response, &proResponse); err != nil {
		panic(err)
	}

	if len(proResponse.Profiles) == 0 {
		return User{}, errors.New("User not found.")
	}

	return proResponse.Profiles[0].User, nil
}
