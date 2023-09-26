package scrap

import (
	"fmt"
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const krishaURL = "https://krisha.kz"

func NewScrap() *[]models.House {
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

func scrapSubPage(c *colly.Collector) *[]models.House {
	houses := make([]models.House, 0)
	c.OnHTML("div.layout__content", func(e *colly.HTMLElement) {
		hmap := make(map[string]string)
		removeTags("a.btm-map", e)
		titles := loopDiv(e, "div.offer__info-title")
		keys := loopDiv(e, "div.offer__advert-short-info")
		price := text(e, "div.offer__price, p.offer__price")
		desc := text(e, "h1")
		for i, title := range titles {
			hmap[title] = keys[i]
		}
		pr, err := parsePrice(price)
		yearofbuild, err := parsePrice(hmap["Год постройки"])
		if err != nil {
			log.Fatal("Cannot convert the string to the int", err)
		}

		house := &models.House{
			Link:               e.Request.URL.String(),
			Intro:              desc,
			Price:              pr,
			City:               trimSpace(hmap["Город"]),
			HouseType:          trimSpace(hmap["Тип дома"]),
			ResidentialComplex: trimSpace(hmap["Жилой комплекс"]),
			YearOfBuild:        uint16(yearofbuild),
			Floor:              trimSpace(hmap["Этаж"]),
			Area:               trimSpace(hmap["Площадь, м²"]),
			Bathroom:           trimSpace(hmap["Санузел"]),
			Ceil:               trimSpace(hmap["Потолки"]),
			FormerHostel:       trimSpace(hmap["Бывшее общжитие"]),
			State:              trimSpace(hmap["Состояние"]),
			CreatedAt:          time.Now().Format("2006-01-02 15:04:05"),
		}
		houses = append(houses, *house)
	})
	return &houses
}
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

func parsePrice(priceStr string) (int, error) {
	// Remove non-numeric characters (e.g., spaces, '〒') using regular expressions
	re := regexp.MustCompile(`[^\d]+`)
	cleanPriceStr := re.ReplaceAllString(priceStr, "")

	// Convert to an integer
	priceInt, err := strconv.Atoi(cleanPriceStr)
	if err != nil {
		return -1, err
	}

	return priceInt, nil
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

func scrapMain(c *colly.Collector) {
	c.OnHTML("div.a-card__header-left", func(e *colly.HTMLElement) {
		link := e.ChildAttrs("a[href].a-card__title", "href")

		err := visitLink(c, krishaURL+link[0])
		if err != nil {
			return
		}
	})
}
func trimSpace(arg string) string {
	return strings.TrimSpace(arg)
}
