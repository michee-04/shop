package model

import (
	"github.com/google/uuid"
	"github.com/michee/e-commerce/database"
	"gorm.io/gorm"
)

type Categorie struct {
	CategorieId string    `gorm:"primary_key; column:categorie_id"`
	Title       string    `gorm:"column:title" json:"title"`
	Article     []Article `gorm:"foreignKey:categorie_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func InitCategroie() {
	database.ConnectDB()
	Db = database.GetDB()
	// Db.Migrator().DropTable()
	if Db != nil {
		err := Db.AutoMigrate(&Categorie{})
		if err != nil {
			panic("Failed to migrate Categorie model: " + err.Error())
		}
	} else {
		panic("DB connection is nil")
	}
}

func (c *Categorie) CreateCategorie() *Categorie {
	c.CategorieId = uuid.New().String()
	Db.Create(c)
	return c
}

func GetCategorie() []Categorie {
	var c []Categorie
	Db.Find(&c)
	return c
}

func GetCategorieId(Id string) (*Categorie, *gorm.DB) {
	var c Categorie
	db := Db.Where("categorie_id=?", Id).First(&c)
	if db.Error != nil {
		return nil, db
	}

	return &c, db
}

func DeleteCategorie(Id string) Categorie {
	var c Categorie
	Db.Where("categorie_id=?", Id).Delete(&c)
	return c
}
