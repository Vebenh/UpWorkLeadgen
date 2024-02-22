package repository

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Db struct {
	Connection *gorm.DB
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

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.dbname"))

	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Customer{}, &SearchQuery{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Db) CreateCustomer(telegramID int64, searchInterval time.Duration) (*Customer, error) {
	customer := &Customer{
		TelegramID:     telegramID,
		SearchInterval: searchInterval,
	}

	if err := db.Connection.Create(customer).Error; err != nil {
		return nil, err
	}

	return customer, nil
}

func (db *Db) GetCustomerByTelegramID(telegramID int64) (*Customer, error) {
	var customer Customer
	err := db.Connection.Preload("SearchQueries").Where("telegram_id = ?", telegramID).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
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
