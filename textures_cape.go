package minecraft

import _ "image/png" // If we work with PNGs we need this

type Cape struct {
	Texture
}

func FetchCapeUUID(uuid string) (Cape, error) {
	cape := &Cape{}

	// Must be careful to not request same profile from session server more than once per ~30 seconds
	sessionProfile, err := GetSessionProfile(uuid)
	if err != nil {
		return *cape, err
	}

	return *cape, cape.FetchWithSessionProfile(sessionProfile, "Cape")
}

func FetchCapeUsernameMojang(username string) (Cape, error) {
	cape := &Cape{}

	return *cape, cape.FetchWithUsernameMojang(username, "Cape")
}

func FetchCapeUsernameS3(username string) (Cape, error) {
	cape := &Cape{}

	return *cape, cape.FetchWithUsernameS3(username, "Cape")
}
