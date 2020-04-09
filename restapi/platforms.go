package restapi

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Platform represents a Vimond platform, such as TV4.
type Platform struct {
	ID   int
	Name string
}

// Platforms returns the list of available platforms.
func (c *Client) Platforms(ctx context.Context) ([]Platform, error) {
	path := "/api/admin/platforms"

	resp, err := c.get(ctx, path, url.Values{})
	if err != nil {
		return nil, err
	}
	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 64)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUnknown
	}

	return parsePlatforms(resp.Body)
}

func parsePlatforms(r io.Reader) ([]Platform, error) {
	var resp []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, err
	}

	platforms := make([]Platform, 0, len(resp))

	for n := range resp {
		platforms = append(platforms, Platform{
			ID:   resp[n].ID,
			Name: resp[n].Name,
		})
	}

	return platforms, nil
}
