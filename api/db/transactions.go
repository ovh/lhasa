package db

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// TransactionManager allows to control a DB transaction
type TransactionManager interface {
	DB() *gorm.DB
	Transaction(func(*gorm.DB) error, logrus.FieldLogger) error
}

type transactionManager struct {
	db *gorm.DB
}

// NewTransactionManager returns a new TransactionManager
func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

// DB returns the database backend
func (tm *transactionManager) DB() *gorm.DB {
	return tm.db
}

// Transaction embed a callback in a db transaction
// if the callback returns an error, the transaction is rollbacked
func (tm *transactionManager) Transaction(callback func(db *gorm.DB) error, log logrus.FieldLogger) error {
	log.Debug("opening transaction")
	tx := tm.db.Begin()
	tx.SetLogger(getLogger(log))
	if err := tx.Error; err != nil {
		return err
	}
	if err := callback(tx); err != nil {
		if err := tx.Rollback().Error; err != nil {
			log.WithError(err).Warn("transaction has not been rollbacked")
			return err
		}
		log.Warnf("transaction has been rollbacked")
		return err
	}
	log.Debug("committing transaction")
	err := tx.Commit().Error
	if err == nil {
		log.Debug("transaction has been committed")
	}
	return err
}
