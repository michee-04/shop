package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/michee/e-commerce/model"
	"github.com/michee/e-commerce/utils"
	"gorm.io/gorm"
)

func CreateCategorie(w http.ResponseWriter, r *http.Request) {
	categorie := model.Categorie{}
	utils.ParseBody(r, categorie)
	c := categorie.CreateCategorie()
	res, _ := json.Marshal(c)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetCategorie(w http.ResponseWriter, r *http.Request) {
	c := model.GetCategorie()
	res, _ := json.Marshal(c)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetCategorieId(w http.ResponseWriter, r *http.Request) {
	categorieId := chi.URLParam(r, "categorieId")
	c, _ := model.GetCategorieId(categorieId)
	res, _ := json.Marshal(c)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateCategorie(w http.ResponseWriter, r *http.Request) {
	categorieUpdate := &model.Categorie{}
	utils.ParseBody(r, &categorieUpdate)

	categorieId := chi.URLParam(r, "categorieId")

	var c model.Categorie
	err := model.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("categorie_id=?", categorieId).First(&c).Error; err != nil {
			return err
		}

		if categorieUpdate.Title != "" {
			c.Title = categorieUpdate.Title
		}

		return tx.Save(&c).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "hero not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update hero: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res, err := json.Marshal(c)
	if err != nil {
		http.Error(w, "Failed update hero", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteCategorie(w http.ResponseWriter, r *http.Request) {
	categorieId := chi.URLParam(r, "categorieId")
	c := model.DeleteCategorie(categorieId)

	utils.RespondWithJSON(w, http.StatusOK, "Delete Categorie successuful", c)
}
