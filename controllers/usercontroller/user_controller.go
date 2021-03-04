package usercontroller

import (
	"b-pay/config/auth"
	"b-pay/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterForm binds the data from the Registration Form to the struct.
type RegisterForm struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// LoginForm binds the data from the Login form to the struct.
type LoginForm struct {
	Email      string `form:"email" binding:"required"`
	Password   string `form:"password" binding:"required"`
	Remembered bool   `form:"remember"`
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

// RegisterUserHandler inputs the form-data into the database.
//
// 1. The functions will get the form data.
// If there is an error, the func will send an error to the front.
//
// 2. Password will be encrypted.
//
// 3. Send the data to the model to be saved to the database.
func RegisterUserHandler(c *gin.Context) {
	// Check if User is already logged in.
	token := c.Request.Header.Get("token")
	if token != "" {
		returnErrorAndAbort(c, http.StatusForbidden, "User is already logged in.")
		return
	}
	// Binds the form-data to `input` variable
	var input RegisterForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	// Password encryption using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest,
			fmt.Sprintf("ERROR: Could not encrypt password. %s", err.Error()),
		)
		return
	}

	// Assign the input into user variable. Used for storing the data in database.
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	// Stores user data.
	if err := user.StoreUser(); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "ok",
		"msg":  "user successfully registered.",
	})
	return
}

// LoginHandler handles the Login feature.
//
// 1. Binds the data from the login form
//
// 2. Check if user with inputted email exists
//
// 3. Check password
//
// 4. Generate JWT Token
//
// 5. Send the Token to the Header.
func LoginHandler(c *gin.Context) {
	// Check whether user is logged in.
	token := c.Request.Header.Get("token")
	if token != "" {
		returnErrorAndAbort(c, http.StatusForbidden, "User is already logged in.")
		return
	}
	// Bind input from the Login form.
	var input LoginForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(input.Email) < 3 || len(input.Password) < 3 {
		returnErrorAndAbort(c, http.StatusBadRequest, "Email, and Password length must be more than 3")
		return
	}

	userEmail := models.User{
		Email: input.Email,
	}

	// Check if user with inputted email exists.
	user := userEmail.GetUserByEmail()
	if user == nil {
		returnErrorAndAbort(c, http.StatusNotFound,
			fmt.Sprintf("ERROR: User with email %s does not exist.", input.Email),
		)
		return
	}

	// Check if inputted password is the same as the User's stored password.
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(input.Password))
	if err != nil {
		returnErrorAndAbort(c, http.StatusUnauthorized, "Password invalid.")
		return
	}

	var expirationHours = 24
	if input.Remembered {
		// Expired in 1 year.
		expirationHours = 8760
	}

	// JWT generation
	jwtWrapper := auth.JwtWrapper{
		SecretKey:       os.Getenv("JWT_SECRET"),
		Issuer:          "AuthService",
		ExpirationHours: int64(expirationHours),
	}

	signedToken, err := jwtWrapper.GenerateToken(user.Email)
	if err != nil {
		returnErrorAndAbort(c, http.StatusInternalServerError, "Error signing token.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     signedToken,
		"userID":    user.ID,
		"userEmail": user.Email,
		"userName":  user.Name,
	})

	return
}
