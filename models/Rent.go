package models

import "github.com/lib/pq"

type BuyHouse struct {
	Residential   *string        `json:"residental"`
	AreaFrom      *int           `json:"area_from"`
	AreaTo        *int           `json:"area_to"`
	Rooms         pq.StringArray `json:"rooms"`
	PriceFrom     *int64         `json:"price_from"`
	PriceTo       *int64         `json:"price_to"`
	FloorFrom     *int           `json:"floor_from"`
	FloorTo       *int           `json:"floor_to"`
	Furnished     pq.StringArray `json:"furnished"`
	FromOwner     *bool          `json:"from_owner"`
	NotFirstFloor *bool          `json:"not_first_floor"`
	NotLastFloor  *bool          `json:"not_last_floor"`
	AllowAnimal   *bool          `json:"allow_animal"`
	AllowKids     *bool          `json:"allow_kids"`
	Text          *string        `json:"text"`
}
