package controllers

import (
	"net/http"

	"github.com/nitomibimo/CRUD-GO/api/response"
)

// Home - Index api
func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}
