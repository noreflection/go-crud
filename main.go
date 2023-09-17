package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var (
	db  *gorm.DB
	err error
)

func main() {
	_ = ensureDatabaseExists()
	// Initialize the database connection
	dsn := "host=localhost user=postgres password=password dbname=core-service-db port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}
	// Migrate the database
	db.AutoMigrate(&Customer{}, &Order{})

	// Create a new Gin router
	r := gin.Default()

	// Define CRUD routes for Customer
	r.GET("/customers", getCustomers)
	r.GET("/customers/:id", getCustomer)
	r.POST("/customers", createCustomer)
	r.PUT("/customers/:id", updateCustomer)
	r.DELETE("/customers/:id", deleteCustomer)

	// Define CRUD routes for Order
	r.GET("/orders", getOrders)
	r.GET("/orders/:id", getOrder)
	r.POST("/orders", createOrder)
	r.PUT("/orders/:id", updateOrder)
	r.DELETE("/orders/:id", deleteOrder)

	// Start the server
	r.Run(":8080")
}

func ensureDatabaseExists() error {
	fmt.Printf("ensuring")
	// Database connection parameters
	connectionString := "user=username dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if the database exists
	query := "SELECT 1 FROM pg_database WHERE datname = 'core-service-db'"
	var exists int
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if exists != 1 {
		// Create the database if it doesn't exist
		_, err = db.Exec("CREATE DATABASE core-service-db")
		if err != nil {
			return err
		}
		fmt.Println("Database 'core-service-db' created.")
	} else {
		fmt.Println("Database 'core-service-db' already exists.")
	}

	return nil
}

// ... (Other CRUD handler functions)
func getCustomers(c *gin.Context) {
	var customers []Customer
	db.Find(&customers)
	c.JSON(http.StatusOK, customers)
}

func createCustomer(c *gin.Context) {
	var customer Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&customer)
	c.JSON(http.StatusCreated, customer)
}

func getOrders(c *gin.Context) {
	var orders []Order
	db.Find(&orders)
	c.JSON(http.StatusOK, orders)
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&order)
	c.JSON(http.StatusCreated, order)
}

func getCustomer(c *gin.Context) {
	var customer Customer
	id := c.Param("id")

	if err := db.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func updateCustomer(c *gin.Context) {
	var customer Customer
	id := c.Param("id")

	if err := db.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&customer)
	c.JSON(http.StatusOK, customer)
}

func deleteCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer Customer

	if err := db.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	db.Delete(&customer)
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}

func getOrder(c *gin.Context) {
	var order Order
	id := c.Param("id")

	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func updateOrder(c *gin.Context) {
	var order Order
	id := c.Param("id")

	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&order)
	c.JSON(http.StatusOK, order)
}

func deleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order

	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	db.Delete(&order)
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
