package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
	"github.com/valdemarceccon/crypto-bot-erp/backend/scrapper"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbConfig := repository.PostgresConfigFromEnv()
	fmt.Println(dbConfig.String())
	db, err := sql.Open("pgx", dbConfig.String())

	if err != nil {
		log.Fatal(err)
	}
	userRepo := repository.NewUserPsql(db)
	bbClient := scrapper.NewByBitScrapper(&userRepo)

	log.Println("Starting data fetching.")
	err = bbClient.Run()
	if err != nil {
		log.Fatal(err)
	}

}
