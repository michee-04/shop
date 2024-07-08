package model

import (
	"github.com/google/uuid"
	"github.com/michee/e-commerce/database"
	"gorm.io/gorm"
)

type Hero struct {
	HeroId      string `gorm:"primary_key; column:hero_id"`
	Title       string `gorm:"column:title" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	ImageUrl    string `gorm:"column:image_url" json:"image_url"`
}


func init(){
	database.ConnectDB()
	Db = database.GetDB()
	// Db.Migrator().DropTable(&Hero{})
	// if Db != nil {
	// 	err := Db.AutoMigrate(&Hero{})
	// 	if err != nil {
	// 		panic("Failed to migrate hero model: "+ err.Error())
	// 	}
	// } else {
	// 	panic("DB connection is nil")
	// }
}

func (h *Hero) CreateHero() *Hero {
	h.HeroId = uuid.New().String()
	Db.Create(&h)
	return h
}

func GetHero() []Hero{
	var h []Hero
	Db.Find(&h)
	return h
}

func GetHeroById(Id string) (*Hero, *gorm.DB) {
	var h Hero
	db := Db.Where("hero_id=?", Id).First(&h)
	if db.Error != nil {
		return nil, db
	}
	return &h, db
}

func DeleteHero(Id string) Hero {
	var h Hero
	Db.Where("hero_id=?", Id).Delete(&h)
	return h
}