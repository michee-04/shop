package model

import (
	"github.com/google/uuid"
	"github.com/michee/e-commerce/database"
	"gorm.io/gorm"
)

type Article struct {
	ArticleId   string    `gorm:"primary_key; column:article_id"`
	Title       string    `gorm:"column:title" json:"title"`
	Description string    `gorm:"column:description" json:"description"`
	Color       string    `gorm:"column:color" json:"color"`
	Size        string    `gorm:"column:size" json:"size"`
	CategorieId string    `gorm:"column:categorie_id" json:"categorie_id"`
	Categorie   Categorie `gorm:"foreignKey:categorie_id" json:"-"`
}

func InitArticle() {
	database.ConnectDB()
	Db = database.GetDB()
	// Db.Migrator().DropTable(&Article{})
	if Db != nil {
		err := Db.AutoMigrate(&Article{})
		if err != nil {
			panic("Failed to migrate Article Model: " + err.Error())
		}
	} else {
		panic("DB connection is nil")
	}
}

func (a *Article) CreateArticle() *Article {
	a.ArticleId = uuid.New().String()
	Db.Create(a)
	return a
}

func GetArticle() []Article {
	var a []Article
	Db.Preload("Categorie").Find(&a)
	return a
}

func GetArticleById(Id string) (*Article, *gorm.DB) {
	var a Article
	db := Db.Preload("Categorie").Where("article_id=?", Id).First(&a)
	if db.Error != nil {
		return nil, db
	}

	return &a, db
}

func DeleteArticle(Id string) Article {
	var a Article
	Db.Where("article_id=?", Id).Delete(&a)
	return a
}
