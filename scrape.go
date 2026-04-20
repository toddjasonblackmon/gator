package main

import (
    "context"
    "fmt"
    "os"
	"github.com/toddjasonblackmon/gator/internal/database"
    "time"
    "database/sql"
)

func scrapeFeeds(s *state) {
	ctx := context.Background()

    feed, err := s.db.GetNextFeedToFetch(ctx)
    if err != nil {
		fmt.Println("Unable to get next feed")
		os.Exit(1)
    }

    err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams {
                                        ID: feed.ID, 
                                        LastFetchedAt: sql.NullTime{
                                            Time: time.Now(),
                                            Valid: true}})
    if err != nil {
		fmt.Println("Unable to mark feed as fetched")
		os.Exit(1)
    }

    rss, err := fetchFeed(ctx, feed.Url)
    if err != nil {
		fmt.Println("Unable to fetch feed")
		os.Exit(1)

    }
    
    fmt.Println(rss.Channel.Title)
	for idx := range rss.Channel.Item {
        title := rss.Channel.Item[idx].Title
        if len(title) > 0 {
            fmt.Printf("  - %s\n", rss.Channel.Item[idx].Title)
        }
	}
    fmt.Print("\n")

}
