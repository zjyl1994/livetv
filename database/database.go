package database

import (
	"errors"
	"sync"

	"github.com/zjyl1994/livetv/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db                 *gorm.DB
	configCache        sync.Map
	ErrConfigNotFound  = errors.New("config not found")
	defaultConfigValue = map[string]string{
		"base_url": "http://127.0.0.1:9000",
		"password": "password",
	}
)

func Init() (err error) {
	db, err = gorm.Open(sqlite.Open(utils.DataDir("./livetv.db")), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Config{}, &Channel{})
	if err != nil {
		return err
	}
	for key, valueDefault := range defaultConfigValue {
		var valueInDB Config
		err = db.Where("name = ?", key).First(&valueInDB).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				configCache.Store(key, valueDefault)
			} else {
				return err
			}
		} else {
			configCache.Store(key, valueInDB.Data)
		}
	}
	return nil
}
