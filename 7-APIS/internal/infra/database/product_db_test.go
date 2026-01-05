package database

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/diogokimisima/goexpert/7-APIS/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	// Setup - create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	// Auto migrate the schema to create the products table
	db.AutoMigrate(&entity.Product{})

	// Initialize the product repository
	productDB := NewProduct(db)

	// Create a new product
	product, err := entity.NewProduct("Test Product", 10.0)
	assert.NoError(t, err)
	assert.NotNil(t, product)

	// Test the Create method
	err = productDB.Create(product)
	assert.NoError(t, err)

	// Verify that the product was created correctly by retrieving it
	productFound, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.NotNil(t, productFound)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, "Test Product", productFound.Name)
	assert.Equal(t, 10.0, productFound.Price)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(
			fmt.Sprintf("Product %d", i),
			rand.Float64()*100,
		)
		assert.NoError(t, err)
		db.Create(product)
	}
	productDB := NewProduct(db)
	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 21", products[0].Name)
}

func TestFindByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	productDB := NewProduct(db)

	// Create a test product
	product := &entity.Product{
		ID:    uuid.New(),
		Name:  "Test Product",
		Price: 100.0,
	}

	// Save the product to the database
	err = productDB.Create(product)
	assert.NoError(t, err)

	// Test case 1: Find existing product
	foundProduct, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.NotNil(t, foundProduct)
	assert.Equal(t, product.ID, foundProduct.ID)
	assert.Equal(t, product.Name, foundProduct.Name)
	assert.Equal(t, product.Price, foundProduct.Price)

	// Test case 2: Try to find non-existent product
	nonExistentID := uuid.New().String()
	nonExistentProduct, err := productDB.FindByID(nonExistentID)
	assert.Error(t, err)
	assert.Nil(t, nonExistentProduct)
	assert.Contains(t, err.Error(), "record not found")
}

func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	productDB := NewProduct(db)

	// Create a test product
	product, err := entity.NewProduct("Original Product", 50.0)
	assert.NoError(t, err)

	// Save the product to the database
	err = productDB.Create(product)
	assert.NoError(t, err)

	// Update the product
	product.Name = "Updated Product"
	product.Price = 75.0

	err = productDB.Update(product)
	assert.NoError(t, err)

	// Retrieve the updated product
	updatedProduct, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Updated Product", updatedProduct.Name)
	assert.Equal(t, 75.0, updatedProduct.Price)

	// Test case: Try to update a non-existent product
	nonExistentProduct := &entity.Product{
		ID:    uuid.New(),
		Name:  "Non-existent Product",
		Price: 100.0,
	}

	err = productDB.Update(nonExistentProduct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestProductDelete(t *testing.T) {
	// Setup do banco de dados
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	// Cria uma instância do repositório
	productDB := NewProduct(db)

	// Cria um produto para testar
	product, err := entity.NewProduct("Test Product", 10.0)
	assert.NoError(t, err)

	// Salva o produto no banco
	err = productDB.Create(product)
	assert.NoError(t, err)

	// Testa se consegue encontrar o produto
	productFound, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, product.ID, productFound.ID)

	// Testa a deleção do produto
	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)

	// Verifica se o produto foi realmente deletado
	_, err = productDB.FindByID(product.ID.String())
	assert.Error(t, err) // Deve retornar erro, pois o produto não existe mais

	// Testa tentativa de deleção de produto inexistente
	err = productDB.Delete("invalid-id")
	assert.Error(t, err)
}
