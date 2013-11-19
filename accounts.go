package minecraft

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type User struct {
	Id   string
	Name string
}

type ProfileResposne struct {
	Size int
	User struct {
		Users []User
	}
}

func GetUser(username string) User {
	fUser := &User{Name: username}
	buf, _ := json.Marshal(fUser)
	body := bytes.NewBuffer(buf)

	resp, err := http.Post("https://api.mojang.com/profiles/page/1", "application/json", body)
	if err != nil {
		panic("POST failed")
	}
	response, _ := ioutil.ReadAll(resp.Body)

	proResponse := ProfileResposne{}
	jsonErr := json.Unmarshal(response, &proResponse)
	if jsonErr != nil {
		panic("Invalid JSON response")
	}

	return proResponse.User.Users[0]
}
