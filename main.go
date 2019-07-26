package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pmoorani/booksAPI/config"
	"github.com/pmoorani/booksAPI/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pmoorani/booksAPI/controllers"
	"github.com/pmoorani/booksAPI/database"
	"github.com/pmoorani/booksAPI/models"

	uuid "github.com/satori/go.uuid"
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
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", conf.DBConfig.DBHost, conf.DBConfig.DBPort, conf.DBConfig.DBUsername, conf.DBConfig.DBPassword, conf.DBConfig.DBName)
	database.DB, err = gorm.Open(conf.DBConfig.DBType, connectionString)
	if err != nil {
		panic(err.Error())
	}
	defer database.DB.Close()

	// Parsing UUID from string input
	u2, err := uuid.FromString("e76d832c-ae61-4a68-8615-f8942ec64c66")
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	fmt.Printf("Successfully parsed: %s", u2)

	// Migrate the schema
	//database.DB.Debug().DropTableIfExists(&models.User{})
	//database.DB.Debug().DropTableIfExists(&models.Status{})

	database.DB.Debug().AutoMigrate(&models.Author{}, &models.Book{})
	database.DB.Debug().AutoMigrate(&models.User{}, &models.Claims{})
	database.DB.Debug().AutoMigrate(&models.Status{}, &models.Priority{})
	database.DB.Debug().AutoMigrate(&models.Task{})
	//database.DB.Debug().Create(&models.Status{Status: "Backlog", StatusDE: "Auftragsbestand"})
	//database.DB.Debug().Create(&models.Status{Status: "In Progress", StatusDE: "In Bearbeitung"})
	//database.DB.Debug().Create(&models.Status{Status: "Completed", StatusDE: "Abgeschlossen"})
	//database.DB.Debug().Create(&models.Priority{Priority: "Low", PriorityDE: "Niedrig"})
	//database.DB.Debug().Create(&models.Priority{Priority: "Medium", PriorityDE: "Mittel"})
	//database.DB.Debug().Create(&models.Priority{Priority: "High", PriorityDE: "High"})

	port := os.Getenv("PORT")

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	router.Use(middlewares.TokenAuthMiddleware())
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg":     "Welcome to our world!",
				"success": 1,
			})
		})

		books := api.Group("/books")
		{
			books.GET("/", controllers.GetAllBooks)
			books.GET("/:id", controllers.GetBook)
			books.POST("/", controllers.CreateBook)
			books.PUT("/:id", controllers.UpdateBook)
			books.DELETE("/:id", controllers.DeleteBook)
		}

		tasks := api.Group("/tasks")
		{
			tasks.GET("/", controllers.GetAllTasks)
			tasks.GET("/:uuid", controllers.GetTask)
			tasks.POST("/", controllers.CreateTask)
			tasks.PUT("/:uuid", controllers.UpdateTask)
			tasks.DELETE("/:uuid", controllers.DeleteTask)
		}

		users := api.Group("/users")
		{
			users.GET("/", controllers.GetAllUsers)
			users.GET("/:uuid", controllers.GetUser)
			users.PUT("/:uuid", controllers.UpdateUser)
			users.DELETE("/:uuid", controllers.DeleteUser)
		}

		statuses := api.Group("/statuses")
		{
			statuses.GET("/", controllers.GetAllStatuses)
		}

		priorities := api.Group("/priorities")
		{
			priorities.GET("/", controllers.GetAllPriorities)
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
