package minecraft

import _ "image/png"

type Cape struct {
	Texture
}

func FetchCapeUUID(uuid string) (*Cape, error) {
	cape := &Cape{}

	// Must be careful to not request same profile from session server more than once per ~30 seconds
	sessionProfile, err := GetSessionProfile(uuid)
	if err != nil {
		return cape, err
	}

	return cape, cape.fetchWithSessionProfile(sessionProfile, "Cape")
}

func FetchCapeUsernameMojang(username string) (*Cape, error) {
	cape := &Cape{}

	return cape, cape.fetchWithUsernameMojang(username, "Cape")
}

func FetchCapeUsernameS3(username string) (*Cape, error) {
	cape := &Cape{}

	return cape, cape.fetchWithUsernameS3(username, "Cape")
}
