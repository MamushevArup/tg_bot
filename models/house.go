package models

type House struct {
	Intro              string `json:"intro"`
	Price              uint64 `json:"price"`
	City               string `json:"city"`
	HouseType          string `json:"house_type,omitempty"`
	ResidentialComplex string `json:"residential_complex,omitempty"`
	YearOfBuild        uint16 `json:"year_of_build,omitempty"`
	Floor              string `json:"floor,omitempty"`
	Area               string `json:"area,omitempty"`
	Bathroom           string `json:"bathroom,omitempty"`
	Ceil               string `json:"ceil,omitempty"`
	FormerHostel       string `json:"former_hostel,omitempty"`
	State              string `json:"state,omitempty"`
}
