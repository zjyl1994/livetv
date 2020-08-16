package global

import (
	"github.com/jinzhu/gorm"
	"github.com/zjyl1994/livetv/model"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func InitDB(filepath string) (err error) {
	DB, err = gorm.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&model.Config{}, &model.Channel{}).Error
}
