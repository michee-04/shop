package model

import (
	"github.com/google/uuid"
	"github.com/michee/e-commerce/database"
	"gorm.io/gorm"
)

type Banner struct {
	BannerId string `gorm:"primary_key;column:banner_id"`
	Title    string `gorm:"column:title" json:"title"`
	BgColor  string `gorm:"column:bg_color" json:"bg_color"`
}

func init() {
	database.ConnectDB()
	Db = database.GetDB()
	Db.Migrator().DropTable(&Banner{})
	if Db != nil {
		err := Db.AutoMigrate(&Banner{})
		if err != nil {
			panic("Failed to migrate Banner model: " + err.Error())
		}
	} else {
		panic("DB connection is nil")
	}
}

func (b *Banner) CreateBanner() *Banner {
	b.BannerId = uuid.New().String()
	Db.Create(b)
	return b
}

func GetBanner() []Banner {
	var b []Banner
	Db.Find(&b)
	return b
}

func GetBannerById(Id string) (*Banner, *gorm.DB) {
	var b Banner
	db := Db.Where("banner_id=?", Id).First(&b)
	if db.Error != nil {
		return nil, db
	}

	return &b, db
}


func DeleteBanner(Id string) Banner{
	var b Banner
	Db.Where("banner_id=?", Id).Delete(&b)
	return b
}