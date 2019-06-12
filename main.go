package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pmoorani/booksAPI/config"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pmoorani/booksAPI/controllers"
	"github.com/pmoorani/booksAPI/database"
	"github.com/pmoorani/booksAPI/models"
)

var err error

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found!")
	}
}

func main() {
	// Load ENV Config
	conf := config.New()

	// Connection String
	connectionString := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", conf.DBConfig.DBUsername, conf.DBConfig.DBPassword, conf.DBConfig.DBName)
	database.DB, err = gorm.Open("mysql", connectionString)
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	//database.DB.DropTableIfExists(&models.Author{}, &models.Book{})
	database.DB.Debug().AutoMigrate(&models.Author{}, &models.Book{}, &models.User{}, &models.Claims{})
	/*database.DB.Create(&models.Author{Fname:"raj", Lname:"moorani"})
	database.DB.First(&author)

	database.DB.Create(&models.Book{Isbn:"443212", Title:"Some book 1", AuthorID: author.ID})
	database.DB.Create(&models.Book{Isbn:"534124", Title:"Some book 2", AuthorID: author.ID})*/

	router := gin.Default()
	router.Use(TokenAuthMiddleware())
	api := router.Group("/api")
	{
		books := api.Group("/books")
		{
			books.GET("/", controllers.GetAllBooks)
			books.GET("/:id", controllers.GetBook)
			books.POST("/", controllers.CreateBook)
			books.PUT("/:id", controllers.UpdateBook)
			books.DELETE("/:id", controllers.DeleteBook)
		}

		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)
	}
	router.Run()
}
