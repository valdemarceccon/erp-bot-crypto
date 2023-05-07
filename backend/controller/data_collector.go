package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller/schema"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/scrapper"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store"
)

type DataCollector struct {
	UserStore store.User
	ApiStore  store.ApiKey
}

func NewDataCollector(userStore store.User, apiStore store.ApiKey) *DataCollector {
	return &DataCollector{
		UserStore: userStore,
		ApiStore:  apiStore,
	}
}

func (dc *DataCollector) RunNow(c *fiber.Ctx) error {
	startDate := c.Params("startDate")
	endDate := c.Params("endDate")
	username := c.Params("username")

	if startDate == "" || endDate == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing start date or end date")
	}
	layout := "2006-01-02"

	parsedStartDate, err := time.Parse(layout, startDate)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "start date is incorrect")
	}

	parsedEndDate, err := time.Parse(layout, endDate)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "end date is incorrect")
	}

	if !parsedEndDate.Equal(parsedStartDate) && parsedEndDate.Before(parsedStartDate) {
		return fiber.NewError(fiber.StatusBadRequest, "end date should be equal of after start date")
	}

	var user *model.User
	if username != "" {
		user, err = dc.UserStore.ByUsername(username)

		if err != nil {
			if errors.Is(err, store.ErrUserNotFound) {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("username '%s' does not exists", username))
			}
		}
	}

	err = dc.startCollecting(user, parsedStartDate, parsedEndDate)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("unexpected error: %v", err))
	}

	return c.Status(fiber.StatusCreated).JSON(schema.MessageResponse{
		Message: "Started collecting api",
	})
}

func (dc *DataCollector) startCollecting(user *model.User, start, end time.Time) error {
	var userId uint32 = 0
	if user != nil {
		userId = user.Id
	}
	sc := scrapper.NewByBitScrapper(dc.UserStore, dc.ApiStore)
	go sc.Run(userId, start, end)

	return store.ErrNotImplemented

}
