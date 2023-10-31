package utils

import (
	json2 "encoding/json"
	"fmt"
)

func ConvertToJSONO(value any) (string, error) {
	json, err := json2.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func EnRusHouse(m map[string]interface{}) map[string]interface{} {
	hmap := map[string]string{
		"link":                "Ссылка",
		"intro":               "Заголовок",
		"price":               "Цена",
		"city":                "Город",
		"house_type":          "Тип дома",
		"residential_complex": "Жилой комплекс",
		"year_of_build":       "Год постройки",
		"floor":               "Этаж",
		"area":                "Площадь",
		"bathroom":            "Санузел",
		"ceil":                "Потолки",
		"state":               "Состояние",
		"former_hostel":       "Бывшее общежитие",
	}
	translatedMap := make(map[string]interface{})

	for key, value := range m {
		translatedKey, ok := hmap[key]
		if !ok {
			translatedKey = key
		}
		translatedMap[translatedKey] = value
	}

	return translatedMap
}

func EnRusUser(m map[string]interface{}) map[string]interface{} {
	hmap := map[string]string{
		"buy_or_rent":              "Тип",
		"type_item":                "Недвижимость",
		"city":                     "Город",
		"rooms":                    "Кол-во комнат",
		"type_house":               "Типы домов",
		"year_of_built_from":       "Год постройки от",
		"year_of_built_to":         "Год постройки до",
		"price_from":               "Цена от",
		"price_to":                 "Цена до",
		"floor_from":               "Этаж от",
		"floor_to":                 "Этаж до",
		"checkbox_not_first_floor": "Не первый этаж",
		"checkbox_not_last_floor":  "Не последний этаж",
		"checkbox_from_owner":      "От хозяина",
		"checkbox_new_building":    "Новостройка",
		"check_real_estate":        "От крыши агента",
		"floor_in_the_house_from":  "Этаж в доме от",
		"floor_in_the_house_to":    "Этаж в доме до",
		"total_area":               "Площадь от",
		"area_to":                  "Площадь до",
		"kitchen_area":             "Площадь кухни от",
		"kitchen_area_to":          "Площадь кухни до",
	}
	fmt.Println(m)
	translatedMap := make(map[string]interface{})

	for key, value := range m {
		translatedKey, ok := hmap[key]
		if !ok {
			continue
		}
		translatedMap[translatedKey] = value
	}

	return translatedMap
}
