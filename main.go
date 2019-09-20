package main

import (
	"fmt"
	"github.com/pmoorani/tasksAPI/models"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pmoorani/tasksAPI/config"
	"github.com/pmoorani/tasksAPI/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pmoorani/tasksAPI/controllers"
	"github.com/pmoorani/tasksAPI/database"

	uuid "github.com/satori/go.uuid"
)

var err error

func init() {
	fmt.Println("init()")
	env := os.Getenv("TMS_ENV")
	fmt.Printf(env)
	if "development" == env {
		godotenv.Load(".env.local")
	}

	if "production" == env {
		godotenv.Load()
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	// Load ENV Config
	conf := config.New()

	// Connection String
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", conf.DBConfig.DBHost, conf.DBConfig.DBPort, conf.DBConfig.DBUsername, conf.DBConfig.DBPassword, conf.DBConfig.DBName, conf.DBConfig.DBSSLMode)
	fmt.Println(connectionString)
	database.DB, err = gorm.Open(conf.DBConfig.DBType, connectionString)
	if err != nil {
		panic(err.Error())
	}
	defer database.DB.Close()

	database.DB.LogMode(true)
	// Parsing UUID from string input
	u2, err := uuid.FromString("e76d832c-ae61-4a68-8615-f8942ec64c66")
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	fmt.Printf("Successfully parsed: %s", u2)

	// Migrate the schema

	if !database.DB.HasTable(&models.User{}) {
		database.DB.CreateTable(&models.User{})
	}

	if !database.DB.HasTable(&models.Claims{}) {
		database.DB.CreateTable(&models.Claims{})
	}

	if !database.DB.HasTable(&models.Task{}) {
		database.DB.CreateTable(&models.Task{})

	}

	if !database.DB.HasTable(&models.Status{}) {
		database.DB.CreateTable(&models.Status{})
		database.DB.Debug().Create(&models.Status{Name: "Backlog", NameDE: "Auftragsbestand"})
		database.DB.Debug().Create(&models.Status{Name: "InProgress", NameDE: "In Bearbeitung"})
		database.DB.Debug().Create(&models.Status{Name: "Completed", NameDE: "Abgeschlossen"})
	}

	if !database.DB.HasTable(&models.Priority{}) {
		database.DB.CreateTable(&models.Priority{})
		database.DB.Debug().Create(&models.Priority{Name: "Low", NameDE: "Niedrig"})
		database.DB.Debug().Create(&models.Priority{Name: "Medium", NameDE: "Mittel"})
		database.DB.Debug().Create(&models.Priority{Name: "High", NameDE: "High"})
	}

	// database.DB.Debug().DropTableIfExists(&models.User{})
	//database.DB.Debug().DropTableIfExists(&models.Task{})
	//database.DB.Debug().DropTableIfExists(&models.Status{})
	//database.DB.Debug().DropTableIfExists(&models.Priority{})

	//database.DB.Debug().CreateTable(&models.Author{}, &models.Book{})
	//database.DB.Debug().CreateTable(&models.User{}, &models.Claims{})
	// database.DB.Debug().AutoMigrate(&models.Status{}, &models.Priority{})
	// database.DB.Debug().AutoMigrate(&models.Task{})

	//database.DB.Debug().Create(&models.Status{Name: "Backlog", NameDE: "Auftragsbestand"})
	//database.DB.Debug().Create(&models.Status{Name: "InProgress", NameDE: "In Bearbeitung"})
	//database.DB.Debug().Create(&models.Status{Name: "Completed", NameDE: "Abgeschlossen"})
	//database.DB.Debug().Create(&models.Priority{Name: "Low", NameDE: "Niedrig"})
	//database.DB.Debug().Create(&models.Priority{Name: "Medium", NameDE: "Mittel"})
	//database.DB.Debug().Create(&models.Priority{Name: "High", NameDE: "High"})

	port := os.Getenv("PORT")

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.TokenAuthMiddleware())

	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg":     "Welcome to our world!",
				"host": os.Getenv("HOSTNAME"),
				"success": 1,
			})
		})

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
