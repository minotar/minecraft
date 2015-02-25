// Minecraft Textures
package minecraft

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	_ "image/png"
	"io"
	"net/http"
)

type SessionProfileTextureProperty struct {
	TimestampMs uint64 `json:"timestamp"`
	ProfileUUID string `json:"profileId"`
	ProfileName string `json:"profileName"`
	IsPublic    bool   `json:"isPublic"`
	Textures    struct {
		Skin struct {
			URL string `json:"url"`
		} `json:"SKIN"`
		Cape struct {
			URL string `json:"url"`
		} `json:"CAPE"`
	} `json:"textures"`
}

func decodeTextureProperty(sessionProfile SessionProfileResponse) (SessionProfileTextureProperty, error) {
	var texturesProperty *SessionProfileProperty
	for _, v := range sessionProfile.Properties {
		if v.Name == "textures" {
			texturesProperty = &v
			break
		}
	}

	if texturesProperty == nil {
		return SessionProfileTextureProperty{}, errors.New("No textures property")
	}

	profileTextureProperty := SessionProfileTextureProperty{}
	err := json.NewDecoder(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(texturesProperty.Value))).Decode(&profileTextureProperty)
	if err != nil {
		return SessionProfileTextureProperty{}, errors.New("Error decoding texture property. (" + err.Error() + ")")
	}

	return profileTextureProperty, nil
}

// decodeTextureURL will return a texture URL string when supplied with a UUID
// and a type (Skin|Cape).
func decodeTextureURL(uuid string, textureType string) (string, error) {
	sessionProfile, err := GetSessionProfile(uuid)
	if err != nil {
		return "", err
	}

	profileTextureProperty, err := decodeTextureProperty(sessionProfile)
	if err != nil {
		return "", errors.New("Error decoding texture property for \"" + uuid + "\" (" + err.Error() + ")")
	}

	textureURL := profileTextureProperty.Textures.Skin.URL
	if textureType != "Skin" {
		textureURL = profileTextureProperty.Textures.Cape.URL
	}

	return textureURL, nil
}

// fetchTexture will return a Response.Body when supplied with texture URL.
// Remember to close the response!
func fetchTexture(textureURL string) (io.ReadCloser, error) {
	resp, err := http.Get(textureURL)
	if err != nil {
		return resp.Body, errors.New("Error retrieving texture (" + err.Error() + ")")
	}

	if resp.StatusCode != http.StatusOK {
		return resp.Body, errors.New("Error retrieving texture (HTTP " + resp.Status + ")")
	}

	return resp.Body, nil
}
