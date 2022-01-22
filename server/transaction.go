package server

import (
	"fmt"
	"net/http"
)

func (s *Server) handleTransactionConfirm(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	fmt.Println("Confirm transaction")

	for k, v := range r.Form {
		fmt.Printf("%s = %s\n", k, v)
	}

}
