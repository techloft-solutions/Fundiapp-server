package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/direction"
	"github.com/andrwkng/hudumaapp/model"
)

var subscription = app.Subscription{
	SubscriptionID:  "dfabd9de-02cb-4928-9a1c-1f135f7e105a",
	Plan:            "Monthly",
	PlanName:        "Montly subscription",
	Price:           "Ksh 999/Month",
	PaymentMethod:   "Lipa Na M-PESA",
	AutoRenew:       true,
	Status:          "active",
	NextBillingDate: "2022-01-30",
}

var paymentMethods = []app.PaymentMethod{
	{
		ID:     "488a6d76-8b34-4e2f-85af-f611d173470c",
		Name:   "Lipa na M-PESA",
		Method: "mpesa",
		Logo:   "https://mpasho254.files.wordpress.com/2018/11/mpesa.png",
		Number: "071234567",
		Type:   "mobile",
	},
	{
		ID:     "f7260d08-64d5-472d-98f7-fba688410fd1",
		Name:   "Visa",
		Method: "visa",
		Logo:   "https://e7.pngegg.com/pngimages/594/666/png-clipart-visa-logo-credit-card-debit-card-payment-card-bank-visa-blue-text.png",
		Number: "**** **** **** 9439",
		Type:   "card",
	},
}

func (s *Server) handleDeletePaymentMethods(w http.ResponseWriter, r *http.Request) {
	handleSuccessMsg(w, "Payment method deleted successfuly")
}

func (s *Server) handlePaymentMethods(w http.ResponseWriter, r *http.Request) {
	var plans []app.PaymentMethod = paymentMethods
	/*
		plans = append(plans, app.PaymentMethod{
			ID:     "488a6d76-8b34-4e2f-85af-f611d173470c",
			Name:   "Lipa na M-PESA",
			Method: "mpesa",
			Logo:   "https://mpasho254.files.wordpress.com/2018/11/mpesa.png",
			Number: "071234567",
			Type:   "mobile",
		})

		plans = append(plans, app.PaymentMethod{
			ID:     "f7260d08-64d5-472d-98f7-fba688410fd1",
			Name:   "Visa",
			Method: "visa",
			Logo:   "https://e7.pngegg.com/pngimages/594/666/png-clipart-visa-logo-credit-card-debit-card-payment-card-bank-visa-blue-text.png",
			Number: "**** **** **** 9439",
			Type:   "card",
		})
	*/
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(plans); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		return
	}
}

func (s *Server) handleAddMpesaPayment(w http.ResponseWriter, r *http.Request) {
	var paymentMethod = model.PaymentMethod{
		ID:          "4984c20d-3a8c-4679-a613-9c470243d0a7",
		Method:      "mpesa",
		PhoneNumber: "0712345678",
		Type:        "mobile",
	}
	handleSuccessMsgWithRes(w, "Payment method added successfuly", paymentMethod)
}

func (s *Server) handleCancelSubscription(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("success"); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		return
	}
}

func (s *Server) handleMyActiveSubscription(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(subscription); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		return
	}
}

func (s *Server) handleMySubscriptions(w http.ResponseWriter, r *http.Request) {
	var subscriptions []app.Subscription

	subscriptions = make([]app.Subscription, 0)
	subscriptions = append(subscriptions, subscription)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(subscriptions); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		return
	}
}

func (s *Server) handleSubscribePage(w http.ResponseWriter, r *http.Request) {
	var page = app.SubscriptionPage{
		Plans:          plans,
		PaymentMethods: paymentMethods,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(page); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		return
	}

}

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
	//err = mpesa.STKPush()
	//err = mpesa.C2BRegisterURL()
	//err = mpesa.C2BSimulate()
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Payment attempt failed", http.StatusInternalServerError)
		return
	}
	/*
		err = mpesa.GetPayment()
		if err != nil {
			return
		}
	*/
	/*
		err = s.SubSvc.CreateSubscription(r.Context(), &subscription)
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			handleError(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	*/
	subscription.SubscriptionID = "48aa8ba5-07e6-4916-8791-47785c8690a7"
	subscription.ClientID = "9f241cde-23f7-4f59-b942-b78f99b5d28a"
	subscription.StartsAt = time.Now().UTC().Format("2006-01-02 15:04:05")
	subscription.PaymentID = "dd86f1ec-7b69-49d2-9903-88451d5b1903"
	subscription.Status = "active"
	subscription.BillingCycles = 30

	handleSuccessMsgWithRes(w, "Subscription created successfully", subscription)
}

func (s *Server) handleSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptions := make([]model.Subscription, 0)
	handleSuccess(w, subscriptions)
}

func (s *Server) handleSubscriptionCancel(w http.ResponseWriter, r *http.Request) {
	var direction = direction.South
	fmt.Println(direction)
}
