package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pmoorani/booksAPI/config"
	"github.com/pmoorani/booksAPI/database"

	"github.com/pmoorani/booksAPI/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Register(c *gin.Context) {
	validate := validator.New()
	fmt.Println("inside register()")
	var user models.User
	// Get the JSON body and decode into credentials
	err := c.ShouldBindJSON(&user)
	fmt.Println(&user)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP err
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Something went wrong!"})
		return
	}

	err = validate.Struct(user)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"msg": validationErrors.Error()})
		return
	}

	if resp, ok := user.Validate(); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": resp, "success": 0})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err = database.DB.Debug().Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Some error occurred!", "success": 0})
		return
	}

	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{"msg": "User account has been created.", "user": user, "success": 1})
}

func Login(c *gin.Context) {
	validate = validator.New()
	conf := config.New()
	SecretKey := conf.SecretKey

	var user models.User
	var userFromDB models.User

	// Get the JSON body and decode into credentials
	err := c.ShouldBindJSON(&user)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP err
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Something went wrong!"})
		return
	}

	err = validate.Struct(user)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"msg": validationErrors.Error()})
		return
	}

	// Get all matched records
	fmt.Printf("Username = %s", user.Username)
	err = database.DB.Debug().Where("username = ?", user.Username).First(&userFromDB).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"msg": "User not found!", "success": 0})
		}
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Something went wrong!"})
		return
	}
	fmt.Println(&userFromDB)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userFromDB.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid Username/Password"})
		return
	}

	userFromDB.Password = ""


	// Declare the expiration time of the token
	// here, we have kept it as 30 minutes
	expirationTime := time.Now().Add(30 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &models.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT String
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	userFromDB.Token = tokenString

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	c.JSON(http.StatusOK, gin.H{
		"user": userFromDB,
		"token": tokenString,
		"expires": expirationTime,
	})
}

