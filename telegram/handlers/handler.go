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

type tempCity struct {
	buttonCity  string
	realEstate  bool
	newBuilding bool
	fromOwner   bool
	notLast     bool
	notFirst    bool
}

var tempC = &tempCity{}

func (b *Bot) HandleUpdate(update *tgbotapi.Update, user *models.User, db *postgres.Sql, last *[]string) {
	if len(*last) > 1 {
		temp := (*last)[len(*last)-1]
		(*last)[len(*last)-1] = (*last)[0]
		(*last)[0] = temp
		*last = (*last)[:len(*last)-1]
	}
	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			*last = append(*last, update.Message.Command())
			b.handleMessageCommand(update, user, db, last)
		}
		user.Username = update.Message.Chat.UserName
	} else if update.CallbackQuery != nil {
		switch (*last)[len(*last)-1] {
		case "start":
			b.handleCallbackQuery(update, user, db)
		case "region":
			b.handleRegion(update, user, db)
		case "city":
			b.handleCity(update, user, db)
		}
	}

	fmt.Println(last)

}

func (b *Bot) handleMessageCommand(update *tgbotapi.Update, user *models.User, db *postgres.Sql, last *[]string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	b.handleCommand(&msg, update, user, db, last)
	b.sendMessage(&msg)
}

func (b *Bot) handleCommand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, user *models.User, db *postgres.Sql, last *[]string) {
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
		msg.Text = "Выберите город"
		msg.ReplyMarkup = inline.ChooseCity()
	case "region":
		msg.Text = "Выберите район для поиска квартир"
		msg.ReplyMarkup = inline.ChooseRegion(tempC.buttonCity)
		//val := update.CallbackData()
	case "rooms":
		val := update.Message.CommandArguments()
		if val == "" {
			user.Rooms = nil
			err := db.UpdateRooms(nil)
			if err != nil {
				log.Println("Error with updating rooms in the rooms switch case ", err)
				msg.Text = "Не удается обновить кол-во комнат ознакомьтесь с документацией по команде /help и повторите попытку"
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
			if err != nil || curr > 4 || curr < 1 {
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

			user.YearOfBuiltFrom = nil
			err := db.UpdateBuiltFrom(nil)
			if err != nil {
				log.Println("Error with updating builtFrom in the builtFrom switch case ", err)
				msg.Text = "Не удается обновить год постройки от ознакомьтесь с документацией по команде /help и повторите попытку"
			}
			return
		}
		year, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		curr := uint(year)
		err = db.UpdateBuiltFrom(&curr)
		if err != nil {
			log.Println("Error with updating builtFrom in the builtFrom switch case ", err)
			msg.Text = "Не удается обновить год постройки от ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}

		user.UserChoice.YearOfBuiltFrom = &curr
	case "builtto":
		val := update.Message.CommandArguments()
		if val == "" {
			user.YearOfBuiltTo = nil
			err := db.UpdateBuiltTo(nil)
			if err != nil {
				log.Println("Error with updating builtTo in the builtTo switch case ", err)
				msg.Text = "Не удается обновить год постройки до  ознакомьтесь с документацией по команде /help и повторите попытку"
			}
			return
		}
		year, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		t := uint(year)
		err = db.UpdateBuiltTo(&t)
		if err != nil {
			log.Println("Error with updating builtTo in the builtTo switch case ", err)
			msg.Text = "Не удается обновить год постройки до  ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		user.UserChoice.YearOfBuiltTo = &t

	case "pricefrom":
		val := update.Message.CommandArguments()
		if val == "" {

			user.PriceFrom = nil
			err := db.UpdatePriceFrom(nil)
			if err != nil {
				log.Println("Error with updating priceFrom in the priceFrom switch case ", err)
				msg.Text = "Не удается обновить стартовую цену ознакомьтесь с документацией по команде /help и повторите попытку"
			}

			return
		}
		price, err := trimAllSpaces(val)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdatePriceFrom(&price)
		if err != nil {
			log.Println("Error with updating priceFrom in the priceFrom switch case ", err)
			msg.Text = "Не удается обновить стартовую цену ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		user.UserChoice.PriceFrom = &price
	case "priceto":
		val := update.Message.CommandArguments()
		if val == "" {

			user.PriceTo = nil
			err := db.UpdatePriceTo(nil)
			if err != nil {
				log.Println("Error with updating priceTo in the priceTo switch case ", err)
				msg.Text = "Не удается обновить сумму до ознакомьтесь с документацией по команде /help и повторите попытку"
			}

			return
		}
		price, err := trimAllSpaces(val)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdatePriceTo(&price)
		if err != nil {
			log.Println("Error with updating priceTo in the priceTo switch case ", err)
			msg.Text = "Не удается обновить сумму до ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		user.UserChoice.PriceTo = &price
	case "floorfrom":
		val := update.Message.CommandArguments()
		if val == "" {

			user.FloorFrom = nil
			err := db.UpdateFloorFrom(nil)
			if err != nil {
				log.Println("Error with updating floorFrom ", err)
				msg.Text = "Не удается обновить стартовый этаж ознакомьтесь с документацией по команде /help и повторите попытку"
			}

			return
		}
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateFloorFrom(&floor)
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

			user.FloorTo = nil
			err := db.UpdateFloorTo(nil)
			if err != nil {
				log.Println("Error with updating floorFrom ", err)
				msg.Text = "Не удается обновить этаж до  ознакомьтесь с документацией по команде /help и повторите попытку"
			}
			return
		}

		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateFloorTo(&floor)
		if err != nil {
			log.Println("Error with updating floorFrom ", err)
			msg.Text = "Не удается обновить этаж до  ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(floor)
		user.UserChoice.FloorTo = &curr
	case "notfirst":
		f := !tempC.notFirst
		user.UserChoice.CheckboxNotFirstFloor = &f
		err := db.UpdateNotFirstFloor(f)
		tempC.notFirst = f
		if err != nil {
			log.Println("Error with updating notFirst ", err)
			msg.Text = "Не удается обновить параметр не первый этаж ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "notlast":
		l := !tempC.notLast
		user.UserChoice.CheckboxNotLastFloor = &l
		err := db.UpdateNotLastFloor(l)
		tempC.notLast = l
		if err != nil {
			log.Println("Error with updating notLast ", err)
			msg.Text = "Не удается обновить параметр не последний этаж ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "fromowner":
		o := !tempC.fromOwner
		user.UserChoice.CheckboxFromOwner = &o
		err := db.UpdateFromOwner(o)
		tempC.fromOwner = o
		if err != nil {
			log.Println("Error with updating fromOwner ", err)
			msg.Text = "Не удается обновить параметр от владельца ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "newbuilding":
		n := !tempC.newBuilding
		user.UserChoice.CheckboxNewBuilding = &n
		err := db.UpdateNewBuilding(n)
		tempC.newBuilding = n
		if err != nil {
			log.Println("Error with updating newBuilding ", err)
			msg.Text = "Не удается обновить параметр новостройка ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "realestate":
		r := !tempC.realEstate
		user.UserChoice.CheckRealEstate = &r
		err := db.UpdateRealEstate(r)
		tempC.realEstate = r
		if err != nil {
			log.Println("Error with updating realEstate ", err)
			msg.Text = "Не удается обновить параметр от крыша агента ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "floorhfrom":
		val := update.Message.CommandArguments()
		if val == "" {

			user.FloorInTheHouseFrom = nil
			err := db.UpdateFloorInTheHouseFrom(nil)
			if err != nil {
				log.Println("Error with updating floorInTheHouseFrom ", err)
				msg.Text = "Не удается обновить параметр этаж в доме от ознакомьтесь с документацией по команде /help и повторите попытку"
			}

			return
		}
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateFloorInTheHouseFrom(&floor)
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

			user.FloorInTheHouseTo = nil
			err := db.UpdateFloorInTheHouseTo(nil)
			if err != nil {
				log.Println("Error with updating floorInTheHouseTo ", err)
				msg.Text = "Не удается обновить параметр этаж в доме до ознакомьтесь с документацией по команде /help и повторите попытку"
			}

			return
		}
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		err = db.UpdateFloorInTheHouseTo(&floor)
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

			user.AreaFrom = nil
			err := db.UpdateAreaFrom(nil)
			if err != nil {
				log.Println("Error with updating areaFrom ", err)
				msg.Text = "Не удается обновить параметр стартовая площадь ознакомьтесь с документацией по команде /help и повторите попытку"

			}
			return
		}
		user.UserChoice.AreaFrom = &val
		err := db.UpdateAreaFrom(&val)
		if err != nil {
			log.Println("Error with updating areaFrom ", err)
			msg.Text = "Не удается обновить параметр стартовая площадь ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "areato":
		val := update.Message.CommandArguments()
		if val == "" {

			user.AreaTo = nil
			err := db.UpdateAreaTo(nil)
			if err != nil {
				log.Println("Error with updating areaTo ", err)
				msg.Text = "Не удается обновить параметр площадь до ознакомьтесь с документацией по команде /help и повторите попытку"

			}
			return
		}
		user.UserChoice.AreaTo = &val
		err := db.UpdateAreaTo(&val)
		if err != nil {
			log.Println("Error with updating areaTo ", err)
			msg.Text = "Не удается обновить параметр площадь до ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "kitfrom":
		val := update.Message.CommandArguments()
		if val == "" {

			user.KitchenAreaFrom = nil
			err := db.UpdateKitchenFrom(nil)
			if err != nil {
				log.Println("Error with updating kitchen ", err)
				msg.Text = "Не удается обновить параметр площадь кухни от ознакомьтесь с документацией по команде /help и повторите попытку"
			}
			return
		}

		user.UserChoice.KitchenAreaFrom = &val
		err := db.UpdateKitchenFrom(&val)
		if err != nil {
			log.Println("Error with updating kitchen ", err)
			msg.Text = "Не удается обновить параметр площадь кухни от ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "kitto":
		val := update.Message.CommandArguments()
		if val == "" {

			user.KitchenAreaTo = nil
			err := db.UpdateKitchenTo(nil)
			if err != nil {
				log.Println("Error with updating kitchenTo ", err)
				msg.Text = "Не удается обновить параметр площадь кухни до ознакомьтесь с документацией по команде /help и повторите попытку"

			}
			return
		}
		user.UserChoice.KitchenAreaTo = &val
		err := db.UpdateKitchenTo(&val)
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
		b.sendMessage(msg)
		b.tickStart(msg, update, db, last)
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
	region := inline.CollectButtonData(inline.ChooseRegion(user.City))
	if db.GetCity() != "" {
		msg.Text = "Для изменения города используйте команду /city\nДля изменения региона используйте команду /region"
		b.sendMessage(&msg)
		return
	}
	if city[update.CallbackQuery.Data] {
		user.City = update.CallbackQuery.Data
		msg.Text = "Выберите один из районов города "
		tempC.buttonCity = user.City
		msg.ReplyMarkup = inline.ChooseRegion(user.City)
		b.sendMessage(&msg)
	} else if region[update.CallbackQuery.Data] {
		checkForSkip(update, user)
		user.City = update.CallbackQuery.Data
		db.IntroDataStartCommand(user)
		msg.Text = "Отлично теперь ознакомьтесь с командами в меню и настраивайте фильтры"
		b.sendMessage(&msg)
	}
}

func (b *Bot) handleRegion(update *tgbotapi.Update, user *models.User, db *postgres.Sql) {
	region := inline.ChooseRegion(user.City)
	data := inline.CollectButtonData(region)
	if data[update.CallbackQuery.Data] {
		checkForSkip(update, user)
		err := db.UpdateCity(update.CallbackQuery.Data)
		if err != nil {
			log.Println("Cannot update region in the handleRegion function")
			return
		}
	}
}

func (b *Bot) handleCity(update *tgbotapi.Update, user *models.User, db *postgres.Sql) {
	chatID := update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	city := inline.ChooseCity()
	msg.ReplyMarkup = city
	data := inline.CollectButtonData(city)
	if data[update.CallbackQuery.Data] {
		checkForSkip(update, user)
		tempC.buttonCity = update.CallbackQuery.Data
		err := db.UpdateCity(update.CallbackQuery.Data)
		if err != nil {
			log.Println("Cannot update region in the handleCity function")
			return
		}
	}
}

func checkForSkip(update *tgbotapi.Update, user *models.User) {
	if update.CallbackQuery.Data == "-" {
		update.CallbackQuery.Data = user.City
	}
}

func (b *Bot) sendMessage(msg *tgbotapi.MessageConfig) {
	_, err := b.bot.Send(msg)
	if err != nil {
		return
	}
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
	val.City = utils.RegionEnRun(val.City)
	js, err := utils.ConvertToJSONO(val)
	if err != nil {
		log.Println("Cannot convert struct to json all command ", err.Error())
	}
	var hmap map[string]interface{}
	if err = json.Unmarshal([]byte(js), &hmap); err != nil {
		log.Println("cannot convert to the map all command ", err.Error())
	}
	newM := utils.EnRusUser(hmap)
	var res strings.Builder
	for s, i := range newM {
		switch t := i.(type) {
		case string:
			res.WriteString(s + " : " + t + "\n")
		case []interface{}:
			var strs []string
			for _, item := range t {
				if str, ok := item.(string); ok {
					strs = append(strs, str)
				}
				// Add handling for other types within []interface{} if needed
			}
			jn := strings.Join(strs, ",")
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
}
func (b *Bot) tickStart(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, db *postgres.Sql, last *[]string) {
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
	fmt.Println(last)
	for range ticker.C {
		if (*last)[len(*last)-1] == "stop" {
			msg.Text = "Вы остановили поиск квартир"
			b.sendMessage(msg)
			break
		}
		fmt.Println("here in the ticker loop")
		houses, err := init.NewScrap(&dups)
		if err != nil {
			msg.Text = "Возникли проблемы с веб-сайтом Krisha.kz ожидайте..."
			break
		}

		b.rangeHouses(houses, update)
	}
}
