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

func (b *Bot) HandleUpdate(update *tgbotapi.Update, user *models.User, sentSecondInlineKeyboard, cityChecker map[int64]bool, db *postgres.Sql) {
	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			b.handleMessageCommand(update, user, cityChecker, db)
		}
		user.Username = update.Message.Chat.UserName
	} else if update.CallbackQuery != nil {
		b.handleCallbackQuery(update, user, sentSecondInlineKeyboard)
	}
	fmt.Println(user)
}

func (b *Bot) handleMessageCommand(update *tgbotapi.Update, user *models.User, cityChecker map[int64]bool, db *postgres.Sql) {

	if update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			b.handleCommand(&msg, update, user, cityChecker, db)
			b.sendMessage(&msg)
		}
	}
}

func (b *Bot) handleCommand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, user *models.User, cityChecker map[int64]bool, db *postgres.Sql) {
	//chatid := update.Message.Chat.ID
	command := update.Message.Command()
	stopTicker := make(chan struct{})
	switch command {
	case "start":
		msg.Text = "Рад видеть тебя здесь.\nДля начала тебе нужно выбрать следущее\nПосле этого ознакомьтесь с документацие по команде /help"
		msg.ReplyMarkup = inline.BuyOrRent()
	case "help":
		val, err := files.ReadTXT("utils/texts/text.txt")
		if err != nil {
			log.Println("Error with something", err)
		}
		msg.Text = val
	case "city":
		val := update.Message.CommandArguments()
		flag := false
		for k, v := range listCities() {
			if k == val {
				flag = true
				val = v
			}
		}
		if !flag {
			msg.Text = "Возможно вы неправильно ввели название города или на данный момент этот город недоступен. \nПовторите попытку"
		}
		user.City = val
		msg.Text = "Отлично!\nТеперь вы можете запустить поиск командой /run или продолжите настройку\nДоступные команды доступны /help"
	case "rooms":
		val := update.Message.CommandArguments()
		arrOfRooms := strings.Split(val, ",")
		for _, s := range arrOfRooms {
			curr, err := strconv.Atoi(s)
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
		msg.Text = "Отлично!\nТеперь вы можете запустить поиск командой /run или продолжите настройку\nДоступные команды доступны /help"
	case "builtf":
		val := update.Message.CommandArguments()
		year, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		user.UserChoice.YearOfBuiltFrom = uint(year)
	case "builtt":
		val := update.Message.CommandArguments()
		year, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		user.UserChoice.YearOfBuiltTo = uint(year)

	case "pricef":
		val := update.Message.CommandArguments()

		price, err := trimAllSpaces(val)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		user.UserChoice.PriceFrom = price
	case "pricet":
		val := update.Message.CommandArguments()
		price, err := trimAllSpaces(val)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		user.UserChoice.PriceTo = price
	case "floorf":
		val := update.Message.CommandArguments()
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		user.UserChoice.FloorFrom = uint(floor)
	case "floort":
		val := update.Message.CommandArguments()
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		user.UserChoice.FloorTo = uint(floor)
	case "checknf":
		user.UserChoice.CheckboxNotFirstFloor = true
	case "checknl":
		user.UserChoice.CheckboxNotLastFloor = true
	case "checkfo":
		user.UserChoice.CheckboxFromOwner = true
	case "checknb":
		user.UserChoice.CheckboxNewBuilding = true
	case "checkre":
		user.UserChoice.CheckRealEstate = true
	case "floorhf":
		val := update.Message.CommandArguments()
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		user.UserChoice.FloorInTheHouseFrom = uint(floor)
	case "floorht":
		val := update.Message.CommandArguments()
		floor, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			log.Println("Cannot convert string to int ", val+" "+err.Error())
			msg.Text = "Кажется вы ввели не подходящее число"
			return
		}
		user.UserChoice.FloorInTheHouseTo = uint(floor)
	case "areaf":
		val := update.Message.CommandArguments()
		user.UserChoice.AreaFrom = val
	case "areat":
		val := update.Message.CommandArguments()
		user.UserChoice.AreaTo = val
	case "kitchenf":
		val := update.Message.CommandArguments()
		user.UserChoice.KitchenAreaFrom = val
	case "kitchent":
		val := update.Message.CommandArguments()
		user.UserChoice.KitchenAreaTo = val
	case "run":
		// Here I need to implement the logic of taking data from the database for the user and
		// check is he/she fill the data that is required by a start command. Do it later
		if user.TypeItem == "" && user.BuyOrRent == "" {
			msg.Text = "Пожалуйста заполните необходимую информацию по команде /start"
			return
		}
		db.GetUser(user)
		c := colly.NewCollector(
			colly.AllowURLRevisit(),
		)
		c.SetRequestTimeout(5 * time.Second)
		init := scrap.New(c, user)
		ticker := time.NewTicker(3 * time.Second)
		go func() {
			select {
			case <-ticker.C:
				for range ticker.C {

					houses, err := init.NewScrap()
					_ = houses
					if err != nil {
						msg.Text = "Возникли проблемы с веб-сайтом Krisha.kz ожидайте..."
						break
					}
					msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Success")
					b.sendMessage(&msg2)
					for _, house := range *houses {
						val, err := utils.ConvertToJSON(house)
						fieldOrder := sliceOrder()
						var hmap map[string]interface{}

						if err != nil {
							log.Println("Cannot convert one element to the json ", err.Error())
						}
						if err = json.Unmarshal([]byte(val), &hmap); err != nil {
							log.Println("Error with converting json to the map")
						}
						output := orderOutput(fieldOrder, hmap)

						msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, output)
						b.sendMessage(&msg2)
					}
				}
			case <-stopTicker:
				ticker.Stop()
				return
			}
		}()
	case "config":
		close(stopTicker)
		msg.Text = "End of the loop in the ticker"
	}
}
func sliceOrder() []string {
	fieldOrder := []string{
		"link",
		"intro",
		"price",
		"city",
		"house_type",
		"residential_complex",
		"year_of_build",
		"area",
		"bathroom",
		"ceil",
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

func (b *Bot) handleCallbackQuery(update *tgbotapi.Update, user *models.User, sentSecondInlineKeyboard map[int64]bool) {
	chatID := update.CallbackQuery.Message.Chat.ID
	if !sentSecondInlineKeyboard[chatID] {
		b.sendSecondInlineKeyboard(chatID)
		sentSecondInlineKeyboard[chatID] = true
		user.BuyOrRent = update.CallbackQuery.Data
	}
	user.TypeItem = update.CallbackQuery.Data
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
