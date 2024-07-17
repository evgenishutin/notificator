package notificator

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/enescakir/emoji"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
)

type NotificatorInterface interface {
	SendMessage(info map[string]interface{})
}

type tg struct {
	TelegramToken string
	ChatID        int64
	Title         string
	Service       notify.Notifier
}

type NotifyService struct {
	Telegram tg
}

func New(token string, chatID int64, serviceName string) (NotificatorInterface, error) {
	telegramService, err := telegram.New(token)
	if err != nil {
		return &NotifyService{}, err
	}

	telegramService.AddReceivers(chatID)
	noti := notify.New()
	noti.UseServices(telegramService)

	tg := tg{
		TelegramToken: token,
		ChatID:        chatID,
		Title:         serviceName,
		Service:       noti,
	}

	return &NotifyService{
		Telegram: tg,
	}, nil
}

func (noti *NotifyService) SendMessage(info map[string]interface{}) {
	timestamp := time.Now()
	var message string
	subject := fmt.Sprintf("%v \n<b>Service : %s</b>", emoji.RedCircle, noti.Telegram.Title)

	for key, value := range info {
		message += fmt.Sprintf("\n<b>%s</b> : %s", key, value)
	}

	message += fmt.Sprintf("\n<b>%s</b> : %s", "time", timestamp.Format(time.RFC3339))

	err := noti.Telegram.Service.Send(
		context.Background(),
		subject,
		message,
	)
	if err != nil {
		log.Println("send notification error:", err.Error())
	}
}
