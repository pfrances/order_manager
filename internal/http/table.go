package http

import (
	"encoding/json"
	"net/http"
	"order_manager/internal/id"
)

func (s *server) registerTableRoutes(r *router) {
	tableRouter := r.group("/table")

	tableRouter.HandleFunc("GET /", s.handleGetTables)
	tableRouter.HandleFunc("POST /", s.handleOpenTable)
	tableRouter.HandleFunc("POST /order", s.handleTakeOrder)
	tableRouter.HandleFunc("POST /close", s.handleCloseTable)
}

func (s *server) handleGetTables(w http.ResponseWriter, r *http.Request) {
	tables, err := s.tableService.FindOpenedTables(r.Context())
	if err != nil {
		s.logger.Errorf("error finding tables: %s\n", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSONBody(w, http.StatusOK, tables)
}

func (s *server) handleOpenTable(w http.ResponseWriter, r *http.Request) {
	table, err := s.tableService.OpenTable(r.Context())
	if err != nil {
		s.logger.Errorf("error creating table: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONBody(w, http.StatusCreated, table)
}

func (s *server) handleCloseTable(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		TableID id.ID `json:"table_id"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Errorf("error decoding request: %s\n", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.tableService.CloseTable(r.Context(), req.TableID); err != nil {
		s.logger.Errorf("error closing table: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handleTakeOrder(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		TableID     id.ID   `json:"table_id"`
		MenuItemIds []id.ID `json:"menu_item_ids"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Errorf("error decoding request: %s\n", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	menuItems, err := s.menuService.FindMenuItems(r.Context(), req.MenuItemIds)
	if err != nil {
		s.logger.Errorf("error finding menu items: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	order, err := s.tableService.TakeOrder(r.Context(), req.TableID, menuItems)
	if err != nil {
		s.logger.Errorf("error taking order: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONBody(w, http.StatusOK, order)
}
