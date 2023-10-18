package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"sync"
)

type Sql struct {
	Db *sql.DB
}

var once sync.Once
var db = new(Sql)

// NewDB returns a singleton instance of *sql.DB
func NewDB() *Sql {
	once.Do(initializeDB)
	return db
}

func initializeDB() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, dbName)
	var err error
	db.Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error with connect to the database:", err)
	}
}
