package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/andrwkng/hudumaapp/direction"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/payments/mpesa"
)

func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	var subscription model.Subscription

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleMysqlErrors(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	if err := json.Unmarshal(jsonStr, &subscription); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	// Make payment
	err = mpesa.MakePayment()
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Payment attempt failed", http.StatusInternalServerError)
		return
	}

	err = mpesa.GetPayment()
	if err != nil {
		return
	}

	/*err = s.SubSvc.CreateSubscription(r.Context(), &subscription)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Subscription created successfully", subscription)
	*/
}

func (s *Server) handleSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptions := make([]model.Subscription, 0)
	handleSuccess(w, subscriptions)
}

func (s *Server) handleSubscriptionCancel(w http.ResponseWriter, r *http.Request) {
	var direction = direction.South
	fmt.Println(direction)
}
