package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/walknight/gotoko/database/seeders"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	//Welcome message
	fmt.Println("Welcome to " + appConfig.AppName)
	//run initilize router
	server.InitializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) InitializeDB(dbConfig DBConfig) {
	var err error
	//setup connection string database with DB Driver check
	if dbConfig.DBDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
		server.DB, err = gorm.Open(mysql.Open(DBURL))
	}
	if dbConfig.DBDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
		server.DB, err = gorm.Open(postgres.Open(DBURL))
	}
	//check if any error when connect to database
	if err != nil {
		fmt.Printf("Cannot connect to %s database", dbConfig.DBDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", dbConfig.DBDriver)
	}
}

func (server *Server) dbMigrate() {
	for _, model := range RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Data migrated successfully.")

}

func (server *Server) initCommands(config AppConfig, dbConfig DBConfig) {
	//initialize db connection
	server.InitializeDB(dbConfig)
	//define command
	cmdApp := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "db:migrate",
				Action: func(ctx *cli.Context) error {
					server.dbMigrate()
					return nil
				},
			},
			{
				Name: "db:seed",
				Action: func(ctx *cli.Context) error {
					err := seeders.DBSeed(server.DB)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
		},
	}
	//run command
	err := cmdApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// fungsi untuk get env dan set default bila param env tidak ada
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func Run() {
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error when loading .env file")
		return
	}
	//get env value
	appConfig.AppName = getEnv("APP_NAME", "GoToko")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "9000")

	dbConfig.DBDriver = getEnv("DB_DRIVER", "mysql")
	dbConfig.DBHost = getEnv("DB_HOST", "localhost")
	dbConfig.DBPort = getEnv("DB_PORT", "3306")
	dbConfig.DBName = getEnv("DB_NAME", "go_toko")
	dbConfig.DBUser = getEnv("DB_USER", "root")
	dbConfig.DBPassword = getEnv("DB_PASSWORD", "")

	flag.Parse()
	arg := flag.Arg(0)
	if arg != "" {
		server.initCommands(appConfig, dbConfig)
	} else {
		server.Initialize(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}
}
