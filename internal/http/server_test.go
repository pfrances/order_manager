package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"order_manager/internal/domain"
	domainHttp "order_manager/internal/http"
	"order_manager/internal/sqlite"
	"testing"

	"github.com/stretchr/testify/require"
)

type nopLogger struct{}

func (nopLogger) Infof(format string, args ...interface{})  {}
func (nopLogger) Errorf(format string, args ...interface{}) {}

type repositories struct {
	Table domain.TableRepository
	Menu  domain.MenuRepository
	Bill  domain.BillRepository
}

func MustNewRepositories(t *testing.T) repositories {
	t.Helper()

	logger := nopLogger{}

	db, err := sqlite.NewDB(":memory:", logger)
	if err != nil {
		t.Fatalf("failed to create db: %s", err)
	}
	tableRepo := sqlite.NewTable(db)
	menuRepo := sqlite.NewMenu(db)
	billRepo := sqlite.NewBill(db)

	return repositories{
		Table: tableRepo,
		Menu:  menuRepo,
		Bill:  billRepo,
	}
}

func MustNewServer(t *testing.T, repos repositories) *domainHttp.Server {
	t.Helper()

	logger := nopLogger{}

	tableService := domain.NewTableService(repos.Table)
	menuService := domain.NewMenuService(repos.Menu)
	billService := domain.NewBillService(repos.Bill)

	return domainHttp.NewServer(logger, tableService, menuService, billService)
}

func MustParseReponse[T any](t *testing.T, w *httptest.ResponseRecorder) (body T, statusCode int) {
	t.Helper()

	res := w.Result()
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var v T
	if err := json.Unmarshal(bodyBytes, &v); err != nil {
		t.Fatalf("failed to parse body: %s", err)
	}

	return v, res.StatusCode
}

func MustUnparseBody[T any](t *testing.T, v T) io.Reader {
	t.Helper()

	bodyBytes, err := json.Marshal(v)
	require.NoError(t, err)

	return bytes.NewReader(bodyBytes)
}
