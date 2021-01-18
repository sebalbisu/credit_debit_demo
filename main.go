package main

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sebalbisu/app/account"
	"github.com/shopspring/decimal"
)

func main() {

	e := echo.New()

	e.GET("/", GetBalanceAction)
	e.GET("/transactions", AllTransactionsAction)
	e.POST("/transactions", CreateTransactionAction)
	e.GET("/transactions/:id", FindTransaction)

	e.Logger.Fatal(e.Start(":1323"))
}

// GetBalanceAction gets the balance
func GetBalanceAction(c echo.Context) error {
	// BalanceJSON type for json response
	type BalanceJSON struct {
		Balance string `json:"balance"`
	}

	return c.JSON(http.StatusOK, BalanceJSON{Balance: account.Status.BalanceFormated()})
}

// AllTransactionsAction list all the transactions
func AllTransactionsAction(c echo.Context) error {
	return c.JSON(http.StatusOK, account.History)
}

// CreateTransactionAction create a transaction
func CreateTransactionAction(c echo.Context) error {
	type TransactionBody struct {
		Kind   string `json:"type"`
		Amount string `json:"amount"`
	}
	txBody := new(TransactionBody)
	if err := c.Bind(txBody); err != nil {
		return err
	}
	if txBody.Kind == "" || txBody.Amount == "" {
		return fmt.Errorf("type and amount parameter is required, not empty")
	}

	amount, err := decimal.NewFromString(txBody.Amount)
	if err != nil {
		return fmt.Errorf("amount is not a decimal")
	}

	txType := account.TransactionType(txBody.Kind)

	tx, err := account.CreateTx(txType, amount)
	if err != nil {
		return err
	}

	x, err := tx.MarshalJSON()
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(x))
}

// FindTransaction finds a tx by id
func FindTransaction(c echo.Context) error {
	id := c.Param("id")
	n := new(big.Int)
	idBig, ok := n.SetString(id, 10)
	if !ok {
		return fmt.Errorf("error not number")
	}
	tx := account.History.Find(*idBig)
	if tx == nil {
		return c.String(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, tx)
}
