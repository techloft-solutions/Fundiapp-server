package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/google/uuid"
)

func (s *Server) handleCategoriesList(w http.ResponseWriter, r *http.Request) {
	var err error
	var categories []*app.Category
	parent_id := r.URL.Query().Get("parent_id")
	// Fetch categories from database.
	switch {
	case parent_id != "":
		categories, err = s.CatSvc.ListCategoriesByParentID(r.Context(), parent_id)
	default:
		categories, err = s.CatSvc.ListCategories(r.Context())
	}
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, categories)
}

func (s *Server) handleCategoryCreate(w http.ResponseWriter, r *http.Request) {
	var category model.Category

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleMysqlErrors(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	if err := json.Unmarshal(jsonStr, &category); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	if err := category.Validate(); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.CatSvc.CreateCategory(r.Context(), &category)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Category created successfully", category)
}

func (s *Server) handleIndustryCreate(w http.ResponseWriter, r *http.Request) {
	var industry model.Industry

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleMysqlErrors(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	if err := json.Unmarshal(jsonStr, &industry); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	if err := industry.Validate(); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.IndSvc.CreateIndustry(r.Context(), &industry)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Industry created successfully", industry)
}

func (s *Server) handleIndustriesList(w http.ResponseWriter, r *http.Request) {
	// Fetch industries from database.
	industries, err := s.IndSvc.ListIndustries(r.Context())

	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, industries)
}

func (s *Server) handleReviewCreate(w http.ResponseWriter, r *http.Request) {
	var review model.Review

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &review); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	if err := review.Validate(); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.RevSvc.CreateReview(r.Context(), &review)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Review created successfully", review)
}

/*
func (s *Server) handleReviewCreate(w http.ResponseWriter, r *http.Request) {
	//var review model.Review
	review := model.Review{}

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &review); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	if err := review.Validate(); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	handleSuccessMsgWithRes(w, "Review created successfully", review)

	err = s.RevSvc.CreateReview(r.Context(), &review)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

}
*/
func (s *Server) handleServiceCreate(w http.ResponseWriter, r *http.Request) {
	var service model.Service

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	// only providers can create services
	provider, err := s.UsrSvc.FindProviderByUserID(r.Context(), userID.String())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleUnathorised(w)
			return
		}
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	providerID, err := uuid.Parse(provider.ID)
	if err != nil {
		log.Println("error parsing providerID:", err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	service.ProviderID = providerID

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &service); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	if err := service.Validate(); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.UsrSvc.CreateService(r.Context(), &service)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Service created successfully", service)
}

func (s *Server) handleMyServices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	resp, err := s.UsrSvc.ListMyServices(ctx, userID.String())
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, resp)
}
