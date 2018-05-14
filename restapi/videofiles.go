package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (c *Client) Videofiles(ctx context.Context, assetID string) (*VideofilesResponse, error) {
	resp, err := c.get(ctx, c.videofilesPath(assetID), url.Values{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.CopyN(ioutil.Discard, resp.Body, 64)
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, ErrUnknown
	}

	var vr VideofilesResponse

	if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
		return nil, err
	}

	return &vr, nil
}

func (c *Client) videofilesPath(assetID string) string {
	return fmt.Sprintf("/api/admin/asset/%s/videofiles", assetID)
}

type VideofilesResponse struct {
	AssetID    int         `json:"assetId"`
	Title      string      `json:"title"`
	Videofiles []Videofile `json:"videofiles"`
}

type Videofile struct {
	Bitrate     int    `json:"bitrate"`
	MediaFormat string `json:"mediaFormat"`
	Scheme      string `json:"scheme"`
	Server      string `json:"server"`
	Base        string `json:"base"`
	URL         string `json:"url"`
}
