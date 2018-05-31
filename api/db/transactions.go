package db

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// TransactionManager allows to control a DB transaction
type TransactionManager interface {
	DB() *gorm.DB
	Transaction(func(*gorm.DB) error) error
}

type transactionManager struct {
	db  *gorm.DB
	log *logrus.Logger
}

// NewTransactionManager returns a new TransactionManager
func NewTransactionManager(db *gorm.DB, log *logrus.Logger) TransactionManager {
	return &transactionManager{db: db, log: log}
}

// DB returns the database backend
func (tm *transactionManager) DB() *gorm.DB {
	return tm.db
}

// Transaction embed a callback in a db transaction
// if the callback returns an error, the transaction is rollbacked
func (tm *transactionManager) Transaction(callback func(db *gorm.DB) error) error {
	tm.log.Debugf("opening transaction")
	tx := tm.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	if err := callback(tx); err != nil {
		if err := tx.Rollback().Error; err != nil {
			tm.log.WithError(err).Warnf("transaction has not been rollbacked")
			return err
		}
		tm.log.Warnf("transaction has been rollbacked")
		return err
	}
	return tx.Commit().Error
}
