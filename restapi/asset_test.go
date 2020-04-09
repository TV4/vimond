package restapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func TestAsset(t *testing.T) {
	t.Run("invalid_asset_id", func(t *testing.T) {
		c := &Client{}

		if _, err := c.Asset(context.Background(), "tv4", "invalid"); err != ErrInvalidAssetID {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("valid_assets", func(t *testing.T) {
		for _, tc := range []struct {
			in  vimondAsset
			out Asset
		}{
			{
				in: vimondAsset{
					AssetTypeID: 10001,
					CategoryID:  10002,
					ChannelID:   10003,
					Duration:    10004.5,
					ID:          10006,
					Title:       "title foo",
				},
				out: Asset{
					AssetTypeID: "10001",
					CategoryID:  "10002",
					ChannelID:   "10003",
					Duration:    10004,
					ID:          "10006",
					Title:       "title foo",
				},
			},
			{
				in: vimondAsset{
					AssetTypeID: 20001,
					CategoryID:  20002,
					ChannelID:   20003,
					Duration:    20004.5,
					ID:          20006,
					Metadata:    vimondAssetMetadata{LouisePressTitle: "louise press title bar"},
					Title:       "asset title bar",
				},
				out: Asset{
					AssetTypeID: "20001",
					CategoryID:  "20002",
					ChannelID:   "20003",
					Duration:    20004,
					ID:          "20006",
					Title:       "asset title bar",
					Metadata:    AssetMetadata{LouisePressTitle: "louise press title bar"},
				},
			},
		} {
			ts, c := testServerAndClient(func(w http.ResponseWriter, r *http.Request) {
				b, _ := json.Marshal(tc.in)
				w.Write(b)
			})
			defer ts.Close()

			t.Run(strconv.Itoa(tc.in.ID), func(t *testing.T) {
				asset, err := c.Asset(context.Background(), "tv4", strconv.Itoa(tc.in.ID))
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if got, want := asset.ID, tc.out.ID; got != want {
					t.Errorf("asset.ID = %q, want %q", got, want)
				}

				if got, want := asset.CategoryID, tc.out.CategoryID; got != want {
					t.Errorf("asset.CategoryID = %q, want %q", got, want)
				}

				if got, want := asset.Title, tc.out.Title; got != want {
					t.Errorf("asset.Title = %q, want %q", got, want)
				}

				if got, want := asset.Metadata.LouisePressTitle, tc.out.Metadata.LouisePressTitle; got != want {
					t.Errorf("a.Metadata = %q, want %q", got, want)
				}
			})
		}
	})
}

func TestImageVersions(t *testing.T) {
	t.Run("TypeURL", func(t *testing.T) {
		for _, tt := range []struct {
			ivs  ImageVersions
			ivt  string
			want string
		}{
			{ImageVersions{}, "", ""},
			{ImageVersions{}, "original", ""},
			{ImageVersions{Images: []Image{{"original", "foo"}}}, "original", "foo"},
			{ImageVersions{Images: []Image{{"original", "foo"}, {"secondary", "bar"}}}, "secondary", "bar"},
		} {
			if got := tt.ivs.TypeURL(tt.ivt); got != tt.want {
				t.Fatalf("tt.ivs.TypeURL(%q) = %q, want %q", tt.ivt, got, tt.want)
			}
		}
	})
}

func TestCategoryIn(t *testing.T) {
	for _, tt := range []struct {
		name     string
		id       string
		category Category
		want     bool
	}{
		{"no_parent_true", "0", Category{ID: "0"}, true},
		{"no_parent_false", "1", Category{ID: "0"}, false},
		{"one_parent_true", "0", Category{ID: "1", Parent: &Category{ID: "0"}}, true},
		{"one_parent_false", "0", Category{ID: "2", Parent: &Category{ID: "1"}}, false},
		{"two_parents_true", "0", Category{ID: "2", Parent: &Category{ID: "1", Parent: &Category{ID: "0"}}}, true},
		{"two_parents_false", "0", Category{ID: "3", Parent: &Category{ID: "2", Parent: &Category{ID: "1"}}}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.In(tt.id); got != tt.want {
				t.Fatalf("tt.category.In(%q) = %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}

func TestLocalizedField(t *testing.T) {
	lf := LocalizedField{
		{"*", "Foo"},
		{"sv_SE", "Bar"},
	}

	if got, want := lf.Value("*"), "Foo"; got != want {
		t.Fatalf(`lf.Value("*") = %q, want %q`, got, want)
	}

	if got, want := lf.Value("da_DK"), "Foo"; got != want {
		t.Fatalf(`lf.Value("da_DK") = %q, want %q`, got, want)
	}

	if got, want := lf.Value("sv_SE"), "Bar"; got != want {
		t.Fatalf(`lf.Value("sv_SE") = %q, want %q`, got, want)
	}
}

type vimondAsset struct {
	AssetTypeID int                 `json:"assetTypeId"`
	CategoryID  int                 `json:"categoryId,attr"`
	ChannelID   int                 `json:"channelId"`
	Duration    float32             `json:"duration"`
	ID          int                 `json:"id,attr"`
	Metadata    vimondAssetMetadata `json:"metadata"`
	Title       string              `json:"title"`
}

type vimondAssetMetadata struct {
	LouisePressTitle string `json:"louisePressTitle,omitempty"`
}
