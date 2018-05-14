package restapi

import (
	"context"
	"encoding/xml"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestAsset(t *testing.T) {
	t.Run("invalid_asset_id", func(t *testing.T) {
		c := &Client{}

		if _, err := c.Asset(context.Background(), "tv4", "invalid"); err != ErrInvalidAssetID {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("valid_assets", func(t *testing.T) {
		for _, ta := range []testAsset{
			{ID: "3846532", CategoryID: "3060070106601", Title: "Pengarna i piracy"},
			{ID: "3800612", CategoryID: "3051302", Title: "I en annan del av Köping: ", PressTitle: "I en annan del av Köping"},
		} {
			ts, c := testServerAndClient(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasSuffix(r.URL.Path, "/publish") {
					xml.NewEncoder(w).Encode(publishing{
						Platform: "tv4",
						Publish:  time.Now(),
					})

					return
				}

				w.Write(testAssetXML(ta))
			})
			defer ts.Close()

			t.Run(ta.ID, func(t *testing.T) {
				a, err := c.Asset(context.Background(), "tv4", ta.ID)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if got, want := a.ID, ta.ID; got != want {
					t.Fatalf("a.ID = %q, want %q", got, want)
				}

				if got, want := a.CategoryID, ta.CategoryID; got != want {
					t.Fatalf("a.CategoryID = %q, want %q", got, want)
				}

				if got, want := a.Title, ta.Title; got != want {
					t.Fatalf("a.Title = %q, want %q", got, want)
				}

				if got, want := a.Metadata.LouisePressTitle, ta.PressTitle; got != want {
					t.Fatalf("a.Metadata.LouisePressTitle = %q, want %q", got, want)
				}
			})
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

type testAsset struct {
	XMLName    xml.Name `xml:"asset"`
	ID         string   `xml:"id,attr"`
	CategoryID string   `xml:"categoryId,attr"`
	Title      string   `xml:"title"`
	PressTitle string   `xml:"metadata>louise-press-title,omitempty"`
}

func testAssetXML(ta testAsset) []byte {
	output, err := xml.Marshal(ta)
	if err != nil {
		return nil
	}

	return output
}
