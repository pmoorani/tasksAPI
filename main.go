package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/pmoorani/booksAPI/config"
	"github.com/pmoorani/booksAPI/middlewares"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pmoorani/booksAPI/controllers"
	"github.com/pmoorani/booksAPI/database"
	"github.com/pmoorani/booksAPI/models"
)

var err error

func init() {
	// Load .env file
	_ = godotenv.Load()
	log.Println("init()")
	//if err := godotenv.Load("sc.env"); err != nil {
	//	log.Print("No .env file found!")
	//	return
	//}
}

func main() {
	// Load ENV Config
	conf := config.New()

	// Connection String
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", conf.DBConfig.DBHost, conf.DBConfig.DBPort, conf.DBConfig.DBUsername, conf.DBConfig.DBPassword, conf.DBConfig.DBName, conf.DBConfig.DBSSLMode)
	database.DB, err = gorm.Open(conf.DBConfig.DBType, connectionString)
	if err != nil {
		panic(err.Error())
	}
	defer database.DB.Close()

	// Migrate the schema
	database.DB.Debug().AutoMigrate(&models.Author{}, &models.Book{}, &models.User{}, &models.Claims{})

	port := os.Getenv("PORT")

	router := gin.Default()
	router.Use(cors.Default())

	router.Use(middlewares.TokenAuthMiddleware())
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
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	if port == "" {
		router.Run(":8080")
	} else {
		router.Run(":" + port)
	}
}
