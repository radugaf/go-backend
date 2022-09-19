package db

import (
	"context"
	"testing"
	"time"

	"github.com/radugaf/simplebank/tools"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) BankAccount {
	user := createRandomUser(t)

	arg := CreateBankAccountParams{
		Owner:    user.Username,
		Balance:  tools.RandomMoney(),
		Currency: tools.RandomCurrency(),
	}
	account, err := testQueries.CreateBankAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateBankAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetBankAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetBankAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateBankAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateBankAccountParams{
		ID:      account1.ID,
		Balance: tools.RandomMoney(),
	}

	account2, err := testQueries.UpdateBankAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteBankAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteBankAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetBankAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.Empty(t, account2)
}

func TestListBankAccounts(t *testing.T) {
	var lastAccount BankAccount

	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListBankAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListBankAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
