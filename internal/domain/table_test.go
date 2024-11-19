package domain_test

import (
	"context"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"order_manager/internal/inmem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsTableValid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			table    domain.Table
		}{
			{
				testName: "opened table with no orders",
				table:    domain.Table{ID: id.New(), Status: domain.TableStatusOpened, Orders: make([]domain.Order, 0)},
			},
			{
				testName: "opened table with orders",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusTaken,
							Preparations: []domain.Preparation{{
								ID:       id.New(),
								Status:   domain.PreparationStatusPending,
								MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
							}},
						},
					},
				},
			},
			{
				testName: " closed table",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusClosed,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusDone,
							Preparations: []domain.Preparation{{
								ID:       id.New(),
								Status:   domain.PreparationStatusServed,
								MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
							}},
						},
					},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				assert.Truef(t, tc.table.IsValid(), "should be valid table: %v", tc.table)
			})
		}
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			table    domain.Table
		}{
			{
				testName: "Invalid table ID",
				table:    domain.Table{ID: id.NilID(), Status: domain.TableStatusOpened, Orders: make([]domain.Order, 0)},
			},
			{
				testName: "orders slice nil",
				table:    domain.Table{ID: id.New(), Status: domain.TableStatusOpened, Orders: nil},
			},
			{
				testName: "order's preparation slice empty",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{ID: id.NilID(), Status: domain.OrderStatusTaken, Preparations: make([]domain.Preparation, 0)},
					},
				},
			},
			{
				testName: "invalid closed table order status",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusClosed,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusTaken,
							Preparations: []domain.Preparation{{
								ID:       id.New(),
								Status:   domain.PreparationStatusPending,
								MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
							}},
						},
					},
				},
			},
			{
				testName: "invalid preparation status of done order",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusDone,
							Preparations: []domain.Preparation{{
								ID:       id.New(),
								Status:   domain.PreparationStatusPending,
								MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
							}},
						},
					},
				},
			},
			{
				testName: "invalid preparation status of aborted order",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusAborted,
							Preparations: []domain.Preparation{{
								ID:       id.New(),
								Status:   domain.PreparationStatusPending,
								MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
							}},
						},
					},
				},
			},
			{
				testName: "table order preparations",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{ID: id.New(), Status: domain.OrderStatusTaken, Preparations: nil},
					},
				},
			},
			{
				testName: "table order preparation",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:           id.New(),
							Status:       domain.OrderStatusTaken,
							Preparations: []domain.Preparation{{ID: id.NilID(), Status: domain.PreparationStatusPending}},
						},
					},
				},
			},
			{
				testName: "table order preparation status",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:           id.New(),
							Status:       domain.OrderStatusTaken,
							Preparations: []domain.Preparation{{ID: id.New(), Status: domain.PreparationStatusServed}},
						},
					},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				assert.Falsef(t, tc.table.IsValid(), "should be invalid table: %v", tc.table)
			})
		}
	})
}

func TestCreateTable(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		table, err := tableService.OpenTable(context.Background())
		require.NoError(t, err, "table creation failed")

		assert.NotEqual(t, table.ID, id.NilID(), "generated table ID is nil")
		assert.Equal(t, domain.TableStatusOpened, table.Status, "invalid table status")
		assert.NotNil(t, table.Orders, "orders not initialized")
		assert.Empty(t, table.Orders, "orders not empty")

		tableRepoTable, err := tableRepo.FindByID(context.Background(), table.ID)
		require.NoError(t, err, "table not correctly saved")
		assert.Equal(t, table, tableRepoTable, "table not correctly saved")
	})

	t.Run("Failure", func(t *testing.T) {
		t.Run("Canceled Context", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := tableService.OpenTable(ctx)
			assert.Equal(t, domain.ErrorCode(err), domain.ECANCELED, "invalid error code")
		})
	})
}

func TestCloseTable(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			table    domain.Table
		}{
			{
				testName: "Table with no orders",
				table:    domain.Table{ID: id.New(), Status: domain.TableStatusOpened, Orders: make([]domain.Order, 0)},
			},
			{
				testName: "Table with orders",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusDone,
							Preparations: []domain.Preparation{
								{ID: id.New(), Status: domain.PreparationStatusServed, MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100}},
							},
						},
					},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				err := tableRepo.Save(context.Background(), tc.table)
				require.NoError(t, err, "initial setup failed")

				err = tableService.CloseTable(context.Background(), tc.table.ID)
				require.NoError(t, err, "close table failed")

				updatedTable, err := tableRepo.FindByID(context.Background(), tc.table.ID)
				require.NoError(t, err, "table not correctly saved")
				assert.Equal(t, domain.TableStatusClosed, updatedTable.Status, "invalid table status")
			})
		}
	})

	t.Run("Failures", func(t *testing.T) {
		tt := []struct {
			testName string
			table    domain.Table
			errCode  string
		}{
			{
				testName: "Table already closed",
				table:    domain.Table{ID: id.New(), Status: domain.TableStatusClosed, Orders: make([]domain.Order, 0)},
				errCode:  domain.EINVALID,
			},
			{
				testName: "Order not done / aborted",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID:     id.New(),
							Status: domain.OrderStatusTaken,
							Preparations: []domain.Preparation{
								{ID: id.New(), Status: domain.PreparationStatusPending, MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100}},
							},
						},
					},
				},
				errCode: domain.EINVALID,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				err := tableRepo.Save(context.Background(), tc.table)
				require.NoError(t, err, "Initial setup failed")

				err = tableService.CloseTable(context.Background(), tc.table.ID)

				require.Error(t, err, "close table should fail")
				assert.Equal(t, domain.ErrorCode(err), tc.errCode, "invalid error code")
			})
		}

		t.Run("Table Not Found", func(t *testing.T) {
			t.Parallel()

			err := tableService.CloseTable(context.Background(), id.New())

			assert.Equal(t, domain.ErrorCode(err), domain.ENOTFOUND, "invalid error code")
		})

		t.Run("Canceled Context", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			err := tableService.CloseTable(ctx, id.New())

			assert.Equal(t, domain.ErrorCode(err), domain.ECANCELED, "invalid error code")
		})
	})
}

func TestTakeOrder(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)

	t.Run("Success", func(t *testing.T) {
		tt := []struct {
			testName string
			table    domain.Table
			items    []domain.MenuItem
		}{
			{
				testName: "First order 1 item",
				table:    domain.Table{ID: id.New(), Orders: make([]domain.Order, 0), Status: domain.TableStatusOpened},
				items:    []domain.MenuItem{{ID: id.New(), Name: "item 1", Price: 100}},
			},
			{
				testName: "Second order 2 items",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{
							ID: id.New(), Status: domain.OrderStatusDone, Preparations: []domain.Preparation{
								{ID: id.New(), MenuItem: domain.MenuItem{ID: id.New(), Name: "item 1", Price: 100}, Status: domain.PreparationStatusServed},
								{ID: id.New(), MenuItem: domain.MenuItem{ID: id.New(), Name: "item 2", Price: 200}, Status: domain.PreparationStatusServed},
							},
						},
					},
				},
				items: []domain.MenuItem{
					{ID: id.New(), Name: "item 3", Price: 300},
					{ID: id.New(), Name: "item 4", Price: 400},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				err := tableRepo.Save(context.Background(), tc.table)
				require.NoError(t, err, "initial setup failed")

				order, err := tableService.TakeOrder(context.Background(), tc.table.ID, tc.items)

				require.NoError(t, err, "take order failed")
				assert.NotEqual(t, id.NilID(), order.ID, "generated order ID is nil")
				assert.Equal(t, domain.OrderStatusTaken, order.Status, "invalid order status")
				assert.Len(t, order.Preparations, len(tc.items), "invalid number of preparations")
				for _, p := range order.Preparations {
					assert.Equal(t, domain.PreparationStatusPending, p.Status, "invalid preparation status")
					assert.Contains(t, tc.items, p.MenuItem, "invalid preparation item")
				}

				updatedTable, err := tableRepo.FindByID(context.Background(), tc.table.ID)
				require.NoError(t, err, "table not correctly saved")
				assert.Len(t, updatedTable.Orders, len(tc.table.Orders)+1, "order not correctly saved")
			})
		}
	})

	t.Run("Failures", func(t *testing.T) {
		tt := []struct {
			testName string
			table    domain.Table
			items    []domain.MenuItem
			errCode  string
		}{
			{
				testName: "Empty items",
				table:    domain.Table{ID: id.New(), Orders: make([]domain.Order, 0), Status: domain.TableStatusOpened},
				items:    make([]domain.MenuItem, 0),
				errCode:  domain.EINVALID,
			},
			{
				testName: "Invalid item",
				table:    domain.Table{ID: id.New(), Orders: make([]domain.Order, 0), Status: domain.TableStatusOpened},
				items:    []domain.MenuItem{{ID: id.NilID(), Name: "test", Price: 100}},
				errCode:  domain.EINVALID,
			},
			{
				testName: "Table closed",
				table:    domain.Table{ID: id.New(), Orders: make([]domain.Order, 0), Status: domain.TableStatusClosed},
				items:    []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
				errCode:  domain.EINVALID,
			},
		}
		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				err := tableRepo.Save(context.Background(), tc.table)
				require.NoError(t, err, "Initial setup failed")

				_, err = tableService.TakeOrder(context.Background(), tc.table.ID, tc.items)

				require.Error(t, err, "take order should fail")
				assert.Equal(t, domain.ErrorCode(err), tc.errCode, "invalid error code")
			})
		}

		t.Run("Canceled Context", func(t *testing.T) {
			t.Parallel()

			items := []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := tableService.TakeOrder(ctx, id.New(), items)

			assert.Equal(t, domain.ErrorCode(err), domain.ECANCELED, "invalid error code")
		})

		t.Run("Table Not Found", func(t *testing.T) {
			t.Parallel()

			items := []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}}

			_, err := tableService.TakeOrder(context.Background(), id.New(), items)

			assert.Equal(t, domain.ErrorCode(err), domain.ENOTFOUND, "invalid error code")
		})
	})
}

func TestStartPreparation(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)

	t.Run("Success", func(t *testing.T) {
		tt := []struct {
			testName             string
			tableFromPreparation func(preparation domain.Preparation) domain.Table
			preparation          domain.Preparation
		}{
			{
				testName: "First preparation",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{ID: id.New(), Status: domain.OrderStatusTaken, Preparations: []domain.Preparation{preparation}},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					Status:   domain.PreparationStatusPending,
					MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
				},
			},
			{
				testName: "Second preparation",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{
								ID:     id.New(),
								Status: domain.OrderStatusTaken,
								Preparations: []domain.Preparation{
									{
										ID:       id.New(),
										Status:   domain.PreparationStatusReady,
										MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
									},
									preparation,
								},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					Status:   domain.PreparationStatusPending,
					MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				table := tc.tableFromPreparation(tc.preparation)
				err := tableRepo.Save(context.Background(), table)
				require.NoError(t, err, "Initial setup failed")

				err = tableService.StartPreparation(context.Background(), tc.preparation.ID)

				require.NoError(t, err, "start preparation failed")
				updatedTable, err := tableRepo.FindByID(context.Background(), table.ID)
				require.NoError(t, err, "table not correctly saved")

				preparation, order, err := updatedTable.ExtractPreparationWithOrder(tc.preparation.ID)
				require.NoError(t, err, "preparation not found")
				assert.Equal(t, domain.PreparationStatusInProgress, preparation.Status, "preparation status not updated")
				assert.Equal(t, domain.OrderStatusTaken, order.Status, "order should not be updated")
			})
		}
	})

	t.Run("Failures", func(t *testing.T) {
		tt := []struct {
			testName             string
			tableFromPreparation func(preparation domain.Preparation) domain.Table
			preparation          domain.Preparation
			errCode              string
		}{
			{
				testName: "Preparation not found",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: make([]domain.Order, 0),
					}
				},
				errCode: domain.ENOTFOUND,
			},
			{
				testName: "Preparation not pending",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{
								ID:           id.New(),
								Status:       domain.OrderStatusTaken,
								Preparations: []domain.Preparation{preparation},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					Status:   domain.PreparationStatusInProgress,
					MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
				},
				errCode: domain.EINVALID,
			},
			{
				testName: "Table closed",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusClosed,
						Orders: []domain.Order{
							{
								ID:           id.New(),
								Status:       domain.OrderStatusDone,
								Preparations: []domain.Preparation{preparation},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					Status:   domain.PreparationStatusServed,
					MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
				},
				errCode: domain.EINVALID,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				table := tc.tableFromPreparation(tc.preparation)
				err := tableRepo.Save(context.Background(), table)
				require.NoError(t, err, "Initial setup failed")

				err = tableService.StartPreparation(context.Background(), tc.preparation.ID)

				require.Error(t, err, "start preparation should fail")
				assert.Equal(t, domain.ErrorCode(err), tc.errCode, "invalid error code")
			})
		}

		t.Run("Canceled Context", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			err := tableService.StartPreparation(ctx, id.New())

			require.Error(t, err, "start preparation should fail")
			assert.Equal(t, domain.ErrorCode(err), domain.ECANCELED, "invalid error code")
		})
	})
}

func TestFinishPreparation(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)

	t.Run("Success", func(t *testing.T) {
		tt := []struct {
			testName             string
			tableFromPreparation func(preparation domain.Preparation) domain.Table
			preparation          domain.Preparation
		}{
			{
				testName: "Preparation in progress",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{
								ID:           id.New(),
								Status:       domain.OrderStatusTaken,
								Preparations: []domain.Preparation{preparation},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					Status:   domain.PreparationStatusInProgress,
					MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				table := tc.tableFromPreparation(tc.preparation)
				err := tableRepo.Save(context.Background(), table)
				require.NoError(t, err, "Initial setup failed")

				err = tableService.FinishPreparation(context.Background(), tc.preparation.ID)

				require.NoError(t, err, "finish preparation failed")

				updatedTable, err := tableRepo.FindByID(context.Background(), table.ID)
				require.NoError(t, err, "table not correctly saved")
				updatedOrder := updatedTable.Orders[0]
				updatedPreparation := updatedOrder.Preparations[0]
				assert.Equal(t, domain.PreparationStatusReady, updatedPreparation.Status, "preparation status not updated")
			})
		}
	})

	t.Run("Failures", func(t *testing.T) {
		tt := []struct {
			testName             string
			tableFromPreparation func(preparation domain.Preparation) domain.Table
			preparation          domain.Preparation
			errCode              string
		}{
			{
				testName: "Preparation not found",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: make([]domain.Order, 0),
					}
				},
				preparation: domain.Preparation{id.NilID(), domain.MenuItem{}, domain.PreparationStatusPending},
				errCode:     domain.ENOTFOUND,
			},
			{
				testName: "Preparation not in progress",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{
								ID:           id.New(),
								Status:       domain.OrderStatusTaken,
								Preparations: []domain.Preparation{preparation},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					Status:   domain.PreparationStatusPending,
					MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
				},
				errCode: domain.EINVALID,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				table := tc.tableFromPreparation(tc.preparation)
				err := tableRepo.Save(context.Background(), table)
				require.NoError(t, err, "Initial setup failed")

				err = tableService.FinishPreparation(context.Background(), tc.preparation.ID)

				require.Error(t, err, "finish preparation should fail")
				assert.Equal(t, domain.ErrorCode(err), tc.errCode, "invalid error code")
			})
		}

		t.Run("Canceled Context", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			err := tableService.FinishPreparation(ctx, id.New())

			require.Error(t, err, "finish preparation should fail")
			assert.Equal(t, domain.ErrorCode(err), domain.ECANCELED, "invalid error code")
		})
	})
}

func TestServePreparation(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)

	t.Run("Success", func(t *testing.T) {
		tt := []struct {
			testName             string
			tableFromPreparation func(preparation domain.Preparation) domain.Table
			preparation          domain.Preparation
			expectedOrderStatus  domain.OrderStatus
		}{
			{
				testName: "Last preparation",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{
								ID:           id.New(),
								Status:       domain.OrderStatusTaken,
								Preparations: []domain.Preparation{preparation},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					Status:   domain.PreparationStatusReady,
					MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
				},
				expectedOrderStatus: domain.OrderStatusDone,
			},
			{
				testName: "Not last preparation",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{
								ID:     id.New(),
								Status: domain.OrderStatusTaken,
								Preparations: []domain.Preparation{
									{
										ID:       id.New(),
										Status:   domain.PreparationStatusReady,
										MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
									},
									preparation,
								},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					Status:   domain.PreparationStatusReady,
					MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
				},
				expectedOrderStatus: domain.OrderStatusTaken,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				table := tc.tableFromPreparation(tc.preparation)
				err := tableRepo.Save(context.Background(), table)
				require.NoError(t, err, "Initial setup failed")

				err = tableService.ServePreparation(context.Background(), tc.preparation.ID)

				require.NoError(t, err, "serve preparation failed")

				updatedTable, err := tableRepo.FindByID(context.Background(), table.ID)
				require.NoError(t, err, "table not correctly saved")
				updatedOrder := updatedTable.Orders[0]
				assert.Equal(t, tc.expectedOrderStatus, updatedOrder.Status, "order status not updated")
			})
		}
	})

	t.Run("Failures", func(t *testing.T) {
		tt := []struct {
			testName             string
			tableFromPreparation func(preparation domain.Preparation) domain.Table
			preparation          domain.Preparation
			errCode              string
		}{
			{
				testName: "Preparation not found",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: make([]domain.Order, 0),
					}
				},
				preparation: domain.Preparation{},
				errCode:     domain.ENOTFOUND,
			},
			{
				testName: "Preparation not ready",
				tableFromPreparation: func(preparation domain.Preparation) domain.Table {
					return domain.Table{
						ID:     id.New(),
						Status: domain.TableStatusOpened,
						Orders: []domain.Order{
							{
								ID:           id.New(),
								Status:       domain.OrderStatusTaken,
								Preparations: []domain.Preparation{preparation},
							},
						},
					}
				},
				preparation: domain.Preparation{
					ID:       id.New(),
					Status:   domain.PreparationStatusPending,
					MenuItem: domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
				},
				errCode: domain.EINVALID,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				table := tc.tableFromPreparation(tc.preparation)
				err := tableRepo.Save(context.Background(), table)
				require.NoError(t, err, "Initial setup failed")

				err = tableService.ServePreparation(context.Background(), tc.preparation.ID)

				require.Error(t, err, "serve preparation should fail")
				assert.Equal(t, domain.ErrorCode(err), tc.errCode, "invalid error code")
			})
		}

		t.Run("Canceled Context", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			err := tableService.ServePreparation(ctx, id.New())

			require.Error(t, err, "serve preparation should fail")
			assert.Equal(t, domain.ErrorCode(err), domain.ECANCELED, "invalid error code")
		})
	})
}
