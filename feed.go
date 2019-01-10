package exrays

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

var appUrl = "https://steamcommunity.com/games/%s/rss/"
var regex = regexp.MustCompile(`<img\s+src="([^"]+)"\s*\\?>`)

type AppImg struct {
	Image string
	Link  string
	Time  time.Time
}

type Feed struct {
	Items []struct {
		Description string `xml:"description"`
		Link        string `xml:"link"`
		PubDate     string `xml:"pubDate"`
	} `xml:"channel>item"`
}

func PullApp(app string) []AppImg {
	url := fmt.Sprintf(appUrl, app)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var rss Feed
	if err = xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		panic(err)
	}

	results := make([]AppImg, 0)

	for _, item := range rss.Items {
		for _, image := range regex.FindAllStringSubmatch(item.Description, -1) {
			var img AppImg
			img.Image = image[1]
			img.Link = item.Link
			img.Time, err = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
			if err != nil {
				panic(err)
			}

			results = append(results, img)
		}
	}

	return results
}
