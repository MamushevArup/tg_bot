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

func (b *Bot) HandleUpdate(update *tgbotapi.Update, user *models.User, db *postgres.Sql) {

	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			b.handleMessageCommand(update, user, db)
		}
		user.Username = update.Message.Chat.UserName
	} else if update.CallbackQuery != nil {
		b.handleCallbackQuery(update, user, db)

	}
}

func (b *Bot) handleMessageCommand(update *tgbotapi.Update, user *models.User, db *postgres.Sql) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	b.handleCommand(&msg, update, user, db)
	b.sendMessage(&msg)
}

func (b *Bot) handleCommand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, user *models.User, db *postgres.Sql) {
	command := update.Message.Command()
	switch command {
	case "start":
		val, err := files.ReadTXT("utils/text/txt/start-msg.txt")
		if err != nil {
			log.Println("Error with reading start-msg.txt ", err.Error())
		}
		msg.Text = val
		msg.ReplyMarkup = inline.ChooseCity()
	case "help":
		val, err := files.ReadTXT("utils/text/txt/text.txt")
		if err != nil {
			log.Println("Error with reading the docs to the bot", err)
		}
		msg.Text = val
	case "city":
		hmap := inline.ChooseCity()
		msg.ReplyMarkup = hmap
		val := update.CallbackQuery.Data
		user.City = val
		err := db.UpdateCity(val)
		if err != nil {
			msg.Text = "Не удается обновить город ознакомьтесь с документацией по команде /help и повторите попытку"
			log.Println("Cannot update city in the city switch case ", err)
			return
		}
		val, err = files.ReadTXT("utils/text/txt/text.txt")
		if err != nil {
			log.Println("Cannot read from the text.txt file ", err.Error())
		}
		msg.Text = val
	case "region":
		msg.Text = "Выберите район для поиска квартир"
		//val := update.CallbackData()
	case "rooms":
		val := update.Message.CommandArguments()
		if val == "" {
			user.UserChoice.Rooms = nil
			err := db.UpdateRooms(nil)
			if err != nil {
				log.Println("Error with updating rooms in the rooms switch case ", err)
			}
			return
		}
		arrOfRooms := strings.Split(val, ",")
		for _, s := range arrOfRooms {
			curr, err := strconv.Atoi(s)
			if err != nil {
				log.Println("Cannot convert room to the int ", err.Error())
			}
			if s != "5+" {
				if err != nil || curr > 5 || curr < 1 {
					msg.Text = "Возможно вы ввели в неправильном формате"
					return
				}
			} else {
				s = "5.100"
			}

		}
		user.UserChoice.Rooms = arrOfRooms
		err := db.UpdateRooms(arrOfRooms)
		if err != nil {
			log.Println("Error with updating rooms in the rooms switch case ", err)
			msg.Text = "Не удается обновить кол-во комнат ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		msg.Text = "Отлично!\nТеперь вы можете запустить поиск командой /run или продолжите настройку\nДоступные команды доступны /help"
	case "type":
		val := update.Message.CommandArguments()
		if val == "" {
			user.UserChoice.TypeHouse = nil
			err := db.UpdateType(nil)
			if err != nil {
				log.Println("Error with updating rooms in the rooms switch case ", err)
			}
			return
		}
		arrOfTypeHouse := strings.Split(val, ",")
		for _, s := range arrOfTypeHouse {
			curr, err := strconv.Atoi(s)
			if err != nil || curr > 3 || curr < 1 {
				msg.Text = "Возможно вы ввели в неправильном формате"
				return
			}
		}
		user.UserChoice.TypeHouse = arrOfTypeHouse
		err := db.UpdateType(arrOfTypeHouse)
		if err != nil {
			log.Println("Error with updating rooms in the rooms switch case ", err)
			msg.Text = "Не удается обновить тип дома ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		msg.Text = "Отлично!\nТеперь вы можете запустить поиск командой /run или продолжите настройку\nДоступные команды доступны /help"
	case "builtfrom":
		val := update.Message.CommandArguments()
		if val == "" {
			var n uint
			err := db.UpdateBuiltFrom(n)
			if err != nil {
				log.Println("Error with updating builtFrom in the builtFrom switch case ", err)
			}
			user.UserChoice.YearOfBuiltFrom = nil
			return

		}
		year, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateBuiltFrom(uint(year))
		if err != nil {
			log.Println("Error with updating builtFrom in the builtFrom switch case ", err)
			msg.Text = "Не удается обновить год постройки от ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(year)
		user.UserChoice.YearOfBuiltFrom = &curr
	case "builtto":
		val := update.Message.CommandArguments()
		if val == "" {
			var n uint
			err := db.UpdateBuiltTo(n)
			if err != nil {
				log.Println("Error with updating builtTo in the builtTo switch case ", err)
			}
			user.UserChoice.YearOfBuiltTo = nil
			return
		}
		year, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateBuiltTo(uint(year))
		if err != nil {
			log.Println("Error with updating builtTo in the builtTo switch case ", err)
			msg.Text = "Не удается обновить год постройки до  ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(year)
		user.UserChoice.YearOfBuiltTo = &curr

	case "pricefrom":
		val := update.Message.CommandArguments()
		if val == "" {
			var n uint64
			err := db.UpdatePriceFrom(n)
			if err != nil {
				log.Println("Error with updating priceFrom in the priceFrom switch case ", err)
			}
			user.UserChoice.PriceFrom = nil
			return
		}
		price, err := trimAllSpaces(val)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdatePriceFrom(price)
		if err != nil {
			log.Println("Error with updating priceFrom in the priceFrom switch case ", err)
			msg.Text = "Не удается обновить стартовую цену ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		user.UserChoice.PriceFrom = &price
	case "priceto":
		val := update.Message.CommandArguments()
		if val == "" {
			var n uint64
			err := db.UpdatePriceTo(n)
			if err != nil {
				log.Println("Error with updating priceTo in the priceTo switch case ", err)
			}
			user.UserChoice.PriceTo = nil
			return
		}
		price, err := trimAllSpaces(val)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdatePriceTo(price)
		if err != nil {
			log.Println("Error with updating priceTo in the priceTo switch case ", err)
			msg.Text = "Не удается обновить сумму до ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		user.UserChoice.PriceTo = &price
	case "floorfrom":
		val := update.Message.CommandArguments()
		if val == "" {
			var n uint64
			err := db.UpdateFloorFrom(n)
			if err != nil {
				log.Println("Error with updating floorFrom ", err)
			}
			user.UserChoice.FloorFrom = nil
			return
		}
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateFloorFrom(floor)
		if err != nil {
			log.Println("Error with updating floorFrom ", err)
			msg.Text = "Не удается обновить стартовый этаж ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(floor)
		user.UserChoice.FloorFrom = &curr
	case "floorto":
		val := update.Message.CommandArguments()
		if val == "" {
			var n uint64
			err := db.UpdateFloorTo(n)
			if err != nil {
				log.Println("Error with updating floorFrom ", err)

			}
			user.UserChoice.FloorTo = nil
			return
		}

		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateFloorTo(floor)
		if err != nil {
			log.Println("Error with updating floorFrom ", err)
			msg.Text = "Не удается обновить этаж до  ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(floor)
		user.UserChoice.FloorTo = &curr
	case "notfirst":
		flag := true
		user.UserChoice.CheckboxNotFirstFloor = &flag
		err := db.UpdateNotFirstFloor(true)
		if err != nil {
			log.Println("Error with updating notFirst ", err)
			msg.Text = "Не удается обновить параметр не первый этаж ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "notlast":
		flag := true
		user.UserChoice.CheckboxNotLastFloor = &flag
		err := db.UpdateNotLastFloor(true)
		if err != nil {
			log.Println("Error with updating notLast ", err)
			msg.Text = "Не удается обновить параметр не последний этаж ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "fromowner":
		flag := true
		user.UserChoice.CheckboxFromOwner = &flag
		err := db.UpdateFromOwner(true)
		if err != nil {
			log.Println("Error with updating fromOwner ", err)
			msg.Text = "Не удается обновить параметр от владельца ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "newbuilding":
		flag := true
		user.UserChoice.CheckboxNewBuilding = &flag
		err := db.UpdateNewBuilding(true)
		if err != nil {
			log.Println("Error with updating newBuilding ", err)
			msg.Text = "Не удается обновить параметр новостройка ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "realestate":
		flag := true
		user.UserChoice.CheckRealEstate = &flag
		err := db.UpdateRealEstate(true)
		if err != nil {
			log.Println("Error with updating realEstate ", err)
			msg.Text = "Не удается обновить параметр от крыша агента ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "floorhfrom":
		val := update.Message.CommandArguments()
		if val == "" {
			var n uint64
			err := db.UpdateFloorInTheHouseFrom(n)
			if err != nil {
				log.Println("Error with updating floorInTheHouseFrom ", err)

			}
			user.UserChoice.FloorInTheHouseFrom = nil
			return
		}
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateFloorInTheHouseFrom(floor)
		if err != nil {
			log.Println("Error with updating floorInTheHouseFrom ", err)
			msg.Text = "Не удается обновить параметр этаж в доме от ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(floor)
		user.UserChoice.FloorInTheHouseFrom = &curr
	case "floorhto":
		val := update.Message.CommandArguments()
		if val == "" {
			var n uint64
			err := db.UpdateFloorInTheHouseTo(n)
			if err != nil {
				log.Println("Error with updating floorInTheHouseTo ", err)

			}
			user.UserChoice.FloorInTheHouseTo = nil
			return
		}
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateFloorInTheHouseTo(floor)
		if err != nil {
			log.Println("Error with updating floorInTheHouseTo ", err)
			msg.Text = "Не удается обновить параметр этаж в доме до ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(floor)
		user.UserChoice.FloorInTheHouseTo = &curr
	case "areafrom":
		val := update.Message.CommandArguments()
		if val == "" {
			user.UserChoice.AreaFrom = nil
			err := db.UpdateAreaFrom("")
			if err != nil {
				log.Println("Error with updating areaFrom ", err)
			}
			return
		}
		user.UserChoice.AreaFrom = &val
		err := db.UpdateAreaFrom(val)
		if err != nil {
			log.Println("Error with updating areaFrom ", err)
			msg.Text = "Не удается обновить параметр стартовая площадь ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "areato":
		val := update.Message.CommandArguments()
		if val == "" {
			user.UserChoice.AreaTo = nil
			err := db.UpdateAreaTo(val)
			if err != nil {
				log.Println("Error with updating areaTo ", err)
			}
			return
		}
		user.UserChoice.AreaTo = &val
		err := db.UpdateAreaTo(val)
		if err != nil {
			log.Println("Error with updating areaTo ", err)
			msg.Text = "Не удается обновить параметр площадь до ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "kitfrom":
		val := update.Message.CommandArguments()
		if val == "" {
			user.UserChoice.KitchenAreaFrom = nil
			err := db.UpdateKitchenFrom("")
			if err != nil {
				log.Println("Error with updating kitchen ", err)
			}
			return
		}

		user.UserChoice.KitchenAreaFrom = &val
		err := db.UpdateKitchenFrom(val)
		if err != nil {
			log.Println("Error with updating kitchen ", err)
			msg.Text = "Не удается обновить параметр площадь кухни от ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "kitto":
		val := update.Message.CommandArguments()
		if val == "" {
			user.UserChoice.KitchenAreaTo = nil
			err := db.UpdateKitchenTo("")
			if err != nil {
				log.Println("Error with updating kitchenTo ", err)

			}
			return
		}
		user.UserChoice.KitchenAreaTo = &val
		err := db.UpdateKitchenTo(val)
		if err != nil {
			log.Println("Error with updating kitchenTo ", err)
			msg.Text = "Не удается обновить параметр площадь кухни до ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "all":
		val, err := db.GetAll()
		if err != nil {
			log.Println("Cannot convert to the struct all command ", err.Error())
		}
		res := b.allCommand(val)
		msg.Text = res

	case "run":
		state := db.IsRunning()
		if state {
			msg.Text = "Поиск квартир уже идет ожидайте"
			return
		}
		err := db.SetRunning(true)
		if err != nil {
			return
		}
		msg.Text = "Идет поиск подходящих обьявлений\nТелеграм оповестит вас когда мы найдем подходящие обьявления\nДля этого уберите беззвучный режим."
		go func() {
			b.tickStart(msg, update, db)
		}()
	case "stop":
		err := db.SetRunning(false)
		if err != nil {
			return
		}
		msg.Text = "Вы остановили поиск обьявлений теперь можете изменить настройки\nПо комманде /help смотрите все доступные комманды"
	default:
		msg.Text = "Нет такой комманды. Ознакомьтесь с документацией"
	}
}

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
	city := inline.CollectButtonData(inline.ChooseCity())
	region := inline.CollectButtonData(inline.ChooseRegion(update.CallbackQuery.Data))
	if db.GetCity() != "" {
		msg.Text = "Для изменения города используйте команду /city\nДля изменения региона используйте команду /region"
		b.sendMessage(&msg)
		return
	}
	if city[update.CallbackQuery.Data] {
		user.City = update.CallbackQuery.Data
		msg.Text = "Выберите один из районов города "
		msg.ReplyMarkup = inline.ChooseRegion(update.CallbackQuery.Data)
		b.sendMessage(&msg)
	} else if region[update.CallbackQuery.Data] {
		if update.CallbackQuery.Data == "-" {
			update.CallbackQuery.Data = user.City
		}
		user.City = update.CallbackQuery.Data
		db.IntroDataStartCommand(user)
		fmt.Println("Enter here the value")
		msg.Text = update.CallbackQuery.Data
		b.sendMessage(&msg)
	}
}

func (b *Bot) sendMessage(msg *tgbotapi.MessageConfig) {
	_, err := b.bot.Send(msg)
	if err != nil {
		return
	}
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
