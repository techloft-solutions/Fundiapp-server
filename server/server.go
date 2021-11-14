package server

import (
	"fmt"
	"net/http"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/gorilla/mux"
)

type Server struct {
	Addr   string
	server *http.Server
	router *mux.Router
	//ln     net.Listener
	BkSvc      app.BookingService
	CatSvc     app.CategoryService
	PfoSvc     app.PortfolioService
	LocSvc     app.LocationService
	BidSvc     app.BidService
	ReqSvc     app.RequestService
	UsrSvc     app.UserService
	RevSvc     app.ReviewService
	ServiceSvc app.ServiceService
}

func New() *Server {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),
	}
	s.router.Use(middlewares.AuthHandler)
	s.router.HandleFunc("/", handleHello).Methods("GET")
	r := s.router.PathPrefix("/").Subrouter()
	s.registerRoutes(r)
	return s
}

func (s *Server) Start() (err error) {
	//if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
	//	return err
	//}
	return http.ListenAndServe(s.Addr, s.router)
	//return nil
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HudumaApp API V1\n")
}

func (s *Server) registerRoutes(r *mux.Router) {
	// Profile
	r.HandleFunc("/profile", s.handleProfileCreate).Methods("POST")
	r.HandleFunc("/profile", s.handleProfileGet).Methods("GET")
	r.HandleFunc("/profile", s.handleProfileUpdate).Methods("PUT")
	//r.HandleFunc("/profile/{id}", s.handleProfileDelete).Methods("DELETE")
	// Providers
	r.HandleFunc("/providers", s.handleProviderList).Methods("GET")
	r.HandleFunc("/providers/{id}", s.handleProviderByID).Methods("GET")
	r.HandleFunc("/providers", s.handleProviderCreate).Methods("POST")
	// Locations
	r.HandleFunc("/locations", s.handleMyLocations).Methods("GET")
	r.HandleFunc("/locations", s.handleLocationCreate).Methods("POST")
	// Categories
	r.HandleFunc("/categories", s.handleCategoriesList).Methods("GET")
	r.HandleFunc("/categories", s.handleCategoryCreate).Methods("POST")
	// Reviews
	r.HandleFunc("/reviews", s.handleReviewCreate).Methods("POST")
	//r.HandleFunc("/reviews", s.handleReviewList).Methods("GET")
	// Services
	r.HandleFunc("/services", s.handleMyServices).Methods("GET")
	r.HandleFunc("/services", s.handleServiceCreate).Methods("POST")
	// Request
	//r.HandleFunc("/requests", s.handleRequestList).Methods("GET")
	//r.HandleFunc("/requests", s.handleRequestCreate).Methods("POST")
	//r.HandleFunc("/requests/{id}", s.handleRequest).Methods("GET")
	// Bookings
	//r.HandleFunc("/bookings/{id}", s.handleBookingByID).Methods("GET")
	//r.HandleFunc("/bookings", s.handleBookingList).Methods("GET")
	//r.HandleFunc("/bookings", s.handleBookingCreate).Methods("POST")
	//r.HandleFunc("/bookings/{id}", s.handleBookingUpdate).Methods("PUT")
	//r.HandleFunc("/bookings/{id}", s.handleBookingDelete).Methods("DELETE")
	// Bids
	//r.HandleFunc("/bids", s.handleBidCreate).Methods("POST")
	//r.HandleFunc("/bids", s.handleBookingList).Methods("GET")
	// Portfolios
	//r.HandleFunc("/portfolios", s.handlePortfolioList).Methods("GET")
	// Payments
	//r.HandleFunc("/payments", s.handlePaymentCreate).Methods("POST")
	//r.HandleFunc("/payments", s.handlePaymentList).Methods("GET")
	// Preferences
	//r.HandleFunc("/preferences", s.handlePreferenceList).Methods("GET")
	//r.HandleFunc("/preferences", s.handlePreferenceCreate).Methods("POST")
}
