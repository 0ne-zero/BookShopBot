package database

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/0ne-zero/BookShopBot/internal/database/model"
	setting "github.com/0ne-zero/BookShopBot/internal/utils/settings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// If it couldn't connect to database, and also it didn't close program, returns nil
func InitializeOrGetDB() *gorm.DB {
	if db == nil {
		// DSN = Data source name (like connection string for database)
		dsn := setting.ReadFieldInSettingData("DSN")

		// For error handling
		var connect_again = true
		for connect_again {
			var err error
			db, err = connectDB(dsn)
			if err != nil {
				// Specific error handling

				//Databse doesn't exists, we have to create the database
				if strings.Contains(err.Error(), "Unknown database") {
					err = createDatabaseFromDSN(dsn)
					if err != nil {
						// Database isn't exists
						// Also we can't create database from dsn
						fmt.Printf("Mentioned database in dsn isn't created,program tried to create that database but it can't do that.\nError: %s\n", err.Error())
						os.Exit(1)
					}
					// Database created in mysql
					// Don't check rest of possible errors and try to connect again
					continue
				}
				// Error handling with error type detection
				switch err.(type) {
				case *net.OpError:
					op_err := err.(*net.OpError)
					// Get TCPAddr if exists
					if tcp_addr, ok := op_err.Addr.(*net.TCPAddr); ok {
						// Check error occurred when we trired to connect to mysql
						if tcp_addr.Port == 3306 {
							// Try to start mysql service
							connect_again = startMySqlService()
						}
					}
				default:
					log.Fatal("Cannot connect to database - " + err.Error() + "\nHint: Maybe you should start database service(deamon)")
				}
			} else {
				// We don't need to try again to connect to database because we are connected
				connect_again = false
			}
		}

		db.Set("gorm:auto_preload", true)
		return getDB()
	} else {
		return getDB()
	}
}
func connectDB(dsn string) (*gorm.DB, error) {
	// Connect to database with gorm
	return gorm.Open(
		// Open Databse
		mysql.New(mysql.Config{DSN: dsn}),
		// Config GORM
		&gorm.Config{
			// Allow create tables with null foreignkey
			DisableForeignKeyConstraintWhenMigrating: true,
			// All Datetime in database is in UTC
			NowFunc:              func() time.Time { return time.Now().UTC() },
			FullSaveAssociations: true,
		})
}

func getDB() *gorm.DB {
	return db
}
func MigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Address{},
		&model.Book{},
		&model.BookCoverType{},
		&model.BookSize{},
		&model.BookAgeCategory{},
		&model.Order{},
		&model.OrderStatus{},
		&model.Cart{},
		&model.CartItem{},
	)
}
func CreateEssensialData(db *gorm.DB) error {
	var err error
	// Book
	err = db.FirstOrCreate(&model.Book{
		Base: model.Base{ID: 1}, ISBN: "978402894626", Title: "Why", Author: "Me", Translator: "Me", PaperType: "Normal", Description: "See why?",
		NumberOfPages: 542, Genre: "History", Censored: false, Publisher: "You", PublishDate: "Now", Price: 430, GoodReadsScore: 4.9, ArezoScore: 5,
		BookCoverTypeID: 1, BookAgeCategoryID: 1, BookSizeID: 5}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.Book{
		Base: model.Base{ID: 2}, ISBN: "978402894626", Title: "Yeah", Author: "Me", Translator: "Me", PaperType: "Normal", Description: "Just yeah",
		NumberOfPages: 542, Genre: "History", Censored: false, Publisher: "You", PublishDate: "Now", Price: 430, GoodReadsScore: 4.9, ArezoScore: 5,
		BookCoverTypeID: 1, BookAgeCategoryID: 1, BookSizeID: 5}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.Book{
		Base: model.Base{ID: 3}, ISBN: "978402894626", Title: "Ofcourse", Author: "Me", Translator: "Me", PaperType: "Normal", Description: "Ofcourse",
		NumberOfPages: 542, Genre: "History", Censored: false, Publisher: "You", PublishDate: "Now", Price: 430, GoodReadsScore: 4.9, ArezoScore: 5,
		BookCoverTypeID: 1, BookAgeCategoryID: 1, BookSizeID: 5}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.Book{
		Base: model.Base{ID: 4}, ISBN: "978402894626", Title: "Ok", Author: "Me", Translator: "Me", PaperType: "Normal", Description: "Ok",
		NumberOfPages: 542, Genre: "History", Censored: false, Publisher: "You", PublishDate: "Now", Price: 430, GoodReadsScore: 4.9, ArezoScore: 5,
		BookCoverTypeID: 1, BookAgeCategoryID: 1, BookSizeID: 5}).Error
	if err != nil {
		return err
	}
	// Order status
	err = db.FirstOrCreate(&model.OrderStatus{Base: model.Base{ID: 1}, Status: "در صف تایید"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.OrderStatus{Base: model.Base{ID: 2}, Status: "در حال بسته بندی"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.OrderStatus{Base: model.Base{ID: 3}, Status: "در حال ارسال"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.OrderStatus{Base: model.Base{ID: 4}, Status: "تحویل داده شده"}).Error
	if err != nil {
		return err
	}
	// Cover type
	err = db.FirstOrCreate(&model.BookCoverType{Base: model.Base{ID: 1}, Type: "گالینور"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookCoverType{Base: model.Base{ID: 2}, Type: "شومیز"}).Error
	if err != nil {
		return err
	}
	// Book size
	err = db.FirstOrCreate(&model.BookSize{Base: model.Base{ID: 1}, Name: "جیبی"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookSize{Base: model.Base{ID: 2}, Name: "پالتویی"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookSize{Base: model.Base{ID: 3}, Name: "رقعی"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookSize{Base: model.Base{ID: 4}, Name: "وزیری"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookSize{Base: model.Base{ID: 5}, Name: "خشتی"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookSize{Base: model.Base{ID: 6}, Name: "رحلی کوچک"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookSize{Base: model.Base{ID: 7}, Name: "رحلی بزرگ"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookSize{Base: model.Base{ID: 8}, Name: "بیاضی"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookSize{Base: model.Base{ID: 9}, Name: "سلطانی"}).Error
	if err != nil {
		return err
	}
	// Book age category
	err = db.FirstOrCreate(&model.BookAgeCategory{Base: model.Base{ID: 1}, Category: "خردسال"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookAgeCategory{Base: model.Base{ID: 2}, Category: "خردسال"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookAgeCategory{Base: model.Base{ID: 3}, Category: "نوجوان"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookAgeCategory{Base: model.Base{ID: 4}, Category: "بزرگسال"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookAgeCategory{Base: model.Base{ID: 5}, Category: "میانسال"}).Error
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model.BookAgeCategory{Base: model.Base{ID: 6}, Category: "کهنسال"}).Error
	if err != nil {
		return err
	}
	return nil
}

// func CreateIfNotExists[m model.Model](model *m)error{

// }
func getDatabaseNameFromDSN(dsn string) string {
	// w = without
	w_user_pass_protocol_ip := dsn[strings.LastIndex(dsn, "/")+1:]
	return w_user_pass_protocol_ip[:strings.LastIndex(w_user_pass_protocol_ip, "?")]
}

func createDatabaseFromDSN(dsn string) error {
	// Create database
	dsn_without_database := dsn[:strings.LastIndex(dsn, "/")] + "/"
	db, err := sql.Open("mysql", dsn_without_database)
	if err != nil {
		if !startMySqlService() {
			fmt.Printf("We can't connect to mysql and we can't even start mysql.service\nError: %s\n", err.Error())
			os.Exit(1)
		}
		db, err = sql.Open("mysql", dsn_without_database)
		if err != nil {
			fmt.Printf("mysql.service is in start mode, but for any reason we can't connect to database\nError: %s\n", err.Error())
			os.Exit(1)
		}
	}
	db_name := getDatabaseNameFromDSN(dsn)
	_, err = db.Exec("CREATE DATABASE " + db_name)
	return err
}
func startMySqlService() bool {
	var service_names = []string{"mysqld.service", "mysql.service"}
	for i := range service_names {
		command := fmt.Sprintf("systemctl start %s", service_names[i])
		_, err := exec.Command("bash", "-c", command).Output()
		if err == nil {
			return true
		}
	}
	return false
}
