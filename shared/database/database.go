package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	// gorm requires a "dialect" is imported to communicate with postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/decentorganization/topaz/shared/models"
)

// Manager is used to access our database across the application
var Manager *gorm.DB

// DataManager ...
type DataManager interface {
	CreateApp(a *models.App) error
	GetApps(as *models.Apps, a *models.App) error
	GetApp(a *models.App) error
}

type dataManager struct {
	db *gorm.DB
}

func (m dataManager) CreateApp(a *models.App) error {
	return m.db.Create(&a).Error
}

func (m dataManager) GetApps(as *models.Apps, a *models.App) error {
	return m.db.Model(&a.User).Order("created_at").Related(&as).Error
}

func (m dataManager) GetApp(a *models.App) error {
	return m.db.Model(&a.User).Related(&a).Error
}

func init() {
	godotenv.Load()

	dbConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	manager, err := gorm.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	Manager = manager

	Manager.LogMode(os.Getenv("GO_ENV") != "production")

	if err := Manager.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	Manager.DB().SetMaxOpenConns(maxOpenConns)
	Manager.DB().SetMaxIdleConns(maxIdleConns)
}

// NewDataManager must be called after init, for now. should replace init
func NewDataManager(db *gorm.DB) DataManager {
	return dataManager{db}
}
