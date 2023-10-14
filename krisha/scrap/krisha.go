package scrap

import (
	"fmt"
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

type Krisha struct {
	Colly *colly.Collector
	User  *models.User
}

/*
	elements := []string{
		"username",
		"buy_or_rent",
		"type_item",
		"city",
		"rooms",
		"type_house",
		"year_of_built_from",
		"year_of_built_to",
		"price_from",
		"price_to",
		"floor_from",
		"floor_to",
		"checkbox_not_first_floor",
		"checkbox_not_last_floor",
		"checkbox_from_owner",
		"checkbox_new_building",
		"check_real_estate",
		"floor_in_the_house_from",
		"floor_in_the_house_to",
	}

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
func (k *Krisha) mapUrls() string {
	var url strings.Builder
	user := k.User.UserChoice
	if user.BuyOrRent != "" {
		url.WriteString(user.BuyOrRent + "/")
	}
	if user.TypeItem != "" {
		url.WriteString(user.TypeItem + "/")
	}
	if user.City != "" {
		url.WriteString(user.City + "?")
	}
	if user.Rooms != nil {
		for _, room := range user.Rooms {
			url.WriteString("das[live.rooms][]=" + room + "&")
		}
	}
	if user.TypeHouse != nil {
		for _, s := range user.TypeHouse {
			url.WriteString("das[flat.building][]=" + s + "&")
		}
	}
	if user.YearOfBuiltFrom != 0 {
		val := strconv.FormatUint(uint64(user.YearOfBuiltFrom), 10)
		url.WriteString("das[house.year][from]=" + val + "&")
	}
	if user.YearOfBuiltTo != 0 {
		val := strconv.FormatUint(uint64(user.YearOfBuiltTo), 10)
		url.WriteString("das[house.year][to]=" + val + "&")
	}
	if user.PriceFrom != 0 {
		val := strconv.FormatUint(user.PriceFrom, 10)
		url.WriteString("das[price][from]=" + val + "&")
	}
	if user.PriceTo != 0 {
		val := strconv.FormatUint(user.PriceTo, 10)
		url.WriteString("das[price][to]=" + val + "&")
	}
	if user.FloorFrom != 0 {
		val := strconv.FormatUint(uint64(user.FloorFrom), 10)
		url.WriteString("das[flat][from]=" + val + "&")
	}
	if user.FloorTo != 0 {
		val := strconv.FormatUint(uint64(user.FloorTo), 10)
		url.WriteString("das[flat][to]=" + val + "&")
	}
	if user.CheckboxNotFirstFloor {
		url.WriteString("das[floor_not_first]=1&")
	}
	if user.CheckboxNotLastFloor {
		url.WriteString("das[floor_not_last]=1&")
	}
	if user.CheckboxFromOwner {
		url.WriteString("das[who]=1&")
	}
	if user.CheckboxNewBuilding {
		url.WriteString("das[novostroiki]=1&")
	}
	if user.CheckRealEstate {
		url.WriteString("das[_sys.fromAgent]=1&")
	}
	if user.FloorInTheHouseFrom != 0 {
		val := strconv.FormatUint(uint64(user.FloorInTheHouseFrom), 10)
		url.WriteString("das[house.floor_num][from]=" + val + "&")
	}
	if user.FloorInTheHouseTo != 0 {
		val := strconv.FormatUint(uint64(user.FloorInTheHouseTo), 10)
		url.WriteString("das[house.floor_num][to]=" + val + "&")
	}
	if user.AreaFrom != "" {
		url.WriteString("das[live.square][from]=" + user.AreaFrom + "&")
	}
	if user.AreaTo != "" {
		url.WriteString("das[live.square][to]=" + user.AreaTo + "&")
	}
	if user.KitchenAreaFrom != "" {
		url.WriteString("das[live.square_k][from]=" + user.KitchenAreaFrom + "&")
	}
	if user.KitchenAreaTo != "" {
		url.WriteString("das[live.square_k][to]=" + user.KitchenAreaTo + "&")
	}
	return url.String()
}

const krishaURL = "https://krisha.kz/"

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
	url := k.mapUrls()
	fmt.Println(url + "++++++++++++++++++++++++++++++++++++")
	err := k.visitLink(krishaURL + url)
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
