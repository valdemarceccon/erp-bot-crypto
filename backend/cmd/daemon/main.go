package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/hirokisan/bybit/v2"
)

type ByBitScrapper struct {
	userRepo repository.User
}

func NewByBitScrapper(userRepo *repository.User) *ByBitScrapper {
	return &ByBitScrapper{
		userRepo: *userRepo,
	}
}

func (bb *ByBitScrapper) Run() error {
	apiKeyList, err := bb.userRepo.ListActiveApiKeys(0)

	if err != nil {
		log.Printf("ERROR: could not get active api keys: %v\n", err)
	}

	log.Printf("INFO: found %d active keys\n", len(apiKeyList))
	bb.ScrapClonedPnL(apiKeyList)

	return nil
}

func (bb *ByBitScrapper) ScrapClonedPnL(apiKeyList []model.ApiKey) {
	now := time.Now()
	days := 10
	yesterdayStart := time.Date(now.Year(), now.Month(), now.Day()-days, 0, 0, 0, 0, time.UTC)
	yesterdayEnd := time.Date(now.Year(), now.Month(), now.Day()-days, 23, 59, 59, 999999999, time.UTC)

	yesterdayStartMs := yesterdayStart.UnixNano() / int64(time.Millisecond)
	yesterdayEndMs := yesterdayEnd.UnixNano() / int64(time.Millisecond)
	for _, apiKey := range apiKeyList {
		bybitClient := bybit.NewClient().WithAuth(apiKey.ApiKey, apiKey.ApiSecret)
		symbol := bybit.SymbolV5ETHUSD
		more := true
		cursor := ""
		for more {
			resp, err := bybitClient.V5().Position().GetClosedPnL(bybit.V5GetClosedPnLParam{
				Category:  bybit.CategoryV5Inverse,
				Symbol:    &symbol,
				StartTime: &yesterdayStartMs,
				EndTime:   &yesterdayEndMs,
				Cursor:    &cursor,
			})

			if err != nil {
				log.Println("ERROR: Was not able to get", apiKey.Id, "closed p&l:", err)
				continue
			}
			log.Printf("INFO: there is '%d' register for the key_id='%d' and user_id='%d'\n", len(resp.Result.List), apiKey.Id, apiKey.UserId)

			err = bb.userRepo.SaveClosedPnL(apiKey.UserId, apiKey.Id, resp.Result.List)
			cursor = resp.Result.NextPageCursor

			more = cursor != ""

			if err != nil {
				log.Println("ERROR: Was not able to save lines from api key id =", apiKey.Id, "closed p&l:", err)
				continue
			}
		}

	}

}

func main() {
	dbConfig := repository.PostgresConfigFromEnv()
	fmt.Println(dbConfig.String())
	db, err := sql.Open("pgx", dbConfig.String())

	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserPsql(db)

	bbClient := &ByBitScrapper{
		userRepo: userRepo,
	}

	log.Println("Starting data fetching.")
	err = bbClient.Run()

	if err != nil {
		log.Fatal(err)
	}

}
