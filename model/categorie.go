package model

type Categorie struct {
	CategorieId string `gorm:"primary_key; column:categorie_id"`
	Title       string `gorm:"column:title" json:"title"`
	Article []Article `gorm:"foreignKey:categorie_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
