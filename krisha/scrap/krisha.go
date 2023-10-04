package scrap

import (
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/gocolly/colly"
	"log"
	"time"
)

type Krisha struct {
	Colly *colly.Collector
	Link  string
}

const krishaURL = "https://krisha.kz"

func New(c *colly.Collector, link string) *Krisha {
	return &Krisha{Colly: c, Link: link}
}

func (k *Krisha) NewScrap() *[]models.House {
	houses := k.scrapSubPage()
	k.scrapMain()
	err := k.visitLink(krishaURL + "/prodazha/kvartiry/")
	if err != nil {
		return houses
	}

	return houses
}

func (k *Krisha) scrapSubPage() *[]models.House {
	houses := make([]models.House, 0)
	k.Colly.OnHTML("div.layout__content", func(e *colly.HTMLElement) {
		hmap := make(map[string]string)
		removeTags(e, "a.btm-map")
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
			YearOfBuild:        yearofbuild,
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
