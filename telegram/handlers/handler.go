package handlers

import (
	"github.com/MamushevArup/krisha-scraper/database/postgres"
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/MamushevArup/krisha-scraper/telegram/inline"
	"github.com/MamushevArup/krisha-scraper/utils/files"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
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
		msg.ReplyMarkup = inline.BuyOrRent()
	case "help":
		val, err := files.ReadTXT("utils/text/txt/text.txt")
		if err != nil {
			log.Println("Error with something", err)
		}
		msg.Text = val
	case "city":
		val := update.Message.CommandArguments()
		if val == "" {
			msg.Text = "Введите город\nПример /city Алматы"
			return
		}
		flag := false
		for k, v := range listCities() {
			if k == val {
				flag = true
				val = v
			}
		}
		if !flag {
			msg.Text = "Возможно вы неправильно ввели название города или на данный момент этот город недоступен. \nПовторите попытку"
			return
		}
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
	case "builtFrom":
		val := update.Message.CommandArguments()
		if val == "" {
			user.YearOfBuiltFrom = nil
			err := db.UpdateBuiltFrom(0)
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
		err = db.UpdateBuiltFrom(uint(year))
		if err != nil {
			log.Println("Error with updating builtFrom in the builtFrom switch case ", err)
			msg.Text = "Не удается обновить год постройки от ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(year)
		user.UserChoice.YearOfBuiltFrom = &curr
	case "builtTo":
		val := update.Message.CommandArguments()
		if val == "" {
			user.YearOfBuiltTo = nil
			err := db.UpdateBuiltTo(0)
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
		err = db.UpdateBuiltTo(uint(year))
		if err != nil {
			log.Println("Error with updating builtTo in the builtTo switch case ", err)
			msg.Text = "Не удается обновить год постройки до  ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(year)
		user.UserChoice.YearOfBuiltTo = &curr

	case "priceFrom":
		val := update.Message.CommandArguments()
		if val == "" {
			user.PriceFrom = nil
			err := db.UpdatePriceFrom(0)
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
		err = db.UpdatePriceFrom(price)
		if err != nil {
			log.Println("Error with updating priceFrom in the priceFrom switch case ", err)
			msg.Text = "Не удается обновить стартовую цену ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		user.UserChoice.PriceFrom = &price
	case "priceTo":
		val := update.Message.CommandArguments()
		if val == "" {
			user.PriceTo = nil
			err := db.UpdatePriceTo(0)
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
		err = db.UpdatePriceTo(price)
		if err != nil {
			log.Println("Error with updating priceTo in the priceTo switch case ", err)
			msg.Text = "Не удается обновить сумму до ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		user.UserChoice.PriceTo = &price
	case "floorFrom":
		val := update.Message.CommandArguments()
		if val == "" {
			user.FloorFrom = nil
			err := db.UpdateFloorFrom(0)
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
		err = db.UpdateFloorFrom(floor)
		if err != nil {
			log.Println("Error with updating floorFrom ", err)
			msg.Text = "Не удается обновить стартовый этаж ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(floor)
		user.UserChoice.FloorFrom = &curr
	case "floorTo":
		val := update.Message.CommandArguments()
		if val == "" {
			user.FloorTo = nil
			err := db.UpdateFloorTo(0)
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
		err = db.UpdateFloorTo(floor)
		if err != nil {
			log.Println("Error with updating floorFrom ", err)
			msg.Text = "Не удается обновить этаж до  ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(floor)
		user.UserChoice.FloorTo = &curr
	case "notFirst":
		flag := true
		user.UserChoice.CheckboxNotFirstFloor = &flag
		err := db.UpdateNotFirstFloor(true)
		if err != nil {
			log.Println("Error with updating notFirst ", err)
			msg.Text = "Не удается обновить параметр не первый этаж ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "notLast":
		flag := true
		user.UserChoice.CheckboxNotLastFloor = &flag
		err := db.UpdateNotLastFloor(true)
		if err != nil {
			log.Println("Error with updating notLast ", err)
			msg.Text = "Не удается обновить параметр не последний этаж ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "fromOwner":
		flag := true
		user.UserChoice.CheckboxFromOwner = &flag
		err := db.UpdateFromOwner(true)
		if err != nil {
			log.Println("Error with updating fromOwner ", err)
			msg.Text = "Не удается обновить параметр от владельца ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "newBuilding":
		flag := true
		user.UserChoice.CheckboxNewBuilding = &flag
		err := db.UpdateNewBuilding(true)
		if err != nil {
			log.Println("Error with updating newBuilding ", err)
			msg.Text = "Не удается обновить параметр новостройка ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "realEstate":
		flag := true
		user.UserChoice.CheckRealEstate = &flag
		err := db.UpdateRealEstate(true)
		if err != nil {
			log.Println("Error with updating realEstate ", err)
			msg.Text = "Не удается обновить параметр от крыша агента ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
	case "floorHouseFrom":
		val := update.Message.CommandArguments()
		if val == "" {
			user.FloorInTheHouseFrom = nil
			err := db.UpdateFloorInTheHouseFrom(0)
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
		err = db.UpdateFloorInTheHouseFrom(floor)
		if err != nil {
			log.Println("Error with updating floorInTheHouseFrom ", err)
			msg.Text = "Не удается обновить параметр этаж в доме от ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(floor)
		user.UserChoice.FloorInTheHouseFrom = &curr
	case "floorHouseTo":
		val := update.Message.CommandArguments()
		if val == "" {
			user.FloorInTheHouseTo = nil
			err := db.UpdateFloorInTheHouseTo(0)
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
		err = db.UpdateFloorInTheHouseTo(floor)
		if err != nil {
			log.Println("Error with updating floorInTheHouseTo ", err)
			msg.Text = "Не удается обновить параметр этаж в доме до ознакомьтесь с документацией по команде /help и повторите попытку"
			return
		}
		curr := uint(floor)
		user.UserChoice.FloorInTheHouseTo = &curr
	case "areaFrom":
		val := update.Message.CommandArguments()
		if val == "" {
			user.AreaFrom = nil
			err := db.UpdateAreaFrom("")
			if err != nil {
				log.Println("Error with updating areaFrom ", err)
				msg.Text = "Не удается обновить параметр стартовая площадь ознакомьтесь с документацией по команде /help и повторите попытку"
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
	case "areaTo":
		val := update.Message.CommandArguments()
		if val == "" {
			user.AreaTo = nil
			err := db.UpdateAreaTo("")
			if err != nil {
				log.Println("Error with updating areaTo ", err)
				msg.Text = "Не удается обновить параметр площадь до ознакомьтесь с документацией по команде /help и повторите попытку"
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
	case "kitFrom":
		val := update.Message.CommandArguments()
		if val == "" {
			user.KitchenAreaFrom = nil
			err := db.UpdateKitchenFrom("")
			if err != nil {
				log.Println("Error with updating kitchen ", err)
				msg.Text = "Не удается обновить параметр площадь кухни от ознакомьтесь с документацией по команде /help и повторите попытку"
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
	case "kitTo":
		val := update.Message.CommandArguments()
		if val == "" {
			user.KitchenAreaTo = nil
			err := db.UpdateKitchenTo("")
			if err != nil {
				log.Println("Error with updating kitchenTo ", err)
				msg.Text = "Не удается обновить параметр площадь кухни до ознакомьтесь с документацией по команде /help и повторите попытку"
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
	case "ch":
		db.Insert()
	case "run":
		msg.Text = "Идет поиск подходящих обьявлений\nТелеграм оповестит вас когда мы найдем подходящие обьявления\nДля этого уберите беззвучный режим."
		go func() {
			b.tickStart(msg, update, db)
		}()
	case "stop":
		msg.Text = "Вы остановили поиск обьявлений теперь можете изменить настройки\nПо комманде /help смотрите все доступные комманды"
	default:
		msg.Text = "Нет такой комманды. Ознакомьтесь с документацией"
	}
}
