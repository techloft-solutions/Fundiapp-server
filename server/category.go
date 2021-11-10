package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
)

func (s *Server) handleCategoriesList(w http.ResponseWriter, r *http.Request) {
	// Fetch categories from database.
	resp, err := s.CatSvc.ListCategories(r.Context())
	if err != nil {
		log.Println(err)
		//handleError(w, errors.New("Something went wrong!"), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (s *Server) handleCategoryCreate(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	category.Name = r.PostFormValue("name")
	category.Profession = r.PostFormValue("profession")
	category.Description = r.PostFormValue("description")
	parentId, err := strconv.Atoi(r.PostFormValue("parent_id"))
	if err != nil {
		log.Println(err)
		//handleError(w, errors.New("Something went wrong!"), http.StatusInternalServerError)
		return
	}
	category.ParentID = parentId
	err = s.CatSvc.CreateCategory(r.Context(), &category)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleUnathorised(w)
		return
	}
	res := app.Category{
		Name: category.Name,
	}
	handleSuccess(w, res)
}
