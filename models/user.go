package models

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/pmoorani/booksAPI/database"

	"strings"

	u "github.com/pmoorani/booksAPI/utils"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	BaseModel
	Username  string `json:"username" validate:"required" gorm:"unique_index"`
	Password  string `json:"password,omitempty" validate:"required"`
	Email     string `json:"email" gorm:"type:varchar(100);unique_index"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Token     string `json:"token" gorm:"-"`
	Tasks     []Task `json:"tasks"`
}

// Create a struct that will be encoded to a JWT
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserId          uuid.UUID `json:"user_id"`
	Username        string    `json:"username"`
	IsAuthenticated bool      `gorm:"-"`
	jwt.StandardClaims
}

// Validate incoming user details
func (user *User) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(user.Email, "@") {
		fmt.Println("Email address is required")

		return u.Message(false, "Email address is required"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password should be min 6 characters"), false
	}

	if exists := database.DB.Where("username = ?", user.Username).First(&user); exists.RowsAffected != 0 {
		return u.Message(false, "Username already exists!"), false
	}

	return u.Message(false, "Requirement passed"), true
}

var users []User

func AllUsers() ([]User, error) {
	err := database.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func FindUserByID(id interface{}) (User, error) {
	var user User
	err := database.DB.Where("id = ?", id).Find(&user).Error
	fmt.Println(err)

	if err != nil {
		return User{}, err
	}
	return user, nil
}
