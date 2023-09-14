package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestUserGetAll(t *testing.T) {

	db, err := initDBToTest()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	// Khởi tạo Models và User
	models := New(db)
	user := User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	users, err := models.User.GetAll()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("Expected 1 user, but got %d", len(users))
	}

	if users[0].ID != user.ID || users[0].Email != user.Email {
		t.Fatalf("Unexpected user data")
	}

}

func TestInsert(t *testing.T) {
	db, err := initDBToTest()
	if err != nil {
		log.Println("ERR: ", err)
		t.Error(err)
	}

	defer db.Close()

	models := New(db)
	user := User{
		Email:     "trevor.sawler@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "verysecret",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newID, err := models.User.Insert(user)
	if err != nil {
		t.Error(err)
	}

	if newID == 0 {
		t.Error(errors.New("Cant insert new user"))
	}
	newUser1, err := models.User.GetOne(newID)
	if err != nil {
		t.Error(err)
	}

	newUser2, err := models.User.GetByEmail(newUser1.Email)
	if err != nil {
		t.Error(err)
	}

	err = models.User.DeleteByID(newUser2.ID)
	if err != nil {
		t.Error(err)
	}

}

func initDBToTest() (*sql.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println("Postgres not ready...")
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
