package global

import (
	"awesomeProject/db"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"strings"
)

var DB *sqlx.DB

func InitDb() {
	log.Println("initialising db")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	DB, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln("failed to connect DB : ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalln("failed to ping DB : ", err)
	}
	log.Println("successfully pinged DB")
	db.MigrationUp(DB)
}

// Tx provides transaction wrapper
func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := DB.Beginx()
	if err != nil {
		log.Println("failed to begin tx : ", err)
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("failed to rollback tx : ", rollbackErr)
			}
			return
		}
		if err := tx.Commit(); err != nil {
			log.Println("failed to commit tx : ", err)
		}
	}()
	err = fn(tx)
	return err
}

func SetupBindVars(stmt, bindVars string, length int) string {
	bindVars += ","
	stmt = fmt.Sprintf(stmt, strings.Repeat(bindVars, length))
	return replaceSQL(strings.TrimSuffix(stmt, ","), "?")
}

func replaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}
