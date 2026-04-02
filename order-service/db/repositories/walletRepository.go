package repositories

import (
	domains2 "Altcover/order-service/domains"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var errWalletNotFound = errors.New("кошелёк не найдена")

type WalletRepository struct {
	connection *pgx.Conn
}

func (w *WalletRepository) AddWallet(ctx context.Context, wallet domains2.Wallet) error {
	query := `
	INSERT INTO wallet (balance, user_id)
	VALUES ($1, $2);
`
	if _, err := w.connection.Exec(ctx, query, wallet.Balance, wallet.UserID); err != nil {
		return fmt.Errorf("wallet repo -> add: %w", err)
	}

	return nil
}

func (w *WalletRepository) GetWalletByID(ctx context.Context, walletID int64) (domains2.Wallet, error) {
	query := `
	SELECT id, balance, user_id
	FROM wallet
	WHERE id = $1;
`
	resultRow := w.connection.QueryRow(ctx, query, walletID)

	var wallet domains2.Wallet

	if err := resultRow.Scan(&wallet.ID, wallet.Balance, wallet.UserID); err != nil {
		return domains2.Wallet{}, fmt.Errorf("wallet repo -> get by id: %w", err)
	}

	return wallet, nil
}

func (w *WalletRepository) GetWallet(ctx context.Context, offset int, limit int) ([]domains2.Wallet, error) {
	query := `
	SELECT id, balance, user_id
	FROM wallet
	OFFSET $1
	LIMIT $2;
`
	resultRows, err := w.connection.Query(ctx, query, offset, limit)
	if err != nil {
		return []domains2.Wallet{}, fmt.Errorf("wallet repo -> get all -> query: %w", err)
	}

	defer resultRows.Close()

	var wallets []domains2.Wallet

	for resultRows.Next() {
		var wallet domains2.Wallet

		if err = resultRows.Scan(&wallet.ID, wallet.Balance, wallet.UserID); err != nil {
			return []domains2.Wallet{}, fmt.Errorf("wallet repo -> get all -> parsing: %w", err)
		}

		wallets = append(wallets, wallet)
	}

	if resultRows.Err() != nil {
		return []domains2.Wallet{}, fmt.Errorf("wallet repo -> get all -> query result: %w", err)
	}

	return wallets, nil
}

func (w *WalletRepository) UpdateWallet(ctx context.Context, wallet domains2.Wallet) error {
	query := `
	UPDATE wallet
	SET balance = $1, user_id = $2
	WHERE id = $3;
`
	tag, err := w.connection.Exec(ctx, query, wallet.Balance, wallet.UserID, wallet.ID)
	if err != nil {
		return fmt.Errorf("wallet repo -> update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errWalletNotFound
	}

	return nil
}

func (w *WalletRepository) DeleteWallet(ctx context.Context, walletID int64) error {
	query := `
	DELETE FROM wallet
	WHERE id = $1;
`
	tag, err := w.connection.Exec(ctx, query, walletID)
	if err != nil {
		return fmt.Errorf("wallet repo -> delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errWalletNotFound
	}

	return nil
}
