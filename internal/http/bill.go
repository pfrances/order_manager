package http

import (
	"encoding/json"
	"net/http"
	"order_manager/internal/id"
)

func (s *server) registerBillRoutes(r *router) {
	billRouter := r.group("/bill")

	billRouter.HandleFunc("POST /", s.handleGenerateBill)
}

func (s *server) handleGenerateBill(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		TableID id.ID `json:"table_id"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Errorf("error decoding request: %s\n", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	table, err := s.tableService.FindTable(r.Context(), req.TableID)
	if err != nil {
		s.logger.Errorf("error finding table: %s\n", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	bill, err := s.billService.GenerateBill(r.Context(), table)
	if err != nil {
		s.logger.Errorf("error generating bill: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONBody(w, http.StatusOK, bill)
}
