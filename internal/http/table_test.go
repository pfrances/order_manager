package http_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func MustPresaveTables(t *testing.T, repos repositories, tables []domain.Table) {
	t.Helper()

	items := make([]domain.MenuItem, 0)
	for _, table := range tables {
		for _, order := range table.Orders {
			for _, preparation := range order.Preparations {
				items = append(items, preparation.MenuItem)
			}
		}
	}

	for _, item := range items {
		err := repos.Menu.SaveItem(context.Background(), item)
		require.NoError(t, err)
	}

	for _, table := range tables {
		err := repos.Table.Save(context.Background(), table)
		require.NoError(t, err)
	}
}

func TestOpenTableHandler(t *testing.T) {
	repos := MustNewRepositories(t)
	s := MustNewServer(t, repos)

	r := httptest.NewRequest(http.MethodPost, "/table", nil)
	w := httptest.NewRecorder()

	s.HandleOpenTable(w, r)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestGetTablesHandler(t *testing.T) {
	tt := []struct {
		testName   string
		tablesInDB []domain.Table
	}{
		{
			testName:   "empty tables",
			tablesInDB: []domain.Table{},
		},
		{
			testName: "one table with no orders",
			tablesInDB: []domain.Table{
				{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{},
				},
			},
		},
		{
			testName: "one table with orders",
			tablesInDB: []domain.Table{
				{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusDone,
							Preparations: []domain.Preparation{
								{
									ID:       id.New(),
									MenuItem: domain.MenuItem{ID: id.New(), Name: "item", Price: 100},
									Status:   domain.PreparationStatusServed,
								},
							},
						},
					},
				},
			},
		},
		{
			testName: "multiple tables",
			tablesInDB: []domain.Table{
				{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{},
				},
				{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusDone,
							Preparations: []domain.Preparation{
								{
									ID:       id.New(),
									MenuItem: domain.MenuItem{ID: id.New(), Name: "item", Price: 100},
									Status:   domain.PreparationStatusServed,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			repos := MustNewRepositories(t)
			s := MustNewServer(t, repos)

			MustPresaveTables(t, repos, tc.tablesInDB)

			r := httptest.NewRequest(http.MethodGet, "/table", nil)
			w := httptest.NewRecorder()

			s.HandleGetTables(w, r)

			body, statusCode := MustParseReponse[[]domain.Table](t, w)

			require.Equal(t, http.StatusOK, statusCode)
			require.ElementsMatch(t, tc.tablesInDB, body)
		})
	}
}

func TestCloseTableHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tt := []struct {
			testName string
			table    domain.Table
		}{
			{
				testName: "one table with no orders",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{},
				},
			},
			{
				testName: "one table with orders",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusDone,
							Preparations: []domain.Preparation{
								{
									ID:       id.New(),
									MenuItem: domain.MenuItem{ID: id.New(), Name: "item", Price: 100},
									Status:   domain.PreparationStatusServed,
								},
							},
						},
					},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				repos := MustNewRepositories(t)
				s := MustNewServer(t, repos)
				MustPresaveTables(t, repos, []domain.Table{tc.table})

				reqBody := fmt.Sprintf(`{"table_id": "%s"}`, tc.table.ID)
				r := httptest.NewRequest(http.MethodPost, "/table/close", strings.NewReader(reqBody))
				w := httptest.NewRecorder()

				s.HandleCloseTable(w, r)

				res := w.Result()

				require.Equal(t, http.StatusNoContent, res.StatusCode)
			})
		}
	})

	t.Run("Failed", func(t *testing.T) {
		tt := []struct {
			testName           string
			table              domain.Table
			expectedStatusCode int
		}{
			{
				testName: "table already closed",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusClosed,
					Orders: []domain.Order{},
				},
				expectedStatusCode: http.StatusBadRequest,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				repos := MustNewRepositories(t)
				s := MustNewServer(t, repos)
				MustPresaveTables(t, repos, []domain.Table{tc.table})

				reqBody := fmt.Sprintf(`{"table_id": "%s"}`, tc.table.ID)
				r := httptest.NewRequest(http.MethodPost, "/table/close", strings.NewReader(reqBody))
				w := httptest.NewRecorder()

				s.HandleCloseTable(w, r)

				res := w.Result()

				require.Equal(t, tc.expectedStatusCode, res.StatusCode)
			})
		}

		t.Run("table not found", func(t *testing.T) {
			repos := MustNewRepositories(t)
			s := MustNewServer(t, repos)

			reqBody := fmt.Sprintf(`{"table_id": "%s"}`, id.New())
			r := httptest.NewRequest(http.MethodPost, "/table/close", strings.NewReader(reqBody))
			w := httptest.NewRecorder()

			s.HandleCloseTable(w, r)

			res := w.Result()

			require.Equal(t, http.StatusNotFound, res.StatusCode)
		})
	})
}

func TestTakeOrder(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tt := []struct {
			testName  string
			table     domain.Table
			menuItems []domain.MenuItem
		}{
			{
				testName: "one table with no orders",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{},
				},
				menuItems: []domain.MenuItem{
					{ID: id.New(), Name: "item", Price: 100},
				},
			},
			{
				testName: "one table with orders",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusDone,
							Preparations: []domain.Preparation{
								{
									ID:       id.New(),
									MenuItem: domain.MenuItem{ID: id.New(), Name: "item", Price: 100},
									Status:   domain.PreparationStatusServed,
								},
							},
						},
					},
				},
				menuItems: []domain.MenuItem{
					{ID: id.New(), Name: "item", Price: 100},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				repos := MustNewRepositories(t)
				s := MustNewServer(t, repos)
				MustPresaveTables(t, repos, []domain.Table{tc.table})

				err := repos.Menu.SaveItem(context.Background(), tc.menuItems[0])
				require.NoError(t, err)
				reqBody := fmt.Sprintf(`{"table_id": "%s", "menu_item_ids": ["%s"]}`, tc.table.ID, tc.menuItems[0].ID)
				r := httptest.NewRequest(http.MethodPost, "/table/order", strings.NewReader(reqBody))
				w := httptest.NewRecorder()

				s.HandleTakeOrder(w, r)

				res := w.Result()

				require.Equal(t, http.StatusOK, res.StatusCode)
			})
		}
	})

	t.Run("Failed", func(t *testing.T) {

		t.Run("menu item not found", func(t *testing.T) {
			repos := MustNewRepositories(t)
			s := MustNewServer(t, repos)

			table := domain.Table{
				ID:     id.New(),
				Status: domain.TableStatusOpened,
				Orders: []domain.Order{},
			}
			MustPresaveTables(t, repos, []domain.Table{table})

			reqBody := fmt.Sprintf(`{"table_id": "%s", "menu_item_ids": ["%s"]}`, table.ID, id.New())
			r := httptest.NewRequest(http.MethodPost, "/table/order", strings.NewReader(reqBody))
			w := httptest.NewRecorder()

			s.HandleTakeOrder(w, r)

			res := w.Result()

			require.Equal(t, http.StatusNotFound, res.StatusCode)
		})

		t.Run("closed table", func(t *testing.T) {
			repos := MustNewRepositories(t)
			s := MustNewServer(t, repos)

			table := domain.Table{
				ID:     id.New(),
				Status: domain.TableStatusClosed,
				Orders: []domain.Order{
					{
						ID:     id.New(),
						Status: domain.OrderStatusDone,
						Preparations: []domain.Preparation{
							{
								ID:       id.New(),
								MenuItem: domain.MenuItem{ID: id.New(), Name: "item", Price: 100},
								Status:   domain.PreparationStatusServed,
							},
						},
					},
				},
			}
			MustPresaveTables(t, repos, []domain.Table{table})

			item := domain.MenuItem{ID: id.New(), Name: "item", Price: 100}
			err := repos.Menu.SaveItem(context.Background(), item)
			require.NoError(t, err)

			reqBody := fmt.Sprintf(`{"table_id": "%s", "menu_item_ids": ["%s"]}`, table.ID, item.ID)
			r := httptest.NewRequest(http.MethodPost, "/table/order", strings.NewReader(reqBody))
			w := httptest.NewRecorder()

			s.HandleTakeOrder(w, r)

			res := w.Result()

			require.Equal(t, http.StatusBadRequest, res.StatusCode)
		})

		t.Run("table not found", func(t *testing.T) {
			repos := MustNewRepositories(t)
			s := MustNewServer(t, repos)

			item := domain.MenuItem{ID: id.New(), Name: "item", Price: 100}
			err := repos.Menu.SaveItem(context.Background(), item)
			require.NoError(t, err)

			reqBody := fmt.Sprintf(`{"table_id": "%s", "menu_item_ids": ["%s"]}`, id.New(), item.ID)
			r := httptest.NewRequest(http.MethodPost, "/table/order", strings.NewReader(reqBody))
			w := httptest.NewRecorder()

			s.HandleTakeOrder(w, r)

			res := w.Result()

			require.Equal(t, http.StatusNotFound, res.StatusCode)
		})

		t.Run("menu item not found", func(t *testing.T) {
			repos := MustNewRepositories(t)
			s := MustNewServer(t, repos)

			table := domain.Table{
				ID:     id.New(),
				Status: domain.TableStatusOpened,
				Orders: []domain.Order{},
			}
			MustPresaveTables(t, repos, []domain.Table{table})

			reqBody := fmt.Sprintf(`{"table_id": "%s", "menu_item_ids": ["%s"]}`, id.New(), id.New())
			r := httptest.NewRequest(http.MethodPost, "/table/order", strings.NewReader(reqBody))
			w := httptest.NewRecorder()

			s.HandleTakeOrder(w, r)

			res := w.Result()

			require.Equal(t, http.StatusNotFound, res.StatusCode)
		})
	})
}

func TestStartPreparation(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tt := []struct {
			testName             string
			tableFromPreparation func(prep domain.Preparation) domain.Table
			preparation          domain.Preparation
		}{
			{
				testName: "one table with one order and one preparation",
				tableFromPreparation: func(prep domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{
								ID:           id.New(),
								Status:       domain.OrderStatusTaken,
								Preparations: []domain.Preparation{prep},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					MenuItem: domain.MenuItem{ID: id.New(), Name: "item", Price: 100},
					Status:   domain.PreparationStatusPending,
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				repos := MustNewRepositories(t)
				s := MustNewServer(t, repos)
				table := tc.tableFromPreparation(tc.preparation)
				MustPresaveTables(t, repos, []domain.Table{table})

				reqBody := fmt.Sprintf(`{"preparation_id": "%s"}`, tc.preparation.ID)
				r := httptest.NewRequest(http.MethodPost, "/table/preparation", strings.NewReader(reqBody))
				w := httptest.NewRecorder()

				s.HandleStartPreparation(w, r)

				res := w.Result()

				require.Equal(t, http.StatusNoContent, res.StatusCode)
			})
		}
	})

	t.Run("Failed", func(t *testing.T) {
		tt := []struct {
			testName             string
			tableFromPreparation func(prep domain.Preparation) domain.Table
			preparation          domain.Preparation
			expectedStatusCode   int
		}{
			{
				testName: "preparation not found",
				tableFromPreparation: func(prep domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: make([]domain.Order, 0),
					}
				},
				preparation:        domain.Preparation{ID: id.NilID()},
				expectedStatusCode: http.StatusNotFound,
			},
			{
				testName: "preparation already started",
				tableFromPreparation: func(prep domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{
								ID:           id.New(),
								Status:       domain.OrderStatusTaken,
								Preparations: []domain.Preparation{prep},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					MenuItem: domain.MenuItem{ID: id.New(), Name: "item", Price: 100},
					Status:   domain.PreparationStatusInProgress,
				},
				expectedStatusCode: http.StatusBadRequest,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				repos := MustNewRepositories(t)
				s := MustNewServer(t, repos)
				table := tc.tableFromPreparation(tc.preparation)
				MustPresaveTables(t, repos, []domain.Table{table})

				reqBody := fmt.Sprintf(`{"preparation_id": "%s"}`, tc.preparation.ID)
				r := httptest.NewRequest(http.MethodPost, "/table/preparation", strings.NewReader(reqBody))
				w := httptest.NewRecorder()

				s.HandleStartPreparation(w, r)

				res := w.Result()

				require.Equal(t, tc.expectedStatusCode, res.StatusCode)
			})
		}
	})
}
