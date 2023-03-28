package models

import (
	"fmt"
	"gin/pkg/logging"
	"gin/pkg/setting"
	"log"
	"time"

	// "github.com/go-sql-driver/mysql"
	// _ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func Setup() {
	var (
		err                                       error
		dbName, user, password, host, tablePrefix string
	)

	// dbType = sec.Key("TYPE").String()
	dbName = setting.DatabaseSetting.Name
	user = setting.DatabaseSetting.User
	password = setting.DatabaseSetting.Password
	host = setting.DatabaseSetting.Host
	tablePrefix = setting.DatabaseSetting.TablePrefix
	//使用gorm.Open连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	fmt.Println(dsn)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix, // 表名前缀
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logging.Fatal(err)
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logging.Fatal(err)
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	updateTimeStampForCreateCallback(db)
	updateTimeStampForUpdateCallback(db)
	// db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}

func updateTimeStampForCreateCallback(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("updateTimeStampForCreateCallback", func(db *gorm.DB) {
		if db.Statement.Schema != nil {
			newTime := time.Now().Unix()
			field := db.Statement.Schema.LookUpField("CreatedOn")
			if field != nil {
				field.Set(db.Statement.Context, db.Statement.ReflectValue, newTime)
			}
			field = db.Statement.Schema.LookUpField("ModifiedOn")
			if field != nil {
				field.Set(db.Statement.Context, db.Statement.ReflectValue, newTime)
			}
		}
	})
}

func updateTimeStampForUpdateCallback(db *gorm.DB) {
	db.Callback().Update().Before("gorm:update").Register("updateTimeStampForUpdateCallback", func(db *gorm.DB) {
		if db.Statement.Schema != nil {
			newTime := time.Now().Unix()
			field := db.Statement.Schema.LookUpField("ModifiedOn")
			if field != nil {
				field.Set(db.Statement.Context, db.Statement.ReflectValue, newTime)
			}
		}
	})
}

func deleteCallback(db *gorm.DB) {
	db.Callback().Delete().Before("gorm:delete").Register("update_deleted_at", func(db *gorm.DB) {
		if db.Statement.Schema != nil {
			field := db.Statement.Schema.LookUpField("DeletedOn")
			if field != nil {
				now := time.Now().Unix()
				db.Statement.SetColumn("DeletedOn", now)
				db.Where(db.Statement.Quote(field.DBName) + " IS NULL")
			} else {
				db = db.Where(db.Statement.Quote(field.DBName) + " IS NULL")
			}
		}
	})
}
