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
	BkSvc   app.BookingService
	CatSvc  app.CategoryService
	PfoSvc  app.PortfolioService
	LocSvc  app.LocationService
	BidSvc  app.BidService
	ReqSvc  app.RequestService
	UsrSvc  app.UserService
	RevSvc  app.ReviewService
	IndSvc  app.IndustryService
	SrchSvc app.SearchService
}

func New() *Server {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),
	}
	s.router.HandleFunc("/", handleHome).Methods("GET")
	// Users
	s.router.HandleFunc("/user", s.handleUserCreate).Methods("POST")
	s.router.HandleFunc("/user", s.handleUserGet).Methods("GET")
	s.router.HandleFunc("/user/validate", s.handleUserValidate).Methods("POST")

	// Tesing
	s.testingRoutes(s.router)

	r := s.router.PathPrefix("/").Subrouter()
	r.Use(middlewares.AuthHandler)
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

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HudumaApp API V1\n")
}

func (s *Server) testingRoutes(r *mux.Router) {

}

func (s *Server) registerRoutes(r *mux.Router) {
	r.HandleFunc("/user/password", s.handlePasswordNew).Methods("POST")
	r.HandleFunc("/user/password", s.handlePasswordChange).Methods("PUT")
	r.HandleFunc("/user/{id}", s.handleUserByID).Methods("GET")
	// Profile
	//r.HandleFunc("/profile", s.handleProfileCreate).Methods("POST")
	r.HandleFunc("/profile", s.handleProfileGet).Methods("GET")
	r.HandleFunc("/profile", s.handleProfileUpdate).Methods("PUT")
	r.HandleFunc("/profile/location", s.handleProfileLocationUpdate).Methods("PUT")
	//r.HandleFunc("/profile/{id}", s.handleProfileDelete).Methods("DELETE")
	// Providers
	r.HandleFunc("/provider", s.handleProviderGet).Methods("GET")
	r.HandleFunc("/providers", s.handleProviderList).Methods("GET")
	r.HandleFunc("/providers/{id}", s.handleProviderByID).Methods("GET")
	r.HandleFunc("/providers", s.handleProviderUpdate).Methods("PUT")
	//r.HandleFunc("/providers/{id}/reviews", s.handleProviderReviews).Methods("GET")
	r.HandleFunc("/providers/{id}/services", s.handleProviderServices).Methods("GET")
	r.HandleFunc("/providers/{id}/portfolios", s.handleProviderPortfolios).Methods("GET")
	r.HandleFunc("/providers/{id}/bookings", s.handleProviderBookings).Methods("GET")
	r.HandleFunc("/providers/{id}/bookings/{id}", s.handleProviderBooking).Methods("GET")
	// Locations
	r.HandleFunc("/locations", s.handleMyLocations).Methods("GET")
	r.HandleFunc("/locations", s.handleLocationCreate).Methods("POST")
	r.HandleFunc("/locations/{id}", s.handleLocationDelete).Methods("DELETE")
	// Categories
	r.HandleFunc("/categories", s.handleCategoriesList).Methods("GET")
	r.HandleFunc("/categories", s.handleCategoryCreate).Methods("POST")
	r.HandleFunc("/categories/root", s.handleCategoriesRoot).Methods("GET")
	// Industries
	r.HandleFunc("/industries", s.handleIndustriesList).Methods("GET")
	r.HandleFunc("/industries", s.handleIndustryCreate).Methods("POST")
	// Reviews
	r.HandleFunc("/reviews", s.handleReviewCreate).Methods("POST")
	//r.HandleFunc("/reviews", s.handleReviewList).Methods("GET")
	// Services
	r.HandleFunc("/services", s.handleMyServices).Methods("GET")
	r.HandleFunc("/services", s.handleServiceCreate).Methods("POST")
	// Request
	r.HandleFunc("/requests", s.handleRequestList).Methods("GET")
	r.HandleFunc("/requests", s.handleRequestCreate).Methods("POST")
	r.HandleFunc("/requests/{id}", s.handleRequest).Methods("GET")
	//r.HandleFunc("/requests/{id}/cancel", s.handleRequestCancel).Methods("PUT")
	r.HandleFunc("/requests/{id}/bids", s.handleRequestBids).Methods("GET")
	// All requests
	r.HandleFunc("/all-requests", s.handleAllRequests).Methods("GET")
	r.HandleFunc("/all-requests/recommended", s.handleRecommendedRequests).Methods("GET")
	r.HandleFunc("/all-requests/categories", s.handleRequestCategories).Methods("GET")
	r.HandleFunc("/all-requests/instant-search", s.handleRequestInstantSearch).Methods("GET")
	r.HandleFunc("/all-requests/search", s.handleRequestSearch).Methods("GET")
	//r.HandleFunc("/all-requests/search", s.handleRequestSearch).Methods("GET")
	// Bookings
	r.HandleFunc("/bookings/{id}", s.handleBookingByID).Methods("GET")
	r.HandleFunc("/bookings", s.handleBookingList).Methods("GET")
	r.HandleFunc("/bookings", s.handleBookingCreate).Methods("POST")
	//r.HandleFunc("/bookings/{id}", s.handleBookingUpdate).Methods("PUT")
	//r.HandleFunc("/bookings/{id}", s.handleBookingDelete).Methods("DELETE")
	// Bids
	r.HandleFunc("/bids", s.handleBidCreate).Methods("POST")
	r.HandleFunc("/bids", s.handleMyBids).Methods("GET")
	r.HandleFunc("/bids/{id}/accept", s.handleAcceptBid).Methods("PUT")
	//r.HandleFunc("/bids/{id}/cancel", s.handleCancelBid).Methods("DELETE")
	// Portfolios
	r.HandleFunc("/portfolios", s.handleMyPortfolio).Methods("GET")
	r.HandleFunc("/portfolios", s.handlePortfolioCreate).Methods("POST")
	r.HandleFunc("/portfolios/{id}", s.handlePortfolio).Methods("GET")
	// Search
	r.HandleFunc("/search", s.handleSearch).Methods("GET")
	// Payments
	//r.HandleFunc("/payments", s.handlePaymentCreate).Methods("POST")
	//r.HandleFunc("/payments", s.handlePaymentList).Methods("GET")
	// Preferences
	//r.HandleFunc("/preferences", s.handlePreferenceList).Methods("GET")
	//r.HandleFunc("/preferences", s.handlePreferenceCreate).Methods("POST")
	r.HandleFunc("/test", s.handleTest).Methods("GET")
}
