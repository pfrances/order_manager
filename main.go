package main

import (
	"context"
	"fmt"
	"io"
	"order_manager/internal/domain"
	"order_manager/internal/http"
	"order_manager/internal/log"
	"order_manager/internal/sqlite"
	"os"
	"os/signal"
)

func run(ctx context.Context, stdout, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := log.New(log.Info, stdout, stderr)

	db, err := sqlite.NewDB("./db", logger)
	if err != nil {
		return err
	}

	tableRepository := sqlite.NewTable(db)
	menuRepository := sqlite.NewMenu(db)
	billRepository := sqlite.NewBill(db)

	tableService := domain.NewTableService(tableRepository)
	menuService := domain.NewMenuService(menuRepository)
	billService := domain.NewBillService(billRepository)

	server := http.NewServer(
		logger,
		tableService,
		menuService,
		billService,
	)

	return server.Run(ctx)
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
