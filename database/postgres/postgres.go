package postgres

import (
	"fmt"
	"github.com/MamushevArup/krisha-scraper/models"
	uuid2 "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type Sql struct {
	Db   *sqlx.DB
	id   uuid2.UUID
	User *models.User
}

var db = new(Sql)

// NewDB returns a singleton instance of *sql.DB
func NewDB() *Sql {
	initializeDB()
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
	db.Db, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error with connect to the database:", err)
	}
	if err = db.Db.Ping(); err != nil {
		log.Panic("Error with ping ", err)
	}
}
