package model

type Claims struct {
	Id           uint   `json:"id" gorm:"primary_key"`
	LoginModelId uint   `json:"loginModelId"`
	Type         string `json:"type"`
	Value        string `json:"value"`
}
