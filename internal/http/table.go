package http

import (
	"encoding/json"
	"net/http"
	"order_manager/internal/id"
)

func (s *Server) registerTableRoutes(r *router) {
	tableRouter := r.group("/table")

	tableRouter.HandleFunc("GET /", s.HandleGetTables)
	tableRouter.HandleFunc("POST /", s.HandleOpenTable)
	tableRouter.HandleFunc("POST /order", s.HandleTakeOrder)
	tableRouter.HandleFunc("POST /close", s.HandleCloseTable)
}

func (s *Server) HandleGetTables(w http.ResponseWriter, r *http.Request) {
	tables, err := s.TableService.FindOpenedTables(r.Context())
	if err != nil {
		s.logger.Errorf("error finding tables: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	writeJSONBody(w, http.StatusOK, tables)
}

func (s *Server) HandleOpenTable(w http.ResponseWriter, r *http.Request) {
	table, err := s.TableService.OpenTable(r.Context())
	if err != nil {
		s.logger.Errorf("error creating table: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	writeJSONBody(w, http.StatusCreated, table)
}

func (s *Server) HandleCloseTable(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		TableID id.ID `json:"table_id"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Errorf("error decoding request: %s\n", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.TableService.CloseTable(r.Context(), req.TableID); err != nil {
		s.logger.Errorf("error closing table: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) HandleTakeOrder(w http.ResponseWriter, r *http.Request) {
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

	menuItems, err := s.MenuService.FindMenuItems(r.Context(), req.MenuItemIds)
	if err != nil {
		s.logger.Errorf("error finding menu items: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	order, err := s.TableService.TakeOrder(r.Context(), req.TableID, menuItems)
	if err != nil {
		s.logger.Errorf("error taking order: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	writeJSONBody(w, http.StatusOK, order)
}

func (s *Server) HandleStartPreparation(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		PreparationID id.ID `json:"preparation_id"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Errorf("error decoding request: %s\n", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err := s.TableService.StartPreparation(r.Context(), req.PreparationID)
	if err != nil {
		s.logger.Errorf("error starting preparation: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) FinishPreparation(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		OrderID id.ID `json:"order_id"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Errorf("error decoding request: %s\n", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err := s.TableService.FinishPreparation(r.Context(), req.OrderID)
	if err != nil {
		s.logger.Errorf("error finishing preparation: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) HandleServePreparation(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		PreparationID id.ID `json:"preparation_id"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Errorf("error decoding request: %s\n", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err := s.TableService.ServePreparation(r.Context(), req.PreparationID)
	if err != nil {
		s.logger.Errorf("error serving preparation: %s\n", err)
		writeError(w, domainErrorToHTTPStatus(err), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
