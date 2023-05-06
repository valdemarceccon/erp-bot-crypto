package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/valdemarceccon/crypto-bot-erp/backend/scrapper"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbConfig := store.PostgresConfigFromEnv()
	fmt.Println(dbConfig.String())
	db, err := sql.Open("pgx", dbConfig.String())

	if err != nil {
		log.Fatal(err)
	}
	userStore := store.NewUserPsql(db)
	apiStore := store.NewApiKeyPsql(db)
	bbClient := scrapper.NewByBitScrapper(userStore, apiStore)

	log.Println("Starting data fetching.")
	err = bbClient.Run()
	if err != nil {
		log.Fatal(err)
	}

}
