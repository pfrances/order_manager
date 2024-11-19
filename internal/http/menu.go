package http

import (
	"encoding/json"
	"net/http"
)

func (s *Server) registerMenuRoutes(r *router) {
	menuRouter := r.group("/menu")

	menuRouter.HandleFunc("POST /item", s.HandleAddMenuItem)
	menuRouter.HandleFunc("GET /item", s.HandleGetMenuItems)
}

func (s *Server) HandleAddMenuItem(w http.ResponseWriter, r *http.Request) {
	type addMenuItemRequest struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}

	var req addMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Errorf("error decoding request: %w\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		s.logger.Errorf("empty name\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Price < 0 {
		s.logger.Errorf("negative price\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, err := s.MenuService.CreateMenuItem(r.Context(), req.Name, req.Price)
	if err != nil {
		s.logger.Errorf("error creating menu item: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	writeJSONBody(w, http.StatusCreated, item)
}

func (s *Server) HandleGetMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := s.MenuService.FindAllMenuItems(r.Context())
	if err != nil {
		s.logger.Errorf("error finding menu items: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	writeJSONBody(w, http.StatusOK, items)
}
