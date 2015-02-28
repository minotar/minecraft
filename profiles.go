package minecraft

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

type User struct {
	UUID     string `json:"id"`
	Username string `json:"name"`
}

type APIProfileResponse struct {
	User
	Legacy bool `json:"legacy"`
	Demo   bool `json:"demo"`
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

// NormalizePlayerForUUID takes either a Username or UUID formatted with or
// without hyphens and returns a uniform response (UUID)
func NormalizePlayerForUUID(player string) (string, error) {
	usernameRegex := regexp.MustCompile("^" + ValidUsernameRegex + "$")
	uuidRegex := regexp.MustCompile("^" + ValidUuidRegex + "$")

	if usernameRegex.MatchString(player) {
		return GetUUID(player)
	} else if uuidRegex.MatchString(player) == true {
		return strings.Replace(player, "-", "", 4), nil
	}

	return "", errors.New("NormalizePlayerForUUID failed: Invalid Username or UUID.")
}

// GetAPIProfile returns the API profile for a given username primarily of use
// for getting the UUID, but can also correct the capitilzation of a username or
// possibly get the account status (legacy or demo) - only included when true
func GetAPIProfile(username string) (APIProfileResponse, error) {
	url := "https://api.mojang.com/users/profiles/minecraft/"
	url += username

	apiBody, err := apiRequest(url)
	if apiBody != nil {
		defer apiBody.Close()
	}
	if err != nil {
		return APIProfileResponse{}, errors.New("GetAPIProfile failed: (" + err.Error() + ")")
	}

	apiProfile := APIProfileResponse{}
	err = json.NewDecoder(apiBody).Decode(&apiProfile)
	if err != nil {
		return APIProfileResponse{}, errors.New("GetAPIProfile failed: Error decoding profile - (" + err.Error() + ")")
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
	if apiBody != nil {
		defer apiBody.Close()
	}
	if err != nil {
		return SessionProfileResponse{}, errors.New("GetSessionProfile failed: (" + err.Error() + ")")
	}

	sessionProfile := SessionProfileResponse{}
	err = json.NewDecoder(apiBody).Decode(&sessionProfile)
	if err != nil {
		return SessionProfileResponse{}, errors.New("GetSessionProfile failed: Error decoding profile - (" + err.Error() + ")")
	}

	return sessionProfile, nil
}
