package restapi

import (
	"fmt"
	"testing"
	"time"
)

func TestAuthorizationHeader(t *testing.T) {
	for i, tt := range []struct {
		method        string
		path          string
		apiKey        string
		secret        string
		now           time.Time
		date          string
		authorization string
	}{
		{
			"GET", "foo", "xyz", "123",
			time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			"Tue, 10 Nov 2009 23:00:00 +0000",
			"SUMO xyz:6hJ55znNIpZBWG6kDweLxr++bWQ=",
		},
		{
			"GET", "foo", "xyz", "123",
			time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC),
			"Wed, 10 Nov 2010 23:00:00 +0000",
			"SUMO xyz:282K7POkdp3X7mas9uVoWYVkByM=",
		},
		{
			"GET", "bar", "abc", "456",
			time.Date(2017, time.April, 18, 12, 0, 0, 0, time.UTC),
			"Tue, 18 Apr 2017 12:00:00 +0000",
			"SUMO abc:sVpYYuC5QuwhY2n7uRUESJ23d5o=",
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			h := authorizationHeader(tt.method, tt.path, tt.now, tt.apiKey, tt.secret)

			if got, want := h.Get("Date"), tt.date; got != want {
				t.Fatalf(`h.Get("Date") = %q, want %q`, got, want)
			}

			if got, want := h.Get("Authorization"), tt.authorization; got != want {
				t.Fatalf(`h.Get("Authorization") = %q, want %q`, got, want)
			}
		})
	}
}

func TestComputeHmacSha1(t *testing.T) {
	for _, tt := range []struct {
		message, secret, want string
	}{
		{"foo", "foo", "sputg0LPIDrWDLp2X2pdKQ6NJiI="},
		{"foo", "bar", "hdFVxV7ShqMAvRzxJN4I2H6RTzo="},
		{"bar", "bar", "Mi4jBzJ4+Eg65AycaLlnIE6Ecs0="},
		{"baz", "bar", "HZy29G84U9+mdRk758n6igvo3iU="},
		{"baz", "foo", "9dtE3Kv9nf2Mxnr6HeCHWFIphKQ="},
	} {
		t.Run(tt.message, func(t *testing.T) {
			if got := computeHmacSha1(tt.message, tt.secret); got != tt.want {
				t.Fatalf("computeHmacSha1(%q, %q) = %q, want %q", tt.message, tt.secret, got, tt.want)
			}
		})
	}
}
