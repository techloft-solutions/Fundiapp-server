package server

import (
	"encoding/json"
	"log"
	"net/http"

	app "github.com/andrwkng/hudumaapp"
)

var plan = app.Plan{
	ID:           "98975604-428c-4335-9474-b99110b2d019",
	Name:         "Weekly",
	Description:  "24km distance",
	Price:        199,
	Currency:     "Ksh",
	Interval:     1,
	IntervalUnit: "week",
}

var plans = []app.Plan{
	plan,
	{
		ID:           "098974ab-441c-4aa9-bf9c-40c6a95bdb5d",
		Name:         "Monthly",
		Description:  "24km distance + distance",
		Price:        999,
		Currency:     "Ksh",
		Interval:     1,
		IntervalUnit: "month",
	},
	{
		ID:           "3363ace9-4c7e-482b-8f3b-ae66f7f88f62",
		Name:         "Yearly",
		Description:  "24km distance + distance + premium",
		Price:        999,
		Currency:     "Ksh",
		Interval:     1,
		IntervalUnit: "year",
	},
}

func (s *Server) handlePlans(w http.ResponseWriter, r *http.Request) {
	/*
		plans, err := s.PlanSvc.GetAllPlans(r.Context())
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			handleError(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		handleSuccess(w, plans)
	*/
	/*
		plans = append(plans, app.Plan{
			ID:           "98975604-428c-4335-9474-b99110b2d019",
			Name:         "Weekly",
			Description:  "24km distance",
			Price:        199,
			Currency:     "Ksh",
			Interval:     1,
			IntervalUnit: "week",
		})

		plans = append(plans, app.Plan{
			ID:           "098974ab-441c-4aa9-bf9c-40c6a95bdb5d",
			Name:         "Monthly",
			Description:  "24km distance + distance",
			Price:        999,
			Currency:     "Ksh",
			Interval:     1,
			IntervalUnit: "month",
		})

		plans = append(plans, app.Plan{
			ID:           "3363ace9-4c7e-482b-8f3b-ae66f7f88f62",
			Name:         "Yearly",
			Description:  "24km distance + distance + premium",
			Price:        999,
			Currency:     "Ksh",
			Interval:     1,
			IntervalUnit: "year",
		})
	*/
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(plans); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		return
	}
}
