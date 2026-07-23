package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"tbo-go-bank/util"
	"testing"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

const (
	dbdriver = "postgres"
	dbsource = "postgresql://root:secret@localhost:5432/bank?sslmode=disable"
)

var testQueries *Queries
var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open(dbdriver, dbsource)
	if err != nil {
		log.Fatal("can not connect to db: ", err)
	}

	testQueries = New(db)
	os.Exit(m.Run())
}

func createTestAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		OwnerName: util.RandomOwnerName(),
		Balance:   decimal.New(2000000000000, -6),
		Currency:  util.RandomCurrencyCode(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Fatalf("failed to create account with params owner: %s, balance: %v,  currency code: %v ", arg.OwnerName, arg.Balance, arg.Currency)
	}

	return acc
}
