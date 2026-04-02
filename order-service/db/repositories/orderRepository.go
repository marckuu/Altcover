package repositories

import (
	domains2 "Altcover/order-service/domains"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var errOrderNotFound = errors.New("заказ не найдена")

type OrderRepository struct {
	connection *pgx.Conn
}

func (o *OrderRepository) AddOrder(ctx context.Context, order domains2.Order) error {
	query := `
	INSERT INTO order (status, user_id, orders_history_id)
	VALUES ($1, $2, $3);
`
	if _, err := o.connection.Exec(ctx, query, order.Status, order.UserID, order.OrdersHistoryID); err != nil {
		return fmt.Errorf("order repo -> add: %w", err)
	}

	return nil
}

func (o *OrderRepository) GetOrderByID(ctx context.Context, orderID int64) (domains2.Order, error) {
	query := `
	SELECT id, status, user_id, orders_history_id
	FROM order
	WHERE id = $1;
`
	resultRow := o.connection.QueryRow(ctx, query, orderID)

	var order domains2.Order

	if err := resultRow.Scan(&order.ID, &order.Status, &order.UserID, &order.OrdersHistoryID); err != nil {
		return domains2.Order{}, fmt.Errorf("order repo -> get by id: %w", err)
	}

	return order, nil
}

func (o *OrderRepository) GetOrder(ctx context.Context, offset int, limit int) ([]domains2.Order, error) {
	query := `
	SELECT id, status, user_id, orders_history_id
	FROM order
	OFFSET $1
	LIMIT $2;
`
	resultRows, err := o.connection.Query(ctx, query, offset, limit)
	if err != nil {
		return []domains2.Order{}, fmt.Errorf("order repo -> get all -> query: %w", err)
	}

	defer resultRows.Close()

	var orders []domains2.Order

	for resultRows.Next() {
		var order domains2.Order

		if err = resultRows.Scan(&order.ID, &order.Status, &order.UserID, &order.OrdersHistoryID); err != nil {
			return []domains2.Order{}, fmt.Errorf("order repo -> get all -> parsing: %w", err)
		}

		orders = append(orders, order)
	}

	if resultRows.Err() != nil {
		return []domains2.Order{}, fmt.Errorf("order repo -> get all -> query result: %w", err)
	}

	return orders, nil
}

func (o *OrderRepository) UpdateOrder(ctx context.Context, order domains2.Order) error {
	query := `
	UPDATE order
	SET status = $1, user_id = $2, orders_history_id = $3
	WHERE id = $4;
`
	tag, err := o.connection.Exec(ctx, query, order.Status, order.UserID, order.OrdersHistoryID, order.ID)
	if err != nil {
		return fmt.Errorf("order repo -> update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errOrderNotFound
	}

	return nil
}

func (o *OrderRepository) DeleteOrder(ctx context.Context, orderID int64) error {
	query := `
	DELETE FROM order
	WHERE id = $1;
`
	tag, err := o.connection.Exec(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("order repo -> delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errOrderNotFound
	}

	return nil
}
