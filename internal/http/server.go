package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type tableService interface {
	FindTable(ctx context.Context, tableID id.ID) (domain.Table, error)
	FindOpenedTables(ctx context.Context) ([]domain.Table, error)
	OpenTable(ctx context.Context) (domain.Table, error)
	CloseTable(ctx context.Context, tableID id.ID) error
	FinishPreparation(ctx context.Context, preparationID id.ID) error
	ServePreparation(ctx context.Context, preparationID id.ID) error
	StartPreparation(ctx context.Context, preparationID id.ID) error
	TakeOrder(ctx context.Context, tableID id.ID, menuItems []domain.MenuItem) (domain.Order, error)
}

type menuService interface {
	FindAllMenuItems(ctx context.Context) ([]domain.MenuItem, error)
	FindMenuItems(ctx context.Context, ids []id.ID) ([]domain.MenuItem, error)
	AddItemToCategory(ctx context.Context, categoryID id.ID, itemID id.ID) error
	CreateCategory(ctx context.Context, name string) (domain.MenuCategory, error)
	CreateMenuItem(ctx context.Context, name string, price int) (domain.MenuItem, error)
}

type billService interface {
	GenerateBill(ctx context.Context, table domain.Table) (domain.Bill, error)
	PayBill(ctx context.Context, billID id.ID, amount int) error
}

type middleware func(http.Handler) http.Handler

type router struct {
	*http.ServeMux
	prefix      string
	middlewares []middleware
}

func newRouter() *router {
	return &router{ServeMux: http.NewServeMux()}
}

func (r *router) group(prefix string, groupMiddleware ...middleware) *router {
	return &router{
		ServeMux:    r.ServeMux,
		prefix:      r.prefix + prefix,
		middlewares: append(r.middlewares, groupMiddleware...),
	}
}

func (r *router) HandleFunc(pattern string, handler http.HandlerFunc) {
	parts := strings.SplitN(pattern, " ", 2)
	if len(parts) != 2 {
		log.Fatalf("invalid pattern: %s\n", pattern)
	}

	method := parts[0]
	path := r.prefix + parts[1]

	fullPattern := method + " " + path

	finalHandler := http.Handler(handler)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = r.middlewares[i](finalHandler)
	}

	r.ServeMux.Handle(fullPattern, finalHandler)
}

type Server struct {
	server *http.Server
	router *router

	logger logger

	TableService tableService
	MenuService  menuService
	BillService  billService

	URL string
}

func NewServer(logger logger, tableService tableService, menuService menuService, billService billService) *Server {
	s := &Server{
		logger:       logger,
		TableService: tableService,
		MenuService:  menuService,
		BillService:  billService,
	}
	router := newRouter().group("/api", s.logMiddleware)
	s.registerTableRoutes(router)
	s.registerMenuRoutes(router)
	s.registerBillRoutes(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	s.server = server
	s.router = router

	s.URL = fmt.Sprintf("http://localhost%s", server.Addr)

	return s
}

func (s *Server) Run(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		s.logger.Infof("listening on %s\n", s.server.Addr)

		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		return s.server.Shutdown(gCtx)
	})

	return g.Wait()
}

func writeJSONBody(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
	}
}

func domainErrorToHTTPStatus(err error) int {
	switch domain.ErrorCode(err) {
	case domain.ENOTFOUND:
		return http.StatusNotFound
	case domain.EINVALID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
}

type logResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &logResponseWriter{w, http.StatusInternalServerError}

		next.ServeHTTP(lrw, r)

		s.logger.Infof("[http] %s %s - %d [%v]\n", r.Method, r.URL.Path, lrw.status, time.Since(start))
	})
}
