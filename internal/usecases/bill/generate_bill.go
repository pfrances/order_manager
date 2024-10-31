package bill

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/money"
	"order_manager/internal/repositories"
)

type GenerateBill struct {
	menuRepo  repositories.MenuRepository
	orderRepo repositories.OrderRepository
	tableRepo repositories.TableRepository
}

func NewGenerateBill(
	menuRepo repositories.MenuRepository,
	orderRepo repositories.OrderRepository,
	tableRepo repositories.TableRepository,
) *GenerateBill {
	return &GenerateBill{
		menuRepo:  menuRepo,
		orderRepo: orderRepo,
		tableRepo: tableRepo,
	}
}

func (g *GenerateBill) Execute(tableID id.ID) (model.Bill, error) {

	bill := model.Bill{ID: id.NilID()}

	table := g.tableRepo.GetTable(tableID)
	if table == nil {
		return bill, fmt.Errorf("%w table %v not found", repositories.ErrNotFound, tableID)
	}

	for _, orderID := range table.OrderIDs {
		order := g.orderRepo.GetOrder(orderID)
		if order == nil {
			return bill, fmt.Errorf("order not found")
		}
		total := 0
		for _, menuItemID := range order.MenuItemIDs {
			menuItem := g.menuRepo.GetMenuItem(menuItemID)
			if menuItem == nil {
				return bill, fmt.Errorf("menu item not found")
			}
			total += menuItem.Price
			bill.Items = append(bill.Items, *menuItem)
		}
		bill.TotalAmount += money.Money(total)
	}

	g.tableRepo.UpdateTable(tableID, func(table *model.Table) error {
		table.BillID = bill.ID
		return nil
	})

	return bill, nil
}
