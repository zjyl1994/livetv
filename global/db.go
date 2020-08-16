package global

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/zjyl1994/livetv/model"
)

var DB *gorm.DB

func InitDB(filepath string) (err error) {
	DB, err = gorm.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&model.Config{}, &model.Channel{}).Error
	if err != nil {
		return err
	}
	for k, v := range defaultConfigValue {
		configItem := model.Config{Name: k, Data: v}
		if DB.NewRecord(configItem) {
			err = DB.Create(&configItem).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
