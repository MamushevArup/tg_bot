package scrap

import (
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func text(e *colly.HTMLElement, query string) string {
	return e.ChildText(query)
}

func loopDiv(e *colly.HTMLElement, query string) []string {
	arr := make([]string, 0)
	e.ForEach(query, func(_ int, element *colly.HTMLElement) {
		val := trimSpace(element.Text)
		arr = append(arr, val)
	})
	return arr
}

func parseInt(priceStr string) (int, error) {
	val := regexp.MustCompile(`[^\d]+`).ReplaceAllString(priceStr, "")
	priceInt, err := strconv.Atoi(val)
	if err != nil {
		return -1, err
	}

	return priceInt, nil
}

func (k *Krisha) visitLink(url string) error {
	err := k.Colly.Visit(url)
	if err != nil {
		log.Fatal("Error while parsing this link", url)
		return err
	}
	return nil
}

func removeTags(e *colly.HTMLElement, goquery string) {
	e.ForEach(goquery, func(_ int, a *colly.HTMLElement) {
		a.DOM.Remove()
	})
}

func (k *Krisha) scrapMain() {
	k.Colly.OnHTML("div.a-card__header-left", func(e *colly.HTMLElement) {
		link := e.ChildAttrs("a[href].a-card__title", "href")

		err := k.visitLink(krishaURL + link[0])
		if err != nil {
			return
		}
	})
}

func trimSpace(arg string) string {
	return strings.TrimSpace(arg)
}
