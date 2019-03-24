package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	// please, do not define constants like this in production
	DbHost     = "db"
	DbUser     = "postgres-dev"
	DbPassword = "password"
	DbName     = "dev"
	Migration  = `CREATE TABLE IF NOT EXISTS users (
id serial PRIMARY KEY,
first_name text NOT NULL,
last_name text NOT NULL,
email text NOT NULL,
password text NOT NULL,
created_time timestamp with time zone DEFAULT current_timestamp)`
)

type User struct {
	firstName string 	`json:"firstName" binding:"required"`
	lastName string `json:"lastName" binding:"required"`
	email string `json:"email" binding:"required"`
	password string `json:"password" binding:"required"`
	createdTime time.Time `json:"createdTime"`
}
// global database connection
var db *sql.DB

func GetALLUsers() ([]User, error) {
	const q = `SELECT first_name, last_name, email, password, created_time FROM users ORDER BY created_time DESC LIMIT 100`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	results := make([]User, 0)

	for rows.Next() {
		var firstName string
		var lastName string
		var email string
		var password string
		var createdTime time.Time
		// scanning the data from the returned rows
		err = rows.Scan(&firstName, &lastName, &email, &password, &createdTime)
		if err != nil {
			return nil, err
		}
		// creating a new result
		user := User{firstName, lastName, email, password, createdTime}
		log.Printf("User is %+v", user)
		results = append(results, user)
	}

	return results, nil
}

func CreateUser(user User) error {
	log.Printf("New user information is %+v", user)
	const q = `INSERT INTO users(first_name, last_name, email, password, created_time) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(q, user.firstName, user.lastName, user.email, user.password, user.createdTime)
	return err
}

func main() {
	var err error
	// create a router with a default configuration
	r := gin.Default()
	// endpoint to retrieve all posted bulletins
	r.GET("/get-all-users", func(context *gin.Context) {
		results, err := GetALLUsers()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
			return
		}
		context.JSON(http.StatusOK, results)
	})
	// endpoint to create a new bulletin
	r.POST("/create-new-user", func(context *gin.Context) {
		var user User
		log.Printf("context is %+v", context)
		// reading the request's body & parsing the json
		if context.Bind(&user) == nil {
			user.createdTime = time.Now()
			if err := CreateUser(user); err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
				return
			}
			context.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}
		// if binding was not successful, return an error
		context.JSON(http.StatusUnprocessableEntity, gin.H{"status": "invalid body"})
	})
	// open a connection to the database
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	// do not forget to close the connection
	defer db.Close()
	// ensuring the table is created
	_, err = db.Query(Migration)
	if err != nil {
		log.Println("failed to run migrations", err.Error())
		return
	}
	// running the http server
	log.Println("running..")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
