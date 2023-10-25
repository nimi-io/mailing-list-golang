package main

import (
	"database/sql"
	"log"
	"mailinlist/jsonapi"
	"mailinlist/mdb"
	"sync"

	_ "github.com/mattn/go-sqlite3"

	"github.com/alexflint/go-arg"
)

var args struct {
	DbPath   string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"`
}

func main() {
	arg.MustParse(&args)

	if args.DbPath == "" {
		args.DbPath = "list.db"
	}
	if args.BindJson == "" {
		args.BindJson = ":8080"
	}
	log.Printf("using database '%v'", args.DbPath)

	db, err := sql.Open("sqlite3", args.DbPath)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mdb.TryCreate(db)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		jsonapi.Serve(db, args.BindJson)
		wg.Done()
	}()
	wg.Wait()
}
