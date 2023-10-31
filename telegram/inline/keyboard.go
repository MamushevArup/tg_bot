package inline

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
func ChooseCity() tgbotapi.InlineKeyboardMarkup {
	var key = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Алматы", "almaty"),
			tgbotapi.NewInlineKeyboardButtonData("Астана", "astana"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Шымкент", "shymkent"),
			tgbotapi.NewInlineKeyboardButtonData("Актау", "aktau"),
		),
	)
	return key
}
func CollectButtonData(v tgbotapi.InlineKeyboardMarkup) map[string]bool {
	set := make(map[string]bool, len(v.InlineKeyboard))
	for _, buttons := range v.InlineKeyboard {
		for _, button := range buttons {
			set[*button.CallbackData] = true
		}
	}
	return set
}
