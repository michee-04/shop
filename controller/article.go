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


func CreateArticle(w http.ResponseWriter, r *http.Request) {
	article := model.Article{}
	utils.ParseBody(r, &article)
	a := article.CreateArticle()
	res, _ := json.Marshal(a)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	c := model.GetArticle()
	res, _ := json.Marshal(c)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetArticleId(w http.ResponseWriter, r*http.Request) {
	articleId := chi.URLParam(r, "articleId")
	a, _ := model.GetArticleById(articleId)
	res, _ := json.Marshal(a)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	articleUpdate := &model.Article{}
	utils.ParseBody(r, &articleUpdate)

	articleId := chi.URLParam(r, "articleId")

	var a model.Article
	err := model.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id=?", articleId).First(&a).Error; err != nil {
			return err
		}

		if articleUpdate.Title != "" {
			a.Title = articleUpdate.Title
		}
		if articleUpdate.Description != "" {
			a.Description = articleUpdate.Description
		}
		if articleUpdate.Color != "" {
			a.Color = articleUpdate.Color
		}
		if articleUpdate.Size != "" {
			a.Size = articleUpdate.Size
		}

		return tx.Save(&a).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "hero not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update hero: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "Failed update hero", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	articleId := chi.URLParam(r, "articleId")
	a := model.DeleteArticle(articleId)

	utils.RespondWithJSON(w, http.StatusOK, "Delete Article successuful", a)
}