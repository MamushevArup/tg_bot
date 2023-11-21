package models

import (
	uuid2 "github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID       uuid2.UUID `db:"id"`
	Username string     `json:"username,omitempty"`
	UserChoice
}
type UserChoice struct {
	City                  string         `db:"city" json:"city,omitempty,omitempty"`
	Rooms                 pq.StringArray `db:"rooms" json:"rooms,omitempty"`
	TypeHouse             pq.StringArray `db:"typehouse" json:"type_house,omitempty"`
	YearOfBuiltFrom       *uint          `db:"yearbuiltfrom" json:"year_of_built_from,omitempty"`
	YearOfBuiltTo         *uint          `db:"yearbuiltto" json:"year_of_built_to,omitempty"`
	PriceFrom             *uint64        `db:"pricefrom" json:"price_from,omitempty"`
	PriceTo               *uint64        `db:"priceto" json:"price_to,omitempty"`
	FloorFrom             *uint          `db:"floorfrom" json:"floor_from,omitempty"`
	FloorTo               *uint          `db:"floorto" json:"floor_to,omitempty"`
	CheckboxNotFirstFloor *bool          `db:"notfirstfloor" json:"checkbox_not_first_floor,omitempty"`
	CheckboxNotLastFloor  *bool          `db:"notlastfloor" json:"checkbox_not_last_floor,omitempty"`
	CheckboxFromOwner     *bool          `db:"fromowner" json:"checkbox_from_owner,omitempty"`
	CheckboxNewBuilding   *bool          `db:"newbuilding" json:"checkbox_new_building,omitempty"`
	CheckRealEstate       *bool          `db:"realestate" json:"check_real_estate,omitempty"`
	FloorInTheHouseFrom   *uint          `db:"floorinthehousefrom" json:"floor_in_the_house_from,omitempty"`
	FloorInTheHouseTo     *uint          `db:"floorinthehouseto" json:"floor_in_the_house_to,omitempty"`
	AreaFrom              *string        `db:"areafrom" json:"total_area,omitempty"`
	AreaTo                *string        `db:"areato" json:"area_to,omitempty"`
	KitchenAreaFrom       *string        `db:"kitchenfrom" json:"kitchen_area,omitempty"`
	KitchenAreaTo         *string        `db:"kitchento" json:"kitchen_area_to,omitempty"`
	Running               *bool          `db:"running" json:"running,omitempty"`
}
