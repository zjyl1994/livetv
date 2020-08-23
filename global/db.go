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
	for key, valueDefault := range defaultConfigValue {
		var valueInDB model.Config
		err = DB.Where("name = ?", key).First(&valueInDB).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				ConfigCache.Store(key, valueDefault)
			} else {
				return err
			}
		} else {
			ConfigCache.Store(key, valueInDB.Data)
		}
	}
	return nil
}
