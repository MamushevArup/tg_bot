package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/MamushevArup/krisha-scraper/database/postgres"
	"github.com/MamushevArup/krisha-scraper/krisha/scrap"
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/MamushevArup/krisha-scraper/telegram/inline"
	"github.com/MamushevArup/krisha-scraper/utils"
	"github.com/MamushevArup/krisha-scraper/utils/files"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

func sliceOrder() []string {
	fieldOrder := []string{
		"Ссылка",
		"Заголовок",
		"Цена",
		"Город",
		"Тип дома",
		"Жилой комплекс",
		"Год постройки",
		"Площадь",
		"Санузел",
		"Потолки",
		"Состояние",
		"Бывшее общежитие",
	}
	return fieldOrder
}

func orderOutput(fieldOrder []string, data map[string]interface{}) string {
	var res string
	for _, ord := range fieldOrder {
		v, ok := data[ord]
		if ok {
			switch curr := v.(type) {
			case string:
				res += ord + " : " + v.(string) + "\n"
			case float64:
				floatStr := strconv.FormatFloat(curr, 'f', -1, 64)
				res += ord + " : " + floatStr + "\n"
			}
		}
	}
	return res
}

func (b *Bot) handleCallbackQuery(update *tgbotapi.Update, user *models.User, db *postgres.Sql) {
	chatID := update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	bor := inline.CollectButtonData(inline.BuyOrRent())
	city := inline.CollectButtonData(inline.ChooseCity())
	if bor[update.CallbackQuery.Data] {
		user.BuyOrRent = update.CallbackQuery.Data
		msg.Text = "Выбери город"
		msg.ReplyMarkup = inline.ChooseCity()
		b.sendMessage(&msg)
	} else if city[update.CallbackQuery.Data] {
		user.City = update.CallbackQuery.Data
		val, err := files.ReadTXT("utils/text/txt/text.txt")
		if err != nil {
			log.Println("Cannot read the after-start.txt file ", err.Error())
		}
		db.IntroDataStartCommand(user)
		msg.Text = "Теперь вы можете начать поиск квартир по команде /run\nИли настроить более точный фильтр для этого ознакомьтесь с документацией по команде /help"
		b.sendMessage(&msg)
		msg.Text = val
		b.sendMessage(&msg)
	}
}

func (b *Bot) sendMessage(msg *tgbotapi.MessageConfig) {
	b.bot.Send(msg)
}
func listCities() map[string]string {
	return map[string]string{"Алматы": "almaty", "Астана": "astana", "Шымкент": "shymkent", "Актау": "aktau"}
}
func trimAllSpaces(s string) (uint64, error) {
	var val string
	for _, i2 := range s {
		if i2 != ' ' {
			val += string(i2)
		}
	}
	res, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		log.Println("error with converting string to the uint")
		return 0, err
	}
	return res, nil
}

func (b *Bot) allCommand(val *models.User) string {

	js, err := utils.ConvertToJSONO(val)
	if err != nil {
		log.Println("Cannot convert struct to json all command ", err.Error())
	}
	var hmap map[string]interface{}
	if err = json.Unmarshal([]byte(js), &hmap); err != nil {
		log.Println("cannot convert to the map all command ", err.Error())
	}
	fmt.Println(hmap)
	newM := utils.EnRusUser(hmap)
	var res strings.Builder
	for s, i := range newM {

		switch t := i.(type) {
		case string:
			res.WriteString(s + " : " + t + "\n")
		case []string:
			jn := strings.Join(t, ",")
			res.WriteString(s + " : " + jn + "\n")
		case float64:
			if t == float64(uint(t)) {
				res.WriteString(s + " : " + strconv.FormatUint(uint64(t), 10) + "\n")
			} else {
				res.WriteString(s + " : " + strconv.FormatFloat(t, 'f', -1, 64) + "\n")
			}
		case bool:
			var tof string
			if t {
				tof = "true"
			} else {
				tof = "false"
			}
			res.WriteString(s + " : " + tof + "\n")
		}
	}
	return res.String()
}

func (b *Bot) rangeHouses(houses *[]models.House, update *tgbotapi.Update) {
	for _, house := range *houses {
		val, err := utils.ConvertToJSONO(house)
		fieldOrder := sliceOrder()
		var hmap map[string]interface{}

		if err != nil {
			log.Println("Cannot convert one element to the json ", err)
		}
		if err = json.Unmarshal([]byte(val), &hmap); err != nil {
			log.Println("Error with converting json to the map")
		}
		newM := utils.EnRusHouse(hmap)
		output := orderOutput(fieldOrder, newM)

		msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, output)
		b.sendMessage(&msg2)
	}
	houses = &[]models.House{}
}
func (b *Bot) tickStart(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, db *postgres.Sql) {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	usr, err := db.GetAll()
	if err != nil {
		log.Println("error with getting all info about user ", err.Error())
	}
	c.SetRequestTimeout(90 * time.Second)
	init := scrap.New(c, usr)
	ticker := time.NewTicker(2 * time.Minute)

	dups := make([]models.House, 30)
	for range ticker.C {

		houses, err := init.NewScrap(dups)
		if err != nil {
			msg.Text = "Возникли проблемы с веб-сайтом Krisha.kz ожидайте..."
			break
		}
		if update.Message != nil {
			if update.Message.Command() == "stop" {
				break
			}
		}
		b.rangeHouses(houses, update)
	}
}
