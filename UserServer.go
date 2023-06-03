package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	Name string `json:"Name"`
}

var (
	hostPG     = os.Getenv("PG_HOST")
	portPG     = os.Getenv("PG_PORT")
	userPG     = os.Getenv("PG_USER")
	passwordPG = os.Getenv("PG_PASSWORD")
	dbnamePG   = os.Getenv("PG_DBNAME")
)

var users []user
var UsePG string = os.Getenv("PostgresDB")

func main() {
	if UsePG == "True" {
		i, err := strconv.Atoi(portPG)
		log.Print(i)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			hostPG, i, userPG, passwordPG, dbnamePG)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sts := `
INSERT INTO usertable(name, ts) VALUES('Miha', $1);`
		//id := 0
		result, err := db.Exec(sts, time.Now().Format(time.RFC3339))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result.RowsAffected())
	}
	router := gin.Default()
	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	router.Use(p.Instrument())
	router.GET("/hello", getHello)
	router.POST("/user", postUser)
	router.GET("/user", getUser)

	router.Run("localhost:8080")
}
func getMetrics(c *gin.Context) {

}
func postUser(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
	if UsePG == "True" {
		i, err := strconv.Atoi(portPG)
		log.Print(i)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			hostPG, i, userPG, passwordPG, dbnamePG)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sts := `INSERT INTO usertable(name, ts) VALUES($1, $2);`
		result, err := db.Exec(sts, newUser.Name, time.Now().Format(time.RFC3339))

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result.RowsAffected())

	} else {
		f, err := os.Create(`user.log`)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		buffer := bufio.NewWriter(f)
		buffer.WriteString(newUser.Name + ": " + time.Now().String())
	}
}
func getUser(c *gin.Context) {
	name := c.Param("Name")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	c.IndentedJSON(http.StatusOK, name)

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}
func getHello(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello Page")
}
