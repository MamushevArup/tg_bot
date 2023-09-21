package scrap

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

const krishaURL = "https://krisha.kz"

func Scrap() {
	c := colly.NewCollector()
	var arrjson []map[string]string
	scrapSubPage(c, arrjson)
	scrapMain(c)
	err := visitLink(c, krishaURL+"/prodazha/kvartiry/")
	if err != nil {
		return
	}
	fmt.Println(arrjson)
}
func visitLink(c *colly.Collector, url string) error {
	err := c.Visit(url)
	if err != nil {
		log.Fatal("Error while parsing this link", url)
		return err
	}
	return nil
}
func removeTags(goquery string, e *colly.HTMLElement) {
	e.ForEach(goquery, func(_ int, a *colly.HTMLElement) {
		// Check if the <a> tag's class meets your removal condition
		a.DOM.Remove()
	})
}
func scrapSubPage(c *colly.Collector, arrjson []map[string]string) {
	c.OnHTML("div.offer__short-description", func(e *colly.HTMLElement) {

		titles := make([]string, 0)
		keys := make([]string, 0)
		hmap := make(map[string]string)
		removeTags("a.btm-map", e)
		e.ForEach("div.offer__info-title", func(_ int, element *colly.HTMLElement) {
			titles = append(titles, element.Text)
		})
		e.ForEach("div.offer__advert-short-info", func(_ int, element *colly.HTMLElement) {
			val := strings.TrimSpace(element.Text)
			keys = append(keys, val)
		})
		for i, title := range titles {
			hmap[title] = keys[i]
		}
		arrjson = append(arrjson, hmap)
	})
}
func scrapMain(c *colly.Collector) {
	c.OnHTML("div.a-card__header-left", func(e *colly.HTMLElement) {
		link := e.ChildAttrs("a[href].a-card__title", "href")
		//e.ForEach("a[href].a-card__title", func(_ int, element *colly.HTMLElement) {
		//	fmt.Println(element.Text)
		//})
		err := visitLink(c, krishaURL+link[0])
		if err != nil {
			return
		}
	})
}
