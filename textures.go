package minecraft

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
	// If we work with PNGs we need this
	_ "image/png"

	"github.com/pkg/errors"
)

type SessionProfileTextureProperty struct {
	TimestampMs uint64 `json:"timestamp"`
	ProfileUUID string `json:"profileId"`
	ProfileName string `json:"profileName"`
	IsPublic    bool   `json:"isPublic"`
	Textures    struct {
		Skin struct {
			Metadata struct {
				Model string `json:"model"`
			} `json:"metadata"`
			URL string `json:"url"`
		} `json:"SKIN"`
		Cape struct {
			URL string `json:"url"`
		} `json:"CAPE"`
	} `json:"textures"`
}

// DecodeTextureProperty takes a SessionProfileResponse and breaks it down into the Skin/Cape URLs for downloading them
func DecodeTextureProperty(sessionProfile SessionProfileResponse) (SessionProfileTextureProperty, error) {
	var texturesProperty *SessionProfileProperty
	for _, v := range sessionProfile.Properties {
		if v.Name == "textures" {
			texturesProperty = &v
			break
		}
	}

	// Is the below tested? Will it defintely be nil !! ??
	if texturesProperty == nil {
		return SessionProfileTextureProperty{}, errors.New("no textures property")
	}

	profileTextureProperty := SessionProfileTextureProperty{}
	// Base64 decode the texturesProperty and further decode the JSON from it into profileTextureProperty
	err := json.NewDecoder(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(texturesProperty.Value))).Decode(&profileTextureProperty)
	if err != nil {
		return SessionProfileTextureProperty{}, errors.Wrap(err, "unable to DecodeTextureProperty")
	}

	return profileTextureProperty, nil
}

func FetchTexturesWithSessionProfile(sessionProfile SessionProfileResponse) (User, Skin, Cape, error) {
	//  We have a sessionProfile!
	user := &User{UUID: sessionProfile.UUID, Username: sessionProfile.Username}
	skin := &Skin{}
	cape := &Cape{}

	profileTextureProperty, err := DecodeTextureProperty(sessionProfile)
	if err != nil {
		return *user, *skin, *cape, errors.Wrap(err, "failed to decode sessionProfile")
	}

	// We got oursleves a profileTextureProperty - now we can get a Skin/Cape

	err = skin.FetchWithTextureProperty(profileTextureProperty, "Skin")
	if err != nil {
		return *user, *skin, *cape, errors.Wrap(err, "not able to retrieve skin")
	}

	err = cape.FetchWithTextureProperty(profileTextureProperty, "Cape")
	if err != nil {
		return *user, *skin, *cape, errors.Wrap(err, "not able to retrieve cape")
	}
	return *user, *skin, *cape, nil
}

// FetchTextures is our silver bullet/wrapper function to get both a Skin and Cape with a
// single request to the Session Servers and return the User details or fall back if have to
func FetchTextures(player string) (User, Skin, Cape, error) {
	user := User{}
	skin := Skin{}
	cape := Cape{}

	var catchErr error

	// NormalizePlayerForUUID takes a Username/UUID and returns a UUID formatted without
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
			} else if errors.Cause(err).Error() == "Cape URL not present" {
				// User has no cape - no worries :)
				return user, skin, cape, nil
			} else if strings.HasPrefix(err.Error(), "not able to retrieve cape") {
				// User likely has no cape - no worries :)
				return user, skin, cape, errors.Wrap(err, "unable to get the cape")
			}
			// Every other error means that we don't have a skin :(
		}
		catchErr = err
	} else {
		catchErr = uuidErr
	}

	if (uuidErr != nil && errors.Cause(uuidErr).Error() != "user not found" &&
		uuidErr.Error() != "unable to NormalizePlayerForUUID due to invalid Username/UUID") ||
		(uuidErr == nil && RegexUsername.MatchString(player)) {
		// If the uuidErr is *not* "User not found" or related to a non username/UUID we
		// know that there must have been another issue (and that the user potentially
		// still exists... somewhere?)
		// Or there was no uuidErr but we have a Username and it broke later on...

		// Let's try another request to the Mojang - this might bypass a rate limit?
		skin, err := FetchSkinUsernameMojang(player)
		if err == nil {
			user.Username = player
			return user, skin, cape, errors.Wrap(catchErr, "falling back to UsernameMojang")
		}

		// A last ditch effort to recall an old skin from S3
		skin, err = FetchSkinUsernameS3(player)
		if err == nil {
			user.Username = player
			return user, skin, cape, errors.Wrap(catchErr, "falling back to UsernameS3")
		}
	}

	// Username or UUID not real or some other issue which we won't recover from.
	// Steve to the rescue!
	skin, err := FetchSkinForSteve()
	if err == nil {
		err = errors.Wrap(catchErr, "falling back to Steve")
	} else {
		err = errors.Wrap(err, "failed to fallback to Steve")
	}
	return user, skin, cape, err
}
