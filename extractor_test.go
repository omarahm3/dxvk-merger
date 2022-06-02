package main

import (
	"fmt"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	link := "https://cdn.discordapp.com/attachments/483946408676818974/977702339811094538/r5apex.dxvk-cache"

	tests := []struct {
		data        string
		linksNumber int
	}{
		{
			data:        fmt.Sprintf("it should extract link only %s not any other text", link),
			linksNumber: 1,
		},
		{
			data: fmt.Sprintf("it should extract link only %snot any other text", link),
			linksNumber: 1,
		},
		{
			data: fmt.Sprintf("it should extract link only%snot any other text", link),
			linksNumber: 1,
		},
		{
			data: fmt.Sprintf("it should extract link only[%s]not any other text", link),
			linksNumber: 1,
		},
		{
			data: fmt.Sprintf("it should extract link only[%s]not any other text%s", link, link),
			linksNumber: 2,
		},
	}

	for _, test := range tests {
		links := ExtractLinks(test.data)

		if len(links) != test.linksNumber {
			t.Errorf("Expected 1 link, got %d", len(links))
		}

    for _, l := range links {
      if l != link {
        t.Errorf("Expected %q, got %q", link, l)
      }
    }
	}
}
