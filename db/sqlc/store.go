package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transaction

type Store struct {
	//Embedding the queries to perform the query operation in different model and this is called composition
	*Queries
	db *sql.DB
}

// Fn to create an store object

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db), // This New fn was created by sqlc
	}
}

// It take a context and a callback function as input, Then it will start a new db transaction and create new transcation with queries object.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// This function will commit or rollback based on the error return by call back function

	// tx,err :=  store.db.BeginTx(ctx, &sql.TxOptions{}) // This option allow us to set a custom isolation level for this transaction if we don't set it explicitly  then the default isolation level of the database  server will be used

	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {

		return err
	}

	q := New(tx) //This is new similar to sqlc , but here we are passing an transaction
	err = fn(q)  // Now we need to query the transaction

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil { // if the callback function return error then need to rollback the transaction
			return fmt.Errorf("TX Error: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	//Finally if  all the transaction are successful then commit the transaction

	return tx.Commit()

	// Note that this function is unexported, because it start with lowercase `e`; don't want the external package to call it directly
	//Instead we will provide an exported function for each specfic transaction

}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

/* TransferTx  perform money transfer from one account to another account.
it create a transfer record, add account entries and update account's balance within a single database transaction*/

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	// var err error

	/*Now here we can see , we are accessing  the result variable of the outer function from inside this callback function similar to arg varibale
	This make the callback function a closure, since go lack supports for generics type  closure is often is used  when we want to get the result from the callback function
	because the callback function itself does not know the exact type of the result it should return*/

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {

			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

		} else {

			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

		}

		return nil

	})

	return result, err

}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return

}
