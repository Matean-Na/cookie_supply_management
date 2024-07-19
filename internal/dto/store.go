package dto

type StoreCreateDTO struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Contact     string `json:"contact"`
	PhoneNumber string `json:"phoneNumber"`
}

type StoreGetDTO struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Contact     string `json:"contact"`
	PhoneNumber string `json:"phoneNumber"`
	Debt        string `json:"debt"`
}
