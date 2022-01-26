package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func check(err error) {
    if err != nil {
        panic(err)
    }
}

type info struct {
    cnt     string
    title   string
    desc    string
}

func main() {

    // Get an url path
    path := "https://flows.nodered.org/search?type=node&sort=downloads&page="

    // Have to stdin total number of pages
    var pages int
    fmt.Scanln(&pages)

    infos := make([]info, 0)

    for i := 0; i < pages; i++ {
        // create context
        ctx, cancel := chromedp.NewContext(context.Background())
        defer cancel()

        ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
        defer cancel()

        // Get HTML of node list
        var html string
        err := chromedp.Run(ctx,
            chromedp.Navigate(path + fmt.Sprintf("%d", pages)),
            chromedp.WaitVisible(`.gistbox-type`, chromedp.NodeVisible, chromedp.ByQuery),
            chromedp.InnerHTML(`.gistlist`, &html, chromedp.NodeVisible, chromedp.ByQuery),
        )
        check(err)

        // Parce HTML
        htmlReader := strings.NewReader(html)
        doc, err := goquery.NewDocumentFromReader(htmlReader)
        check(err)

        // Find the review items
        doc.Find(".gistbox-node").Each(func(cnt int, s *goquery.Selection) {
              // For each item found, get the title
              infos = append(infos, info{
                  cnt: fmt.Sprintf("%d-%d", i, cnt),
                  title: s.Find(".gisttitle").Text(),
                  desc: s.Find(".gistdescription").Text(),
              })
        })
    }

    for _, el := range infos {
        fmt.Printf("%s:\t%s\n%s\n\n", el.cnt, el.title, el.desc)
    }
}
