package server

import (
	"log"
	"net/http"
)

func (s *Server) handlePlans(w http.ResponseWriter, r *http.Request) {
	plans, err := s.PlanSvc.GetAllPlans(r.Context())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, plans)
}
