// Minecraft Textures
package minecraft

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	_ "image/png"
	"strings"
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

// FetchTextures is our silver bullet/wrapper function to get both a Skin and Cape with a
// single request to the Session Servers and return the User details or fall back if have to
func FetchTextures(player string) (User, Skin, Cape, error) {
	user := User{}
	skin := Skin{}
	cape := Cape{}

	// This function takes a Username/UUID and returns a UUID formatted without
	// hyphens. Will error if it can't lookup Username (no account or API error) or
	// if the Username/UUID regex fails
	uuid, uuidErr := NormalizePlayerForUUID(player)
	if uuidErr == nil {
		// We should have a valid UUID which *hopefully* corresponds to a user.
		// Must be careful to not request same profile from session server more than once per ~30 seconds
		sessionProfile, err := GetSessionProfile(uuid)
		if err == nil {
			// We got a sessionProfile (so UUID must be for a user)
			user, skin, cape, err = FetchTexturesWithSessionProfile(sessionProfile)
			if err == nil {
				// We got the skin and cape!
				return user, skin, cape, nil
			} else if strings.HasPrefix(err.Error(), "FetchTexturesWithSessionProfile failed: Unable to retrieve cape") {
				// User likely has no cape - no worries :)
				return user, skin, cape, nil
			}
			// Every other error means that we don't have a skin :(
		}
	}

	if (uuidErr != nil && !strings.HasPrefix(uuidErr.Error(), "GetAPIProfile failed: (apiRequest failed: User not found") &&
		uuidErr.Error() != "NormalizePlayerForUUID failed: Invalid Username or UUID.") ||
		(uuidErr == nil && IsUsername(player)) {
		// If the uuidErr is *not* "User not found" or related to a non username/UUID we
		// know that there must have been another issue (and that the user potentially
		// still exists... somewhere?)
		// Or there was no uuidErr but we have a Username and it broke later on...

		// Let's try another request to the Mojang - this might bypass a rate limit?
		skin, err := FetchSkinUsernameMojang(player)
		if err == nil {
			user.Username = player
			return user, skin, cape, errors.New("FetchTextures fallback to UsernameMojang: (" + uuidErr.Error() + ")")
		}

		skin, err = FetchSkinUsernameS3(player)
		if err == nil {
			user.Username = player
			return user, skin, cape, errors.New("FetchTextures fallback to UsernameS3: (" + uuidErr.Error() + ")")
		}
	}

	// Username or UUID not real or some other issue which we won't recover from.
	// Steve to the rescue!
	skin, err := FetchSkinForSteve()
	if err == nil {
		err = errors.New("FetchTextures fallback to Steve: (" + uuidErr.Error() + ")")
	} else {
		err = errors.New("FetchTextures failed to fallback:  (" + err.Error() + ") (" + uuidErr.Error() + ")")
	}
	return user, skin, cape, err
}

func FetchTexturesWithSessionProfile(sessionProfile SessionProfileResponse) (User, Skin, Cape, error) {
	//  We have a sessionProfile!
	user := &User{UUID: sessionProfile.UUID, Username: sessionProfile.Username}
	skin := &Skin{}
	cape := &Cape{}

	profileTextureProperty, err := DecodeTextureProperty(sessionProfile)
	if err != nil {
		return *user, *skin, *cape, errors.New("FetchTexturesWithSessionProfile failed: Unable to decode sessionProfile (" + err.Error() + ")")
	}

	// We got oursleves a profileTextureProperty - now we can get a Skin/Cape

	err = skin.FetchWithTextureProperty(profileTextureProperty, "Skin")
	if err != nil {
		return *user, *skin, *cape, errors.New("FetchTexturesWithSessionProfile failed: Unable to retrieve skin - (" + err.Error() + ")")
	}

	err = cape.FetchWithTextureProperty(profileTextureProperty, "Cape")
	if err != nil {
		return *user, *skin, *cape, errors.New("FetchTexturesWithSessionProfile failed: Unable to retrieve cape - (" + err.Error() + ")")
	}
	return *user, *skin, *cape, nil
}

func DecodeTextureProperty(sessionProfile SessionProfileResponse) (SessionProfileTextureProperty, error) {
	var texturesProperty *SessionProfileProperty
	for _, v := range sessionProfile.Properties {
		if v.Name == "textures" {
			texturesProperty = &v
			break
		}
	}

	if texturesProperty == nil {
		return SessionProfileTextureProperty{}, errors.New("DecodeTextureProperty failed: No textures property.")
	}

	profileTextureProperty := SessionProfileTextureProperty{}
	err := json.NewDecoder(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(texturesProperty.Value))).Decode(&profileTextureProperty)
	if err != nil {
		return SessionProfileTextureProperty{}, errors.New("DecodeTextureProperty failed: Error decoding texture property - (" + err.Error() + ")")
	}

	return profileTextureProperty, nil
}

// DecodeTextureURL will return a texture URL string when supplied with a
// SessionProfileTextureProperty and a type (Skin|Cape).
func DecodeTextureURL(profileTextureProperty SessionProfileTextureProperty, textureType string) (string, error) {
	textureURL := profileTextureProperty.Textures.Skin.URL
	if textureType != "Skin" {
		textureURL = profileTextureProperty.Textures.Cape.URL
	}

	if textureURL == "" {
		return "", errors.New("DecodeTextureURL failed: " + textureType + " URL is not present.")
	}

	return textureURL, nil
}
