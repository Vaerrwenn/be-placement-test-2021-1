package transactioncontroller

import (
	"b-pay/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateTransactionForm is a struct to bind with the Transaction
// Creation form
type CreateTransactionForm struct {
	SavingID    uint   `form:"saving" binding:"required"`
	Type        string `form:"type" binding:"required"`
	Value       int64  `form:"value" binding:"required"`
	Description string `form:"desc"`
}

// returnErrorAndAbort returns a JSON with error key and text value.
// And then abort any other handlers.
func returnErrorAndAbort(ctx *gin.Context, code int, errorText string) {
	ctx.JSON(code, gin.H{
		"error": errorText,
	})
	ctx.Abort()
}

// CreateTransactionHandler handles Transaction creation
func CreateTransactionHandler(c *gin.Context) {
	var input CreateTransactionForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	if strings.ToUpper(input.Type) != "DEPOSIT" && strings.ToUpper(input.Type) != "WITHDRAWAL" {
		returnErrorAndAbort(c, http.StatusBadRequest, "Wrong Transaction Type. Must be only DEPOSIT or WITHDRAWAL.")
		return
	}

	if strings.ToUpper(input.Type) == "WITHDRAWAL" {
		input.Value = -input.Value
	}

	var saving models.Saving
	savingID := strconv.FormatUint(uint64(input.SavingID), 10)
	// Get the Saving Source
	source := saving.GetSavingByID(savingID)
	if source == nil {
		returnErrorAndAbort(c, http.StatusNotFound, "Could not find Saving.")
		return
	}
	// Balance calculation.
	newBalance := source.Balance + input.Value
	// fmt.Printf("\nSource Balance: %d\nInput Value: %d\nNew Balance: %d\n", source.Balance, input.Value, newBalance)

	if newBalance < 0 {
		returnErrorAndAbort(c, http.StatusNotAcceptable, "Balance can not be lower than 0.")
		return
	}

	transaction := models.Transaction{
		SavingID:    input.SavingID,
		Type:        input.Type,
		Value:       input.Value,
		Description: input.Description,
	}

	if err := transaction.Store(); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	// Change the balance of the source's balance.
	err := source.ChangeBalance(newBalance)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Transaction added successfully.",
	})
	return
}
