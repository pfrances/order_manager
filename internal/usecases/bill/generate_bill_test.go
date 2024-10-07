package bill_test

import (
	"order_manager/internal/model"
	"order_manager/internal/usecases/bill"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBillGenerationSuccess(t *testing.T) {
	table := model.Table{}
	usecase := bill.NewGenerateBill()

	_, err := usecase.Execute(table.ID)
	require.NoError(t, err)
}
