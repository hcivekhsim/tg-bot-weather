package handler

import (
	"fmt"
	"log"
	"math"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hcivekhsim/tg-bot-weather/clients/openweather"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
	}
}
func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.handleUpdate(update)
	}
}
func (h *Handler) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	coordinates, err := h.owClient.Coordinates(update.Message.Text)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смог получить координаты :(")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смог получить погоду :(")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("❄️ Температура в %s: %d°C", update.Message.Text, int(math.Round(weather.Temp))),
	)

	if update.Message.IsCommand() && update.Message.Command() == "time" {
		now := time.Now()
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("Сейчас: %s", now.Format("10:00")),
		)
		h.bot.Send(msg)
		return
	}
	h.bot.Send(msg)
	msg.ReplyToMessageID = update.Message.MessageID
}
