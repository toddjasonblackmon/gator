
package main

import (
    "context"
    "net/http"
    "io"
    "fmt"
    "encoding/xml"
    "html"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

    client := &http.Client{}
    req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("User-Agent", "gator")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if resp.StatusCode > 299 {
		return nil, 
            fmt.Errorf("Response failed with status code: %d and \nbody: %s\n",
		    	resp.StatusCode, dat)
	}
	if err != nil {
		return nil, err
	}

    var feed RSSFeed
    err = xml.Unmarshal(dat, &feed)
    if err != nil {
        return nil, err
    }

    unescapeFeed(&feed)

    return &feed, nil
}

func unescapeFeed(feed *RSSFeed) {
    feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
    feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

    for idx, entry := range feed.Channel.Item {
        feed.Channel.Item[idx].Title = html.UnescapeString(entry.Title)
        feed.Channel.Item[idx].Description = html.UnescapeString(entry.Description)
    }
}


