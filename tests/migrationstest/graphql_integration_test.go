package tests

import (
	gql "OzonTestovoe/internal/graphql"
	"OzonTestovoe/internal/repository/InDb"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testSchema graphql.Schema
var repo *InDb.PsqlRepos

func TestMain(m *testing.M) {
	db := setupTestDB()
	repo = InDb.NewPsqlRepos(db)
	if err := runMigrations(db, "up"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	schema, err := gql.Schema(repo)
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
	testSchema = schema
	code := m.Run()
	if err := runMigrations(db, "down"); err != nil {
		log.Fatalf("failed to rollback migrations: %v", err)
	}
	os.Exit(code)
}

func setupTestDB() *sql.DB {
	connStr := "host=localhost port=4040 user=postgres password=pass dbname=ozontestovoe_test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("ping failed: %v", err)
	}
	return db
}

func runMigrations(db *sql.DB, direction string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../cmd/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if direction == "up" {
		return m.Up()
	}
	if direction == "down" {
		return m.Down()
	}
	return fmt.Errorf("unknown migration direction: %s", direction)
}

func TestCreatePostAndComment(t *testing.T) {
	query := `
	mutation {
		createPost(title:"Test", content:"Content", commentsareallowed:true, postAuthor:"Tester") {
			id
			title
		}
	}`
	params := graphql.Params{Schema: testSchema, RequestString: query}
	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		t.Fatalf("failed to create post: %v", res.Errors)
	}
	post := res.Data.(map[string]interface{})["createPost"].(map[string]interface{})
	postID := post["id"].(int)
	commentQuery := fmt.Sprintf(`
	mutation {
		createComment(content:"Hello", postId:%d, commentAuthor:"User") {
			id
			content
		}
	}`, postID)
	params.RequestString = commentQuery
	res = graphql.Do(params)
	if len(res.Errors) > 0 {
		t.Fatalf("failed to create comment: %v", res.Errors)
	}
	comment := res.Data.(map[string]interface{})["createComment"].(map[string]interface{})
	if comment["content"] != "Hello" {
		t.Fatalf("expected comment content Hello, got %v", comment["content"])
	}
	getQuery := fmt.Sprintf(`{ post(id:%d) { id title comments { content } } }`, postID)
	params.RequestString = getQuery
	res = graphql.Do(params)
	if len(res.Errors) > 0 {
		t.Fatalf("failed to query post: %v", res.Errors)
	}
	comments := res.Data.(map[string]interface{})["post"].(map[string]interface{})["comments"].([]interface{})
	if len(comments) != 1 || comments[0].(map[string]interface{})["content"] != "Hello" {
		t.Fatalf("expected comment Hello, got %+v", comments)
	}
}
