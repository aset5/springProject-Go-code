package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"strconv"
)

type Product struct {
	ID     uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	UserID uint    `json:"user_id"`
}

var db *gorm.DB

func initDatabase() {
	dsn := "host=localhost user=postgres password=andaset2005 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Almaty"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database!")
	}
	err = db.AutoMigrate(&Product{})
	if err != nil {
		panic("Failed to migrate database!")
	}
}

func getProducts(c *gin.Context) {
	var products []Product
	db.Find(&products)

	// Лог для отладки
	for _, p := range products {
		fmt.Printf("ID: %d | Name: %s | Price: %.2f\n", p.ID, p.Name, p.Price)
	}

	c.JSON(http.StatusOK, products)
}

func getProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var product Product
	if result := db.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	db.Create(&product)
	c.JSON(http.StatusCreated, product)
}

func updateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var product Product
	if db.First(&product, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.Name = input.Name
	product.Price = input.Price
	db.Save(&product)
	c.JSON(http.StatusOK, product)
}

func deleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if result := db.Delete(&Product{}, id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func main() {
	initDatabase()

	router := gin.Default()

	router.GET("/products", getProducts)
	router.GET("/products/:id", getProduct)
	router.POST("/products", createProduct)
	router.PUT("/products/:id", updateProduct)
	router.DELETE("/products/:id", deleteProduct)

	router.Run(":8082")
}
