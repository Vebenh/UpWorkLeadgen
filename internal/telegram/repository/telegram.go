package repository

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Db struct {
	Connection *gorm.DB
	Mutex      *sync.RWMutex
}

type Customer struct {
	ID             uint           `gorm:"primaryKey;autoIncrement:true"`
	TelegramID     int64          `gorm:"uniqueIndex;not null"`
	SearchQueries  []*SearchQuery `gorm:"foreignKey:TelegramID;references:TelegramID"`
	SearchInterval time.Duration
}

type SearchQuery struct {
	ID         uint   `gorm:"primaryKey;autoIncrement:true"`
	TelegramID int64  `gorm:"not null"`
	Query      string `gorm:"not null"`
}

func InitDB() (*Db, error) {
	Db := &Db{}

	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.Get("db.host"),
		viper.Get("db.port"),
		viper.Get("db.user"),
		viper.Get("db.password"),
		viper.Get("db.dbname"))

	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Customer{}, &SearchQuery{})
	if err != nil {
		return nil, err
	}

	Db.Connection = db
	Db.Mutex = &sync.RWMutex{}

	return Db, nil
}

func (db *Db) CreateCustomer(telegramID int64, searchInterval time.Duration) error {
	customer := &Customer{
		TelegramID:     telegramID,
		SearchInterval: searchInterval,
	}

	if err := db.Connection.Create(customer).Error; err != nil {
		return err
	}

	return nil
}

func (db *Db) EnsureCustomerExists(telegramID int64) error {
	var customer Customer
	if err := db.Connection.Where("telegram_id = ?", telegramID).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.CreateCustomer(telegramID, 0); err != nil {
				return fmt.Errorf("error creating customer: %w", err)
			}
			return nil
		}
		return fmt.Errorf("error getting customer by telegramID: %w", err)
	}

	return nil
}

func (db *Db) GetCustomerByTelegramID(telegramID int64) (*Customer, error) {
	var customer Customer
	err := db.Connection.Preload("SearchQueries").Where("telegram_id = ?", telegramID).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (db *Db) FetchAllCustomers() ([]*Customer, error) {
	customers := []*Customer{}

	return customers, nil
}

func (db *Db) GetUpdateTime(tgID int64) (time.Duration, error) {
	var customer Customer
	if err := db.Connection.Where("telegram_id = ?", tgID).First(&customer).Error; err != nil {
		return 0, err
	}
	return customer.SearchInterval, nil
}

func (db *Db) SetUpdateTime(tgID int64, interval time.Duration) error {
	return db.Connection.Model(&Customer{}).Where("telegram_id = ?", tgID).Update("search_interval", interval).Error
}

func (db *Db) GetSearchQueries(tgID int64) ([]SearchQuery, error) {
	var queries []SearchQuery
	if err := db.Connection.Where("telegram_id = ?", tgID).Find(&queries).Error; err != nil {
		return nil, err
	}
	return queries, nil
}

func (db *Db) AddSearchQuery(tgID int64, query string) error {
	newQuery := SearchQuery{TelegramID: tgID, Query: query}
	return db.Connection.Create(&newQuery).Error
}

func (db *Db) UpdateSearchQuery(id uint, newQuery string) error {
	return db.Connection.Model(&SearchQuery{}).Where("id = ?", id).Update("query", newQuery).Error
}
