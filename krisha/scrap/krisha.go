package scrap

import (
	"fmt"
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/gocolly/colly"
	"log"
	"time"
)

type Krisha struct {
	Colly *colly.Collector
	User  *models.User
}

/*
url for krisha scrap
order of it doesn't matter
random number or zero is not applicable for this because of error
продажа по комнатам ?das[live.rooms][]=0?das[live.rooms][]=1&?das[live.rooms][]=2&?das[live.rooms][]=3&?das[live.rooms][]=4&?das[live.rooms][]=5&?das[live.rooms][]=5.100&
тип дома ?das[flat.building][]=i++ for every type
yearOfBuilt from ?das[house.year][from]={value} to das[house.year][to]={value}
price from ?das[price][from]={value} to das[price][to]={value}
floor from ?das[flat][from]={value} to das[flat][to]={v}
checkbox for not_first floor ?das[floor_not_first]=1 for not last das[floor_not_last]=1
checkbox for от хозяев ?das[who]=1
от новостроек das[novostroiki]=1
от крыша агентов das[_sys.fromAgent]=1
floor in the house from ?das[house.floor_num][from]={v} to das[house.floor_num][to]={v}
total area from das[live.square][from]={v} to das[live.square][to]={v}
area kitchen from das[live.square_k][from]={v} to das[live.square_k][to]={v}
*/

//func mapUrls() {
//	urls := map[string]string{
//		"city" : "/",
//		"rooms" : "das[live.rooms][]=",
//
//	}
//}

const krishaURL = "https://krisha.kz"

func New(c *colly.Collector, user *models.User) *Krisha {
	return &Krisha{Colly: c, User: user}
}

func (k *Krisha) NewScrap() (*[]models.House, error) {
	houses := k.scrapSubPage()
	seen := make(map[string]bool)
	fmt.Println(seen, "Before")
	removeDuplicates(houses, seen)
	fmt.Println(seen, "After")
	k.scrapMain()
	fmt.Println(houses)
	err := k.visitLink(krishaURL + "/prodazha/kvartiry/almaty")
	if err != nil {
		return nil, err
	}
	return houses, nil
}

func removeDuplicates(houses *[]models.House, seen map[string]bool) *[]models.House {
	result := (*houses)[:0] // Create a new slice with the same underlying array
	for _, house := range *houses {
		if !seen[house.Link] {
			seen[house.Link] = true
			result = append(result, house)
		}
	}
	*houses = result // Update the original array with the unique houses
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
			log.Fatal("Cannot convert the string to the int ", err, hmap["Год постройки"])
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
