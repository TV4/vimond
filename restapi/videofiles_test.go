package restapi

import "testing"

func TestVideofilesPath(t *testing.T) {
	for _, tt := range []struct {
		assetID string
		want    string
	}{
		{"", "/api/admin/asset//videofiles"},
		{"123", "/api/admin/asset/123/videofiles"},
		{"456", "/api/admin/asset/456/videofiles"},
	} {
		c := &Client{}

		if got := c.videofilesPath(tt.assetID); got != tt.want {
			t.Fatalf("c.videofilesPath(%q) = %q, want %q", tt.assetID, got, tt.want)
		}
	}
}
