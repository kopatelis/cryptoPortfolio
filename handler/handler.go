package handler

import (
	"cryptoPortfolio/clients/gecko"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler struct {
	bot         *tgbotapi.BotAPI
	geckoClient *gecko.GeckoClient
}

func New(bot *tgbotapi.BotAPI, geckoClient *gecko.GeckoClient) *Handler {
	return &Handler{
		bot:         bot,
		geckoClient: geckoClient,
	}
}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := h.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		h.handleUpdate(update)
	}
}

func (h *Handler) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	res, err := h.geckoClient.Price(update.Message.Text)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли получить цену")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v: %v", res.Name, res.Price))
	msg.ReplyToMessageID = update.Message.MessageID

	h.bot.Send(msg)
}
