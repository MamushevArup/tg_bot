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

func (k *Krisha) mapUrls() string {
	var url strings.Builder
	user := k.User.UserChoice
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
	if user.YearOfBuiltFrom != nil {
		val := strconv.FormatUint(uint64(*user.YearOfBuiltFrom), 10)
		url.WriteString("das[house.year][from]=" + val + "&")
	}
	if user.YearOfBuiltTo != nil {
		val := strconv.FormatUint(uint64(*user.YearOfBuiltTo), 10)
		url.WriteString("das[house.year][to]=" + val + "&")
	}
	if user.PriceFrom != nil {
		val := strconv.FormatUint(*user.PriceFrom, 10)
		url.WriteString("das[price][from]=" + val + "&")
	}
	if user.PriceTo != nil {
		val := strconv.FormatUint(*user.PriceTo, 10)
		url.WriteString("das[price][to]=" + val + "&")
	}
	if user.FloorFrom != nil {
		val := strconv.FormatUint(uint64(*user.FloorFrom), 10)
		url.WriteString("das[flat][from]=" + val + "&")
	}
	if user.FloorTo != nil {
		val := strconv.FormatUint(uint64(*user.FloorTo), 10)
		url.WriteString("das[flat][to]=" + val + "&")
	}
	if user.CheckboxNotFirstFloor != nil {
		url.WriteString("das[floor_not_first]=1&")
	}
	if user.CheckboxNotLastFloor != nil {
		url.WriteString("das[floor_not_last]=1&")
	}
	if user.CheckboxFromOwner != nil {
		url.WriteString("das[who]=1&")
	}
	if user.CheckboxNewBuilding != nil {
		url.WriteString("das[novostroiki]=1&")
	}
	if user.CheckRealEstate != nil {
		url.WriteString("das[_sys.fromAgent]=1&")
	}
	if user.FloorInTheHouseFrom != nil {
		val := strconv.FormatUint(uint64(*user.FloorInTheHouseFrom), 10)
		url.WriteString("das[house.floor_num][from]=" + val + "&")
	}
	if user.FloorInTheHouseTo != nil {
		val := strconv.FormatUint(uint64(*user.FloorInTheHouseTo), 10)
		url.WriteString("das[house.floor_num][to]=" + val + "&")
	}
	if user.AreaFrom != nil {
		url.WriteString("das[live.square][from]=" + *user.AreaFrom + "&")
	}
	if user.AreaTo != nil {
		url.WriteString("das[live.square][to]=" + *user.AreaTo + "&")
	}
	if user.KitchenAreaFrom != nil {
		url.WriteString("das[live.square_k][from]=" + *user.KitchenAreaFrom + "&")
	}
	if user.KitchenAreaTo != nil {
		url.WriteString("das[live.square_k][to]=" + *user.KitchenAreaTo + "&")
	}
	fmt.Println(url.String(), "URL IN THE MAP URL FUNCTION")
	return url.String()
}

const krishaURL = "https://krisha.kz/"

func New(c *colly.Collector, user *models.User) *Krisha {
	return &Krisha{Colly: c, User: user}
}

func (k *Krisha) NewScrap(dups []models.House) (*[]models.House, error) {
	houses := k.scrapSubPage()

	k.scrapMain()
	url := k.mapUrls()
	err := k.visitLink(krishaURL + "prodazha/kvartiry/" + url)
	if err != nil {
		return nil, err
	}
	removeDuplicates(houses, dups)
	return houses, nil
}

func removeDuplicates(houses *[]models.House, dups []models.House) *[]models.House {
	set := make(map[string]bool, 23)
	for i, house := range *houses {
		if set[house.Link] {
			*houses = append((*houses)[:i], (*houses)[i+1:]...)
		}
		set[house.Link] = true
	}
	inter := make([]models.House, len(*houses))
	copy(inter, *houses)
	for i := len(*houses) - 1; i >= 0; i-- {
		h := (*houses)[i]
		for _, d := range dups {
			if h.Link == d.Link {
				if i < 0 || i >= len(*houses) {
					break
				}
				*houses = append((*houses)[:i], (*houses)[i+1:]...)
				fmt.Println("INSIDE THE REMOVE DUPS", len(*houses), len(dups), h.Link)
			}
		}
	}
	copy(dups, inter)
	fmt.Println(len(*houses))
	fmt.Println(len(inter))
	fmt.Println(len(dups))
	inter = []models.House{}
	return houses
}

func (k *Krisha) scrapSubPage() *[]models.House {
	houses := make([]models.House, 0)
	time.Sleep(5 * time.Second)
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
