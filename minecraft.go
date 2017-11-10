// Package minecraft is a library for interacting with the profiles and skins of Minecraft players
package minecraft

import (
	"io"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
)

const (
	// ValidUsernameRegex is proper Minecraft username regex
	ValidUsernameRegex = `[a-zA-Z0-9_]{1,16}`

	// ValidUUIDRegex is proper Minecraft UUID regex
	ValidUUIDRegex = `[0-9a-f]{32}|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`

	// ValidUsernameOrUUIDRegex is proper Minecraft Username-or-UUID regex
	ValidUsernameOrUUIDRegex = "(" + ValidUUIDRegex + "|" + ValidUsernameRegex + ")"
)

var (
	// RegexUsername is our compiled once Username matching Regex
	RegexUsername = regexp.MustCompile("^" + ValidUsernameRegex + "$")

	// RegexUUID is our compiled once UUID matching Regex
	RegexUUID = regexp.MustCompile("^" + ValidUUIDRegex + "$")

	// RegexUsernameOrUUID is our compiled once Username OR UUID matching Regex
	RegexUsernameOrUUID = regexp.MustCompile("^" + ValidUsernameOrUUIDRegex + "$")
)

// Mojang APIs have fairly standard responses and this makes those requests and
// catches the errors. Remember to close the response!
func apiRequest(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to request URL")
	}

	switch resp.StatusCode {

	case http.StatusOK:
		return resp.Body, nil

	case http.StatusNoContent:
		return resp.Body, errors.New("user not found")

	case http.StatusTooManyRequests:
		return resp.Body, errors.New("rate limited")

	case http.StatusForbidden:
		return resp.Body, errors.New("likely not found, maybe forbidden")

	default:
		return resp.Body, errors.Errorf("apiRequest HTTP %s", resp.Status)
	}
}
