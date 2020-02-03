package restapi

import (
	"strings"
	"testing"
)

func TestParsePlatforms(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		const platformsJSON = `
    [
      {
        "id": 123,
        "name": "foo-name",
        "publishingRegion": "foo-region",
        "platformGroup": "foo-group"
      },
      {
        "id": 234,
        "name": "bar-name",
        "publishingRegion": "bar-region",
        "platformGroup": "bar-group"
      }
    ]`

		platforms, err := parsePlatforms(strings.NewReader(platformsJSON))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wantPlatforms := []Platform{
			{
				ID:   123,
				Name: "foo-name",
			},
			{
				ID:   234,
				Name: "bar-name",
			},
		}

		if got, want := len(platforms), len(wantPlatforms); got != want {
			t.Fatalf("got %d platforms, want %d", got, want)
		}

		for n := range wantPlatforms {
			platform := platforms[n]
			wantPlatform := wantPlatforms[n]

			if got, want := platform.ID, wantPlatform.ID; got != want {
				t.Errorf("platform[%d].ID = %q, want %q", n, got, want)
			}

			if got, want := platform.Name, wantPlatform.Name; got != want {
				t.Errorf("platform[%d].Name = %q, want %q", n, got, want)
			}
		}
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		_, err := parsePlatforms(strings.NewReader("not-json"))

		if err == nil {
			t.Fatal("err is nil")
		}
	})
}
