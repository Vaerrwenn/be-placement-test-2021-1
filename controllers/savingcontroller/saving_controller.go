package savingcontroller

import (
	"b-pay/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateSavingForm is a struct for Create a Saving.
type CreateSavingForm struct {
	Name string `form:"name" binding:"required"`
	PIN  string `form:"pin" binding:"required"`
}

// LoginSavingForm is a struct for accessing a Saving.
type LoginSavingForm struct {
	PIN string `form:"pin" binding:"required"`
}

// returnErrorAndAbort returns a JSON with "error": errorText in it. After that,
// it aborts and stop the running function.
//
// Takes Gin's context, the HTTP Code, and error text.
func returnErrorAndAbort(ctx *gin.Context, code int, errorText string) {
	ctx.JSON(code, gin.H{
		"error": errorText,
	})
	ctx.Abort()
}

// CreateSavingHandler handles Saving creation.
func CreateSavingHandler(c *gin.Context) {
	var input CreateSavingForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get UserID from Header.
	userID, err := strconv.ParseUint(c.Request.Header.Get("userID"), 10, 0)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, "User ID is not found.")
		return
	}

	if len(input.Name) < 3 {
		returnErrorAndAbort(c, http.StatusBadRequest, "Input name must be more than 3 characters")
		return
	}

	_, err = strconv.Atoi(input.PIN)
	if err != nil || len(input.PIN) != 6 {
		returnErrorAndAbort(c, http.StatusBadRequest, "PIN must be numeric with 6 digits.")
		return
	}

	// Encrypt password for an account.
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(input.PIN), bcrypt.DefaultCost)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest,
			fmt.Sprintf("ERROR: Could not encrypt password. %s", err.Error()),
		)
		return
	}

	saving := models.Saving{
		UserID:  uint(userID),
		Name:    input.Name,
		Balance: 0,
		PIN:     hashedPIN,
	}

	if err := saving.Store(); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": saving.ID,
		"msg":  "Saving data is stored successfully.",
	})
	return
}

// IndexSavingHandler handles Savings Index. Shows all of Saving accounts that
// a user has.
func IndexSavingHandler(c *gin.Context) {
	var saving models.Saving

	userID := c.Request.Header.Get("userID")

	result, err := saving.GetSavingsByUserID(userID)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
		"qty":  len(*result),
	})
	return
}

// LoginSavingHandler handles the login process for Saving account before accessing the Saving account.
func LoginSavingHandler(c *gin.Context) {
	savingID := c.Param("id")
	if savingID == "" {
		returnErrorAndAbort(c, http.StatusBadRequest, "Saving ID is empty")
	}

	var input LoginSavingForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := strconv.Atoi(input.PIN)
	if len(input.PIN) != 6 || err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, "PIN must be 6 DIGITS long.")
		return
	}

	var saving models.Saving
	savingPIN := saving.GetPINBySavingID(savingID)

	err = bcrypt.CompareHashAndPassword([]byte(savingPIN), []byte(input.PIN))
	if err != nil {
		returnErrorAndAbort(c, http.StatusForbidden, "PIN is incorrect.")
		return
	}

	// Used as the key to unlock Saving account. Similar to token.
	key := fmt.Sprintf("key_%s", savingPIN)

	c.JSON(http.StatusOK, gin.H{
		"data": savingID,
		"key":  key,
		"msg":  "Successfully logging in to Saving account",
	})
	return
}
