package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/sergey-kvachenok/go-hello/db/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount (t *testing.T) Accounts {
	arg := CreateAccountParams{
		Owner: utils.RandomOwner(),
	Balance:  utils.RabdomMoney(),
	Currency: utils.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, arg.Owner, acc.Owner)
	require.Equal(t, arg.Currency, acc.Currency)
	require.Equal(t, arg.Balance, acc.Balance)

	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)

	return acc
}

func TestCreateAccount(t *testing.T) {
createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
		arg := UpdateAccountParams{
		ID: account1.ID,
		Balance:  utils.RabdomMoney(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)

	require.Equal(t, account2.Balance, arg.Balance)
}


func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
		
	 err := testQueries.DeleteAccount(context.Background(), account1.ID)
		require.NoError(t, err)

	account2, err :=  testQueries.GetAccount(context.Background(), account1.ID)
require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {
	var lastAccount Accounts

	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}