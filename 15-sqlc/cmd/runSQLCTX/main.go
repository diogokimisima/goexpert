package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diogokimisima/15-sqlc/internal/db"

	_ "github.com/go-sql-driver/mysql"
)

type CourseDB struct {
	dbConn *sql.DB
	*db.Queries
}

func NewCourseDB(dbConn *sql.DB) *CourseDB {
	return &CourseDB{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

type CourseParams struct {
	ID          string
	Name        string
	Description sql.NullString
}

type CategoryParams struct {
	ID          string
	Name        string
	Description sql.NullString
}

func (c *CourseDB) callTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := c.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			return fmt.Errorf("error on rollback: %v, original error: %v", errRb, err)
		}
		return err
	}
	return tx.Commit()
}

func (c *CourseDB) CreateCourseAndCategory(ctx context.Context, arg CourseParams, argCategory CategoryParams) error {
	err := c.callTx(ctx, func(q *db.Queries) error {
		var err error
		err = q.CreateCategory(ctx, db.CreateCategoryParams{
			ID:          argCategory.ID,
			Name:        argCategory.Name,
			Description: argCategory.Description,
		})
		if err != nil {
			return err
		}
		err = q.CreateCourse(ctx, db.CreateCourseParams{
			ID:          arg.ID,
			Name:        arg.Name,
			Description: arg.Description,
			CategoryID:  argCategory.ID,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	// queries := db.New(dbConn)

	courseArgs := CourseParams{
		ID:          "123e4567-e89b-12d3-a456-426614174000",
		Name:        "Go Expert",
		Description: sql.NullString{String: "Go Expert Course", Valid: true},
	}

	categoryArgs := CategoryParams{
		ID:          "123e4567-e89b-12d3-a456-426614174001",
		Name:        "Backend",
		Description: sql.NullString{String: "Backend Course", Valid: true},
	}

	courseDB := NewCourseDB(dbConn)
	err = courseDB.CreateCourseAndCategory(ctx, courseArgs, categoryArgs)
	if err != nil {
		panic(err)
	}
	// categories, err := queries.ListCategories(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, category := range categories {
	// 	println(category.Name)
	// }
}
