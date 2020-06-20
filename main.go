package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int `json:"id"`
	Name  string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	createTb := `CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);`

	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table customers", err)
	}
}


func createCustomersHandler(c *gin.Context) {
	cust := Customer{}
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := db.QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3)  RETURNING id", cust.Name, cust.Email, cust.Status)

	err := row.Scan(&cust.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, cust)
}

func getCustomersByIdHandler(c *gin.Context) {
	id := c.Param("id")

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	cust := &Customer{}

	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cust)
}

func getCustomersHandler(c *gin.Context) {

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	customers := []*Customer{}
	for rows.Next() {
		cust := &Customer{}

		err := rows.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		customers = append(customers, cust)
	}

	c.JSON(http.StatusOK, customers)
}

func updateCustomersHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	cust := &Customer{}

	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := c.ShouldBindJSON(cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err = db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := stmt.Exec(id, cust.Name, cust.Email, cust.Status); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cust)
}

func deleteCustomersHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM customers WHERE id = $1")
	if err != nil {
		log.Fatal("can't prepare delete statement", err)
	}

	if _, err := stmt.Exec(id); err != nil {
		log.Fatal("can't execute delete statment", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func authMiddleware(c *gin.Context) {
	fmt.Println("start #middleware")
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you don't have authorization!!"})
		c.Abort()
		return
	}

	c.Next()

	fmt.Println("end #middleware")

}
func setupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/")

	api.Use(authMiddleware)

	api.POST("/customers", createCustomersHandler)
	api.GET("/customers/:id", getCustomersByIdHandler)
	api.GET("/customers", getCustomersHandler)
	api.PUT("/customers/:id", updateCustomersHandler)
	api.DELETE("/customers/:id", deleteCustomersHandler)

	return r
}

func main() {
	r := setupRouter()
	r.Run(":2019")
}
