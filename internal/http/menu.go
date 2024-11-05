package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *server) registerMenuRoutes(r *router) {
	menuRouter := r.group("/menu")

	menuRouter.HandleFunc("POST /item", s.handleAddMenuItem)
}

func (s *server) handleAddMenuItem(w http.ResponseWriter, r *http.Request) {
	type addMenuItemRequest struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}

	var req addMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Errorf("error decoding request: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, err := s.menuService.CreateMenuItem(r.Context(), req.Name, req.Price)
	if err != nil {
		s.logger.Errorf("error creating menu item: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%v", item)))
}
