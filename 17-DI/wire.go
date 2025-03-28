//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/diogokimisima/17-DI/product"
	"github.com/google/wire"
)

var setRepositoryDepenndency = wire.NewSet(product.NewProductRepository,
	wire.Bind(product.ProductRepositoryInterface, new(*product.ProductRepository)))

func NewUseCase(db *sql.DB) *product.ProductUseCase {
	wire.Build(
		setRepositoryDepenndency,
		product.NewProductRepository)
	return &product.ProductUseCase{}
}
