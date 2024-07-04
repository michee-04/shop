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


func CreateBanner(w http.ResponseWriter, r *http.Request) {
	banner := model.Banner{}
	utils.ParseBody(r, &banner)
	b := banner.CreateBanner()
	res, _ := json.Marshal(b)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBanner(w http.ResponseWriter, r *http.Request) {
	b := model.GetBanner()
	res, _ := json.Marshal(b)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBannerById(w http.ResponseWriter, r *http.Request) {
	bannerId := chi.URLParam(r, "bannerId")
	b, _ := model.GetBannerById(bannerId)
	res, _ := json.Marshal(b)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateBanner(w http.ResponseWriter, r *http.Request) {
	bannerUpdate := &model.Banner{}
	utils.ParseBody(r, &bannerUpdate)

	bannerId := chi.URLParam(r, "bannerId")

	var b model.Banner
	err := model.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("banner_id=?", bannerId).First(&b).Error; err != nil {
			return err
		}

		if bannerUpdate.BgColor != "" {
			b.BgColor = bannerUpdate.BgColor
		}
		if bannerUpdate.Title != "" {
			b.Title = bannerUpdate.Title
		}
		

		return tx.Save(&b).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "banner not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update banner: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res, err := json.Marshal(b)
	if err != nil {
		http.Error(w, "Failed update hero", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteBanner(w http.ResponseWriter, r *http.Request) {
	bannerId := chi.URLParam(r, "bannerId")
	b := model.DeleteBanner(bannerId)
	
	utils.RespondWithJSON(w, http.StatusOK, "Delete Banner successuful", b)
}