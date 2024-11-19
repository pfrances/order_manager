package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMenuItem(t *testing.T) {
	repos := MustNewRepositories(t)
	s := MustNewServer(t, repos)

	tt := []struct {
		testName string
		name     string
		price    int
		status   int
	}{
		{testName: "valid item", name: "item1", price: 100, status: http.StatusCreated},
		{testName: "empty name", name: "", price: 100, status: http.StatusBadRequest},
		{testName: "negative price", name: "item2", price: -1, status: http.StatusBadRequest},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			body := fmt.Sprintf(`{"name":"%s","price":%d}`, tc.name, tc.price)
			r := httptest.NewRequest(http.MethodPost, "/menu/item", strings.NewReader(body))
			w := httptest.NewRecorder()

			s.HandleAddMenuItem(w, r)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tc.status, res.StatusCode)

			if tc.status != http.StatusCreated {
				return
			}

			var item domain.MenuItem
			if err := json.NewDecoder(res.Body).Decode(&item); err != nil {
				t.Fatalf("failed to decode response: %s", err)
			}

			assert.Equal(t, tc.name, item.Name)
			assert.Equal(t, tc.price, item.Price)
			assert.NotEqual(t, id.NilID(), item.ID)
		})
	}
}

func TestGetMenuItems(t *testing.T) {
	repos := MustNewRepositories(t)
	s := MustNewServer(t, repos)

	t.Run("no items", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/menu/item", nil)
		w := httptest.NewRecorder()

		s.HandleGetMenuItems(w, r)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode)

		var items []domain.MenuItem
		if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
			t.Fatalf("failed to decode response: %s", err)
		}

		assert.Empty(t, items)
	})

	t.Run("with items", func(t *testing.T) {
		ctx := context.Background()
		item1, err := s.MenuService.CreateMenuItem(ctx, "item1", 100)
		require.NoError(t, err)

		item2, err := s.MenuService.CreateMenuItem(ctx, "item2", 200)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/menu/item", nil)
		w := httptest.NewRecorder()

		s.HandleGetMenuItems(w, r)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode)

		var items []domain.MenuItem
		if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
			t.Fatalf("failed to decode response: %s", err)
		}

		assert.ElementsMatch(t, []domain.MenuItem{item1, item2}, items)
	})
}
