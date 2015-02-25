package minecraft

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type APIProfileResponse struct {
	UUID     string `json:"id"`
	Username string `json:"name"`
	Legacy   bool   `json:"legacy"`
	Demo     bool   `json:"demo"`
}

type SessionProfileResponse struct {
	UUID       string                   `json:"id"`
	Username   string                   `json:"name"`
	Properties []SessionProfileProperty `json:"properties"`
}

type SessionProfileProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetAPIProfile returns the API profile for a given username primarily of use
// for getting the UUID, but can also correct the capitilzation of a username or
// possibly get the account status (legacy or demo) - only included when true
func GetAPIProfile(username string) (APIProfileResponse, error) {
	url := "https://api.mojang.com/users/profiles/minecraft/"
	url += username

	apiBody, err := apiRequest(url)
	defer apiBody.Close()
	if err != nil {
		return APIProfileResponse{}, err
	}

	apiProfile := APIProfileResponse{}
	err = json.NewDecoder(apiBody).Decode(&apiProfile)
	if err != nil {
		return APIProfileResponse{}, errors.New("Error decoding profile. (" + err.Error() + ")")
	}

	return apiProfile, nil
}

// GetUUID returns the UUID for a given username (shorthand for GetAPIProfile)
func GetUUID(username string) (string, error) {
	apiProfile, err := GetAPIProfile(username)
	return apiProfile.UUID, err
}

// GetSessionProfile fetches the session profile of the UUID, this includes
// extra properties for the user (currently just a textures property)
func GetSessionProfile(uuid string) (SessionProfileResponse, error) {
	url := "https://sessionserver.mojang.com/session/minecraft/profile/"
	url += uuid

	apiBody, err := apiRequest(url)
	defer apiBody.Close()
	if err != nil {
		return SessionProfileResponse{}, err
	}

	sessionProfile := SessionProfileResponse{}
	err = json.NewDecoder(apiBody).Decode(&sessionProfile)
	if err != nil {
		return SessionProfileResponse{}, errors.New("Error decoding profile. (" + err.Error() + ")")
	}

	return sessionProfile, nil
}

// Mojang APIs have fairly standard responses and this makes those requests and
// catches the errors. Remember to close the response!
func apiRequest(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return resp.Body, errors.New("User not found. (HTTP 204 No Content)")
	} else if resp.StatusCode == 429 { // StatusTooManyRequests
		return resp.Body, errors.New("Rate limited")
	} else if resp.StatusCode != http.StatusOK {
		return resp.Body, errors.New("Error retrieving profile. (HTTP " + resp.Status + ")")
	}

	return resp.Body, nil
}
