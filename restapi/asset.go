package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Asset returns an asset from the Vimond Rest API
func (c *Client) Asset(ctx context.Context, platform, assetID string) (*Asset, error) {
	if _, err := strconv.Atoi(assetID); err != nil {
		return nil, ErrInvalidAssetID
	}

	resp, err := c.get(ctx, c.assetPath(platform, assetID), url.Values{"expand": {"metadata,category"}})
	if err != nil {
		return nil, err
	}
	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 64)
		resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, ErrUnknown
	}

	return parseAsset(resp.Body)
}

// AssetRaw returns the raw response for an asset from the Vimond Rest API.
// possible values for headerAccept:
// application/json; v=3; charset=utf-8
// application/json; v=2; charset=utf-8
// application/json; charset=utf-8
// application/xml; charset=utf-8
func (c *Client) AssetRaw(ctx context.Context, platform, assetID, headerAccept string) ([]byte, error) {
	if headerAccept != "" {
		c.headerAccept = headerAccept
	}

	resp, err := c.get(ctx, c.assetPath(platform, assetID), url.Values{"expand": {"metadata,category"}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Vimond replies with status %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	return b, nil
}

func (c *Client) assetPath(platform, assetID string) string {
	return fmt.Sprintf("/api/%s/asset/%s", platform, assetID)
}

func parseAsset(r io.Reader) (*Asset, error) {
	type Alias Asset

	var asset struct {
		AssetTypeID int     `json:"assetTypeId"`
		CategoryID  int     `json:"categoryId"`
		ChannelID   int     `json:"channelId"`
		Duration    float32 `json:"duration"`
		ID          int     `json:"id"`

		*Alias
	}

	if err := json.NewDecoder(r).Decode(&asset); err != nil {
		return nil, err
	}

	asset.Alias.AssetTypeID = strconv.Itoa(asset.AssetTypeID)
	asset.Alias.CategoryID = strconv.Itoa(asset.CategoryID)
	asset.Alias.ChannelID = strconv.Itoa(asset.ChannelID)
	asset.Alias.Duration = int(asset.Duration)
	asset.Alias.ID = strconv.Itoa(asset.ID)

	return (*Asset)(asset.Alias), nil
}

// Asset is a Vimond Rest API asset
type Asset struct {
	ID         string `json:"id"`
	ChannelID  string `json:"channelId"`
	CategoryID string `json:"categoryId"`

	AssetTypeID string `json:"assetTypeId"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	Title       string `json:"title"`

	ImageVersions ImageVersions `json:"imageVersions"`

	Archive        bool `json:"archive"`
	Aspect16x9     bool `json:"aspect16x9"`
	AutoDistribute bool `json:"autoDistribute"`
	AutoEncode     bool `json:"autoEncode"`
	AutoPublish    bool `json:"autoPublish"`
	CopyLiveStream bool `json:"copyLiveStream"`
	DRMProtected   bool `json:"drmProtected"`
	Deleted        bool `json:"deleted"`
	ItemsPublished bool `json:"itemsPublished"`
	LabeledAsFree  bool `json:"labeledAsFree"`
	Live           bool `json:"live"`

	Duration int `json:"duration"`
	Views    int `json:"views"`

	AccurateDuration float64 `json:"accurateDuration"`

	CreateTime        time.Time `json:"createTime"`
	ExpireDate        time.Time `json:"expireDate"`
	LiveBroadcastTime time.Time `json:"liveBroadcastTime"`
	UpdateTime        time.Time `json:"updateTime"`

	Metadata AssetMetadata `json:"metadata"`
	Category Category      `json:"category"`
}

// ImageVersions is a slice of ImageVersion
type ImageVersions struct {
	Images []Image `json:"images,omitempty"`
}

// TypeURL returns the URL for the given image version type
func (ivs ImageVersions) TypeURL(ivt string) string {
	for _, iv := range ivs.Images {
		if iv.Type == ivt {
			return iv.URL
		}
	}

	return ""
}

// Image is a version of image representing the asset
type Image struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Category is a category node in the Vimond Rest API category tree
type Category struct {
	Parent *Category `json:"parent"`
	Title  string    `json:"title"`
	ID     string    `json:"id"`
}

// In walks the category tree upwards, looking for the given ID
func (c *Category) In(id string) bool {
	if c.ID == id {
		return true
	}

	if c.Parent != nil {
		return c.Parent.In(id)
	}

	return false
}

// AssetMetadata represents Vimond asset metadata.
type AssetMetadata struct {
	Entries MetadataEntries `json:"entries"`
}

// MetadataEntries is metadata for an Asset in the Vimond Rest API
type MetadataEntries struct {
	Annotags               LocalizedField `json:"annotags,omitempty"`
	AssetLength            LocalizedField `json:"asset-length,omitempty"`
	ContentAPIID           LocalizedField `json:"content-api-id,omitempty"`
	ContentAPISeasonID     LocalizedField `json:"content-api-season-id,omitempty"`
	ContentAPISeriesID     LocalizedField `json:"content-api-series-id,omitempty"`
	ContentSource          LocalizedField `json:"content-source,omitempty"`
	DescriptionShort       LocalizedField `json:"description-short,omitempty"`
	Episode                LocalizedField `json:"episode,omitempty"`
	Genre                  LocalizedField `json:"genre,omitempty"`
	GenreDescription       LocalizedField `json:"genre-description,omitempty"`
	HideAds                LocalizedField `json:"hideAds,omitempty"`
	JuneMediaID            LocalizedField `json:"june-media-id,omitempty"`
	JuneProgramID          LocalizedField `json:"june-program-id,omitempty"`
	LouisePressTitle       LocalizedField `json:"louise-press-title,omitempty"`
	LouiseProductKey       LocalizedField `json:"louise-product-key,omitempty"`
	Season                 LocalizedField `json:"season,omitempty"`
	SeasonID               LocalizedField `json:"season-id,omitempty"`
	SeasonSynopsis         LocalizedField `json:"season-synopsis,omitempty"`
	SeriesDescriptionShort LocalizedField `json:"series-description-short,omitempty"`
	SeriesID               LocalizedField `json:"series-id,omitempty"`
	Title                  LocalizedField `json:"title,omitempty"`
	YouTubeTemplate        LocalizedField `json:"youtube-template,omitempty"`
}

// LocalizedField is field with localized values
type LocalizedField []LocalizedValue

// Value returns the value for the given lang, fallback to *
func (lf LocalizedField) Value(lang string) string {
	def := LocalizedValue{}

	for _, l := range lf {
		if l.Lang == lang {
			return l.Value
		} else if l.Lang == "*" {
			def = l
		}
	}

	return def.Value
}

// LocalizedValue is a representation of a multi-language value
type LocalizedValue struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}
