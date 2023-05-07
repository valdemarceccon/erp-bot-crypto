package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	t := time.Now()
	err = bbClient.Run(0, t, t)
	if err != nil {
		log.Fatal(err)
	}

}
