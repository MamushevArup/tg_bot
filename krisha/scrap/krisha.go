package scrap

import (
	"fmt"
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/MamushevArup/krisha-scraper/utils"
	"github.com/gocolly/colly"
	"log"
	"time"
)

const krishaURL = "https://krisha.kz"

func NewScrap() string {
	c := colly.NewCollector()
	houses := scrapSubPage(c)
	scrapMain(c)
	err := visitLink(c, krishaURL+"/prodazha/kvartiry/")
	if err != nil {
		return ""
	}
	housesJSON, err := utils.ConvertToJSON(houses)
	if err != nil {
		log.Fatal("Cannot convert to the json")
		return ""
	}
	fmt.Print(housesJSON)
	return housesJSON
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
		pr, err := parseInt(price)
		yearofbuild, err := parseInt(hmap["Год постройки"])
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
