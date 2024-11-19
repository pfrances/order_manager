package http

import (
	"encoding/json"
	"net/http"
	"order_manager/internal/domain"
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

	item, err := s.MenuService.CreateMenuItem(r.Context(), req.Name, req.Price)
	if err != nil {
		if domain.ErrorCode(err) == domain.EINVALID {
			s.logger.Errorf("error creating menu item: %s\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.logger.Errorf("error creating menu item: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONBody(w, http.StatusCreated, item)
}

func (s *Server) HandleGetMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := s.MenuService.FindAllMenuItems(r.Context())
	if err != nil {
		s.logger.Errorf("error finding menu items: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONBody(w, http.StatusOK, items)
}
