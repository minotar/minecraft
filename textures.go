// Minecraft Textures
package minecraft

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	_ "image/png"
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
		return SessionProfileTextureProperty{}, errors.New("decodeTextureProperty failed: No textures property.")
	}

	profileTextureProperty := SessionProfileTextureProperty{}
	err := json.NewDecoder(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(texturesProperty.Value))).Decode(&profileTextureProperty)
	if err != nil {
		return SessionProfileTextureProperty{}, errors.New("decodeTextureProperty failed: Error decoding texture property - (" + err.Error() + ")")
	}

	return profileTextureProperty, nil
}

// decodeTextureURL will return a texture URL string when supplied with a
// SessionProfileTextureProperty and a type (Skin|Cape).
func decodeTextureURL(profileTextureProperty SessionProfileTextureProperty, textureType string) (string, error) {
	textureURL := profileTextureProperty.Textures.Skin.URL
	if textureType != "Skin" {
		textureURL = profileTextureProperty.Textures.Cape.URL
	}

	if textureURL == "" {
		return "", errors.New("decodeTextureURL failed: " + textureType + " URL is not present.")
	}

	return textureURL, nil
}
