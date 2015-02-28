// minecraft project minecraft.go
package minecraft

import (
	"errors"
	"io"
	"net/http"
)

const (
	// Proper Minecraft username regex
	ValidUsernameRegex = `[a-zA-Z0-9_]{1,16}`

	// Proper Minecraft UUID regex
	ValidUuidRegex = `[0-9a-f]{32}|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`

	// Minecraft username-or-UUID regex
	ValidUsernameOrUuidRegex = "(" + ValidUuidRegex + "|" + ValidUsernameRegex + ")"
)

// Mojang APIs have fairly standard responses and this makes those requests and
// catches the errors. Remember to close the response!
func apiRequest(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("apiRequest failed: Unable to Get URL - (" + err.Error() + ")")
	}

	if resp.StatusCode == http.StatusNoContent {
		return resp.Body, errors.New("apiRequest failed: User not found - (HTTP 204 No Content)")
	} else if resp.StatusCode == 429 { // StatusTooManyRequests
		return resp.Body, errors.New("apiRequest failed: Rate limited")
	} else if resp.StatusCode != http.StatusOK {
		return resp.Body, errors.New("apiRequest failed: Error retrieving profile - (HTTP " + resp.Status + ")")
	}

	return resp.Body, nil
}
