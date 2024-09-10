package website

import "testing"

func TestNormalizeURL(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"remove trailing slash", "https://www.example.com/", "https://www.example.com"},
		{"normalize case", "HTTPS://www.example.com", "https://www.example.com"},
		{"remove query string", "https://www.example.com?query=string", "https://www.example.com"},
		{"remove anchor", "https://www.example.com#contact", "https://www.example.com"},
		{"criteria combined", "https://www.example.com/?test=test#contact", "https://www.example.com"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeURL(tt.input); got != tt.expected {
				t.Errorf("normalizeURL() = %s, want: %s", got, tt.expected)
			}
		})
	}
}
