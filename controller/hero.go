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

func CreateHero(w http.ResponseWriter, r *http.Request) {
	hero := model.Hero{}
	utils.ParseBody(r, hero)
	h := hero.CreateHero()
	res, _ := json.Marshal(h)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetHero(w http.ResponseWriter, r *http.Request) {
	h := model.GetHero()
	res, _ := json.Marshal(h)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetHeroById(w http.ResponseWriter, r *http.Request) {
	heroId := chi.URLParam(r, "heroId")
	h, _ := model.GetHeroById(heroId)
	res, _ := json.Marshal(h)
	w.Header().Set("content-type", "application?json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateHero(w http.ResponseWriter, r *http.Request) {
	heroUpdate := &model.Hero{}
	utils.ParseBody(r, &heroUpdate)

	heroId := chi.URLParam(r, "heroId")

	var h model.Hero
	err := model.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("hero_id=?", heroId).First(&h).Error; err != nil {
			return err
		}

		if heroUpdate.Title != "" {
			h.Title = heroUpdate.Title
		}
		if heroUpdate.Description != "" {
			h.Description = heroUpdate.Description
		}
		if heroUpdate.ImageUrl != "" {
			h.ImageUrl = heroUpdate.ImageUrl
		}

		return tx.Save(&h).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "hero not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update hero: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res, err := json.Marshal(h)
	if err != nil {
		http.Error(w, "Failed update hero", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteHero(w http.ResponseWriter, r *http.Request) {
	heroId := chi.URLParam(r, "heroId")
	h := model.DeleteHero(heroId)
	
	utils.RespondWithJSON(w, http.StatusOK, "Delete hero successuful", h)
}