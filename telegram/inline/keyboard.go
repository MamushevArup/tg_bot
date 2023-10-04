package inline

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func BuyOrRent() tgbotapi.InlineKeyboardMarkup {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Продажа", "prodazha"),
			tgbotapi.NewInlineKeyboardButtonData("Аренда", "arenda"),
		),
	)
	return keyboard
}
func TypeItem() tgbotapi.InlineKeyboardMarkup {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Квартиры", "kvartiry"),
			tgbotapi.NewInlineKeyboardButtonData("Дома", "doma"),
			tgbotapi.NewInlineKeyboardButtonData("Дачи", "dachi"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Бизнес", "biznes"),
			tgbotapi.NewInlineKeyboardButtonData("Коммерческая", "kommercheskaya-nedvizhimost"),
		),
	)
	return numericKeyboard
}
