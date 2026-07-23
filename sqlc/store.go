package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
)

type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates new store acc
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// executeTx executes transaction function
func (s *Store) executeTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("rb: %v, fn: %v", rbErr, err)
		}
		return err
	}
	return tx.Commit()
}

type ExecuteTransfereParams struct {
	FromAccountID int64           `json:"fromAccountId"`
	ToAccountID   int64           `json:"toAccountId"`
	Amount        decimal.Decimal `json:"amount"`
}

type ExecuteTransferResult struct {
	ExecutedTransfer  Transfer `json:"executedTransfer"`
	FromAccountyEntry Entry    `json:"fromAccountyEntry"`
	ToAccountEntry    Entry    `json:"toAccountEntry"`
	FromAccount       Account  `json:"fromAccount"`
	ToAccount         Account  `json:"toAccount"`
}

func (s *Store) ExecuteTransfer(ctx context.Context, params ExecuteTransfereParams) (*ExecuteTransferResult, error) {
	var result ExecuteTransferResult

	err := s.executeTx(ctx, func(q *Queries) error {
		var txError error
		result.ExecutedTransfer, txError = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		})

		if txError != nil {
			return txError
		}

		result.FromAccountyEntry, txError = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount:    params.Amount.Neg(),
		})
		if txError != nil {
			return txError
		}

		result.ToAccountEntry, txError = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount:    params.Amount,
		})
		if txError != nil {
			return txError
		}

		if params.FromAccountID < params.ToAccountID {

			result.FromAccount, result.ToAccount, txError = transferUpdateBalance(ctx, transferUpdateBalanceParams{
				FirstAccountId:  params.FromAccountID,
				FirstAmount:     params.Amount.Neg(),
				SecondAccountId: params.ToAccountID,
				SecondAmount:    params.Amount,
			}, q)
			if txError != nil {
				return txError
			}

		} else {
			result.ToAccount, result.FromAccount, txError = transferUpdateBalance(ctx, transferUpdateBalanceParams{
				FirstAccountId:  params.ToAccountID,
				FirstAmount:     params.Amount,
				SecondAccountId: params.FromAccountID,
				SecondAmount:    params.Amount.Neg(),
			}, q)
			if txError != nil {
				return txError
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

type transferUpdateBalanceParams struct {
	FirstAccountId  int64
	FirstAmount     decimal.Decimal
	SecondAccountId int64
	SecondAmount    decimal.Decimal
}

func transferUpdateBalance(ctx context.Context, params transferUpdateBalanceParams, q *Queries) (firstAccount Account, secondAccount Account, err error) {
	firstAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: params.FirstAmount,
		ID:     params.FirstAccountId,
	})
	if err != nil {
		return
	}

	secondAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: params.SecondAmount,
		ID:     params.SecondAccountId,
	})
	if err != nil {
		return
	}

	return
}
