package scrapper

import (
	"log"
	"time"

	"github.com/hirokisan/bybit/v2"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store"
)

type ByBitScrapper struct {
	userStore store.User
	apiStore  store.ApiKey
}

func NewByBitScrapper(userStore store.User, apiStore store.ApiKey) *ByBitScrapper {
	return &ByBitScrapper{
		userStore: userStore,
		apiStore:  apiStore,
	}
}

func (bb *ByBitScrapper) Run(userId uint32, start, end time.Time) error {

	apiKeyList, err := bb.apiStore.ListActive(userId)

	if err != nil {
		log.Printf("ERROR: could not get active api keys: %v\n", err)
	}

	log.Printf("INFO: found %d active keys\n", len(apiKeyList))
	bb.ScrapClonedPnL(apiKeyList, start, end)

	return nil
}

func (bb *ByBitScrapper) ScrapClonedPnL(apiKeyList []model.ApiKey, start, end time.Time) {
	yesterdayStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	yesterdayEnd := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, time.UTC)

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

			err = bb.userStore.SaveClosedPnL(apiKey.UserId, apiKey.Id, resp.Result.List)
			cursor = resp.Result.NextPageCursor

			more = cursor != ""

			if err != nil {
				log.Println("ERROR: Was not able to save lines from api key id =", apiKey.Id, "closed p&l:", err)
				continue
			}
		}

	}

}
