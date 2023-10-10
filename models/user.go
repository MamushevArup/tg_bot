package models

type User struct {
	Username string `json:"username"`
	UserChoice
}
type UserChoice struct {
	BuyOrRent             string   `json:"buy_or_rent"`
	TypeItem              string   `json:"type_item"`
	City                  string   `json:"city,omitempty"`
	Rooms                 []string `json:"rooms,omitempty"`
	TypeHouse             []string `json:"type_house,omitempty"`
	YearOfBuiltFrom       uint     `json:"year_of_built_from,omitempty"`
	YearOfBuiltTo         uint     `json:"year_of_built_to,omitempty"`
	PriceFrom             uint64   `json:"price_from,omitempty"`
	PriceTo               uint64   `json:"price_to,omitempty"`
	FloorFrom             uint8    `json:"floor_from,omitempty"`
	FloorTo               uint8    `json:"floor_to,omitempty"`
	CheckboxNotFirstFloor bool     `json:"checkbox_not_first_floor,omitempty"`
	CheckboxNotLastFloor  bool     `json:"checkbox_not_last_floor,omitempty"`
	CheckboxFromOwner     bool     `json:"checkbot_from_owner,omitempty"`
	CheckboxNewBuilding   bool     `json:"checkbox_new_building,omitempty"`
	CheckRealEstate       bool     `json:"check_real_estate,omitempty"`
	FloorInTheHouseFrom   uint8    `json:"floor_in_the_house_from,omitempty"`
	FloorInTheHouseTo     uint8    `json:"floor_in_the_house_to,omitempty"`
	AreaFrom              string   `json:"total_area,omitempty"`
	AreaTo                string   `json:"area_to,omitempty"`
	KitchenAreaFrom       string   `json:"kitchen_area,omitempty"`
	KitchenAreaTo         string   `json:"kitchen_area_to,omitempty"`
}
