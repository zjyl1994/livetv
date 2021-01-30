package database

import (
	"errors"

	"gorm.io/gorm"
)

func GetConfig(key string) (string, error) {
	if confValue, ok := configCache.Load(key); ok {
		return confValue.(string), nil
	} else {
		var value Config
		err := db.Where("name = ?", key).First(&value).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", ErrConfigNotFound
			} else {
				return "", err
			}
		} else {
			configCache.Store(key, value.Data)
			return value.Data, nil
		}
	}
}

func SetConfig(key, value string) error {
	data := Config{Name: key, Data: value}
	err := db.Save(&data).Error
	if err == nil {
		configCache.Store(key, value)
	}
	return err
}
