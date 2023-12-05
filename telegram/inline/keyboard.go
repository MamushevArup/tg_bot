package inline

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ChooseCity() tgbotapi.InlineKeyboardMarkup {
	var key = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Алматы", "almaty"),
			tgbotapi.NewInlineKeyboardButtonData("Астана", "astana"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Шымкент", "shymkent"),
		),
	)
	return key
}
func ChooseRegion(city string) tgbotapi.InlineKeyboardMarkup {
	var almaty = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ауэзовский", "almaty-aujezovskij"),
			tgbotapi.NewInlineKeyboardButtonData("Алатауский", "almaty-alatauskij"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Алмалинский", "almaty-almalinskij"),
			tgbotapi.NewInlineKeyboardButtonData("Бостандыкский", "almaty-bostandykskij"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Жетысуский", "almaty-zhetysuskij"),
			tgbotapi.NewInlineKeyboardButtonData("Медеуский", "almaty-medeuskij"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Наурызбайский", "almaty-nauryzbajskiy"),
			tgbotapi.NewInlineKeyboardButtonData("Турксибский", "almaty-turksibskij"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пропустить", "-"),
		),
	)
	var astana = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Алматы", "astana-almatinskij"),
			tgbotapi.NewInlineKeyboardButtonData("Есильский", "astana-esilskij"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Нура", "astana-nura"),
			tgbotapi.NewInlineKeyboardButtonData("Байконур", "r-n-bajkonur"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сарыарка", "astana-saryarkinskij"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пропустить", "-"),
		),
	)
	var shymkent = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Абайский", "shymkent-abajskij"),
			tgbotapi.NewInlineKeyboardButtonData("Аль-Фарабийский", "shymkent-al-farabijskij"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Енбекшинский", "shymkent-enbekshinskij"),
			tgbotapi.NewInlineKeyboardButtonData("Туран", "shymkent-turan"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Каратауский", "karatauskij"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пропустить", "-"),
		),
	)
	hmap := map[string]tgbotapi.InlineKeyboardMarkup{
		"almaty":   almaty,
		"astana":   astana,
		"shymkent": shymkent,
	}
	return hmap[city]
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
