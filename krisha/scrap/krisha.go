package scrap

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

const krishaURL = "https://krisha.kz"

func NewScrap() *[]map[string]string {
	c := colly.NewCollector()
	arrjson := scrapSubPage(c)
	scrapMain(c)
	err := visitLink(c, krishaURL+"/prodazha/kvartiry/")
	if err != nil {
		return nil
	}
	fmt.Println(arrjson)
	return arrjson
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
		a.DOM.Remove()
	})
}
func scrapSubPage(c *colly.Collector) *[]map[string]string {
	var arrjson []map[string]string
	c.OnHTML("div.layout__content", func(e *colly.HTMLElement) {
		titles := make([]string, 0)
		keys := make([]string, 0)
		hmap := make(map[string]string)
		removeTags("a.btm-map", e)
		e.ForEach("div.offer__info-title", func(i int, element *colly.HTMLElement) {
			fmt.Println(element, "------")
			titles = append(titles, element.Text)
		})
		e.ForEach("div.offer__advert-short-info", func(_ int, element *colly.HTMLElement) {
			val := strings.TrimSpace(element.Text)
			keys = append(keys, val)
		})
		price := e.ChildText("div.offer__price, p.offer__price")
		desc := e.ChildText("h1")
		fmt.Println(desc)
		hmap["Ввод"] = desc
		hmap["Цена"] = price
		for i, title := range titles {
			hmap[title] = keys[i]
		}
		hmap["Ссылка"] = e.Request.URL.String()
		arrjson = append(arrjson, hmap)
	})

	return &arrjson
}
func scrapMain(c *colly.Collector) {
	c.OnHTML("div.a-card__header-left", func(e *colly.HTMLElement) {
		link := e.ChildAttrs("a[href].a-card__title", "href")

		err := visitLink(c, krishaURL+link[0])
		if err != nil {
			return
		}
	})
}
