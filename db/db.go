package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.cloudfoundry.org/lager"
	"gopkg.in/retry.v1"

	mysqld "github.com/go-sql-driver/mysql"
	"github.com/tocy1/toggl/api/model"
	"github.com/tocy1/toggl/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariaDBDataStore struct {
	db     *gorm.DB
	logger lager.Logger
}

var ErrNotFound = errors.New("object not found")

func NewMariaDBDataStore(c config.MariaDB, logger lager.Logger) *MariaDBDataStore {
	l := logger.Session("mariadb")
	ds := MariaDBDataStore{logger: l}
	var err error

	dsn := ds.getDSN(c)
	for a := retry.Start(connectRetryStrategy(), nil); a.Next(); {
		ds.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		} else {
			ds.logger.Error("connect.retry", err)
		}
	}

	if err != nil {
		ds.logger.Fatal("connect", err)
	}

	err = ds.db.AutoMigrate(&model.Deck{})
	if err != nil {
		ds.logger.Fatal("migrate.Deck", err)
	}

	return &ds
}

func (ds *MariaDBDataStore) getDSN(c config.MariaDB) string {
	cfg := mysqld.NewConfig()
	cfg.DBName = c.Database
	cfg.ParseTime = true
	cfg.User = c.Username
	cfg.Passwd = c.Password
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%v:%v", c.Host, c.Port)
	dsn := cfg.FormatDSN()
	return dsn

}

func connectRetryStrategy() retry.Strategy {
	return retry.LimitTime(30*time.Second,
		retry.Exponential{
			Initial: 1 * time.Second,
			Factor:  1.5,
		},
	)
}

func (ds *MariaDBDataStore) CreateDeck(ctx context.Context, Deck model.Deck) (model.Deck, error) {
	tx := ds.db.WithContext(ctx).Create(&Deck)
	if tx.Error != nil {
		return Deck, tx.Error
	}

	return Deck, nil
}

func (ds *MariaDBDataStore) DeleteDeck(ctx context.Context, Deck model.Deck) error {
	return ds.db.WithContext(ctx).Delete(&Deck).Error
}

func (ds *MariaDBDataStore) GetDeck(ctx context.Context, deckID string) (deck model.Deck, err error) {
	tx := ds.db.WithContext(ctx).Where("id = ?", deckID).First(&deck)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return deck, ErrNotFound
		}
		return deck, tx.Error
	}

	return deck, nil
}

func (ds *MariaDBDataStore) UpdateDeck(ctx context.Context, deck model.Deck) error {
	tx := ds.db.WithContext(ctx).Where("id = ?", deck.Id).Save(&deck)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
