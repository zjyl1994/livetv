package model

type Config struct {
	Name string `gorm:"primary_key"`
	Data string
}
