package belajargolangdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id, name) VALUES ('subai','Subai')"
	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert success")
}

func TestQuery(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println("id:", id)
		fmt.Println("name:", name)
	}
}
func TestQueryComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var created_at time.Time
		var birth_date sql.NullTime
		var married bool
		err = rows.Scan(&id, &name, &email, &balance, &rating, &birth_date, &married, &created_at)
		if err != nil {
			panic(err)
		}
		fmt.Println("id:", id)
		fmt.Println("name:", name)
		if email.Valid {
			fmt.Println("email:", email.String)
		}
		fmt.Println("balance:", balance)
		fmt.Println("rating:", rating)
		if birth_date.Valid {
			fmt.Println("birth_date:", birth_date.Time)
		}
		fmt.Println("married:", married)
		fmt.Println("created_at:", created_at)
	}
}

func TestSQLInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #"
	password := "salah"

	script := "SELECT username FROM user WHERE username='" + username +
		"' AND password='" + password + "' LIMIT 1"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Login success", username)

	} else {
		fmt.Println("Gagal login")
	}
}
func TestSQLInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #"
	password := "admin"

	script := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"
	rows, err := db.QueryContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Login success", username)

	} else {
		fmt.Println("Gagal login")
	}
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "Subairi@gmail.com"
	comment := "Test KOMEN"
	script := "INSERT INTO comments(email, comment) VALUES (?,?)"
	result, err := db.ExecContext(ctx, script, email, comment)
	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Insert success with id", insertId)
}
func TestPrepareStatement(t *testing.T) {
	group := sync.WaitGroup{}
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO comments(email, comment) VALUES (?,?)"
	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}

	defer statement.Close()

	for i := 0; i < 10; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			email := "Subairi" + strconv.Itoa(i) + "@gmail.com"
			comment := "Test KOMEN ke-" + strconv.Itoa(i)
			result, err := statement.ExecContext(ctx, email, comment)
			if err != nil {
				panic(err)
			}

			insertId, err := result.LastInsertId()
			if err != nil {
				panic(err)
			}

			fmt.Println("Insert success with id", insertId)
		}()
	}
	group.Wait()

}

func TestTransaction(t *testing.T) {
	group := sync.WaitGroup{}
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	group.Add(1)
	go func() {
		defer group.Done()
		id := "joko" + time.Now().String()
		name := "Joko"
		script := "INSERT INTO customer(id, name) VALUES (?,?)"
		_, err := tx.ExecContext(ctx, script, id, name)
		if err != nil {
			panic(err)
		}
		fmt.Println("Insert customer success")
	}()

	for i := 0; i < 10; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			script := "INSERT INTO comments(email, comment) VALUES (?,?)"
			email := "Subairi" + strconv.Itoa(i) + "@gmail.com"
			comment := "Test KOMEN ke-" + strconv.Itoa(i)
			result, err := tx.ExecContext(ctx, script, email, comment)
			if err != nil {
				panic(err)
			}

			insertId, err := result.LastInsertId()
			if err != nil {
				panic(err)
			}

			fmt.Println("Insert success with id", insertId)
		}()
	}
	group.Wait()

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
}
