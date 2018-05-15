package restapi

import (
	"context"
	"encoding/json"
	"encoding/xml"
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

	resp, err := c.get(ctx, c.assetPath(platform, assetID), url.Values{
		"expand": {"metadata,category"},
	})
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

	var a Asset
	if err := xml.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, err
	}

	return &a, nil
}

func (c *Client) assetPath(platform, assetID string) string {
	return fmt.Sprintf("/api/%s/asset/%s", platform, assetID)
}

// Asset is a Vimond Rest API asset
type Asset struct {
	ID         string `xml:"id,attr" json:"id"`
	ChannelID  string `xml:"channelId,attr" json:"channel_id"`
	CategoryID string `xml:"categoryId,attr" json:"category_id"`

	AssetTypeID string `xml:"assetTypeId" json:"asset_type_id"`
	Description string `xml:"description" json:"description"`
	ImageURL    string `xml:"imageUrl" json:"image_url"`
	Title       string `xml:"title" json:"title"`

	ImageVersions []ImageVersion `xml:"imageVersions>image"`

	Archive        bool `xml:"archive" json:"archive"`
	Aspect16x9     bool `xml:"aspect16x9" json:"aspect_16x9"`
	AutoDistribute bool `xml:"autoDistribute" json:"auto_distribute"`
	AutoEncode     bool `xml:"autoEncode" json:"auto_encode"`
	AutoPublish    bool `xml:"autoPublish" json:"auto_publish"`
	CopyLiveStream bool `xml:"copyLiveStream" json:"copy_live_stream"`
	DRMProtected   bool `xml:"drmProtected" json:"drm_protected"`
	Deleted        bool `xml:"deleted" json:"deleted"`
	ItemsPublished bool `xml:"itemsPublished" json:"items_published"`
	LabeledAsFree  bool `xml:"labeledAsFree" json:"labeled_as_free"`
	Live           bool `xml:"live" json:"live"`

	Duration int `xml:"duration" json:"duration"`
	Views    int `xml:"views" json:"views"`

	AccurateDuration float64 `xml:"accurateDuration" json:"accurate_duration"`

	CreateTime        time.Time `xml:"createTime" json:"create_time"`
	ExpireDate        time.Time `xml:"expireDate" json:"expire_date"`
	LiveBroadcastTime time.Time `xml:"liveBroadcastTime" json:"live_broadcast_time"`
	UpdateTime        time.Time `xml:"updateTime" json:"update_time"`

	Metadata AssetMetadata `xml:"metadata" json:"metadata"`
	Category Category      `xml:"category" json:"category"`
}

// ImageVerion is a version of image representing the asset
type ImageVersion struct {
	Type string `xml:"type,attr" json:"type"`
	URL  string `xml:"url" json:"url"`
}

// Category is a category node in the Vimond Rest API category tree
type Category struct {
	Parent *Category `xml:"parent" json:"parent"`
	Title  string    `xml:"title" json:"title"`
	ID     string    `xml:"id,attr" json:"id"`
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

// AssetMetadata is metadata for an Asset in the Vimond Rest API
type AssetMetadata struct {
	ContentSource          string         `xml:"content-source" json:"content_source,omitempty"`
	Genre                  string         `xml:"genre" json:"genre,omitempty"`
	LouisePressTitle       string         `xml:"louise-press-title" json:"louise_press_title,omitempty"`
	LouiseProductKey       string         `xml:"louise-product-key" json:"louise_product_key,omitempty"`
	LouiseProgramType      string         `xml:"louise-program-type" json:"louise_program_type,omitempty"`
	Episode                json.Number    `xml:"episode" json:"episode,omitempty"`
	Season                 json.Number    `xml:"season" json:"season,omitempty"`
	SeasonSynopsis         LocalizedField `xml:"season-synopsis" json:"season_synopsis,omitempty"`
	SeriesDescriptionShort LocalizedField `xml:"series-description-short" json:"series_description_short,omitempty"`
	GenreDescription       LocalizedField `xml:"genre-description" json:"genre_description,omitempty"`
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

// LocalizedValue is a representation of a parsed multi-language value
type LocalizedValue struct {
	Lang  string `xml:"lang,attr" json:"lang"`
	Value string `xml:",chardata" json:"value"`
}

type publishing struct {
	Platform string    `xml:"publishing>platform"`
	Publish  time.Time `xml:"publishing>publish"`
}
