package repository

import (
	"github.com/op/go-logging"
	"gorm.io/gorm"
)

var log = logging.MustGetLogger("log")

type IDatabase interface {
	Exec(sql string, args ...interface{}) *gorm.DB
	Raw(sql string, args ...interface{}) *gorm.DB
}
