package proxy

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Test_Proxy(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		proxy := NewProxy(ctx, "127.0.0.1", ":3306", true)
		proxy.EnableDecoding()
		err := proxy.Start("3336")
		if err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()

	time.Sleep(2 * time.Second)

	db, err := sql.Open("mysql", "root:secret@tcp(localhost:3336)/mercadobitcoin")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed ot connect to db: %s", err.Error())
	}

	type User struct {
		Id   int64
		Name string
	}

	sql := "SELECT id, name FROM users"
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatalf("Failed to query db: %s (%s)", sql, err.Error())
	}

	if rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			log.Fatalf("Failed to scan row: %s", err.Error())
		}

		log.Printf("User fetched, id: %d, name: %s", user.Id, user.Name)
	}

	if rows.Err(); err != nil {
		log.Fatalf("Failed fetch all data: %s", err.Error())
	}

	cancel()
	wg.Wait()
}
