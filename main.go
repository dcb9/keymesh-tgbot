package main

import (
	"log"
	"fmt"
	"os"
	"errors"
	"strings"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"encoding/json"
)

var userMap = map[string]string{
	"111": "Howard",
	"222": "Ceo",
	"333": "Bob",
}

var (
	tgToken = os.Getenv("TG_TOKEN")
	tgHookURL = os.Getenv("TG_HOOK_URL")
)

var bot *tgbotapi.BotAPI

func init() {
	log.Printf("start to NewBotAPI with tgToken %s\n", tgToken)
	var err error
	bot, err = tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Fatal("NewBotAPI error: ", err)
	}

	log.Print("set bot.Debug to true")
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(tgHookURL))
	if err != nil {
		log.Fatal("bot.SetWebhook", err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal("bot.GetWebhookInfo", err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("[Telegram callback failed]%s", info.LastErrorMessage)
	}
}

var ErrMessageIsEmpty= errors.New("message could not be empty")

func errResponse(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body: err.Error(),
	}, nil
}

func parseBody(body []byte, update *tgbotapi.Update) error {
	err := json.Unmarshal(body, &update)
	if err != nil {
		return err
	}

	if update.Message == nil {
		return ErrMessageIsEmpty
	}

	return nil
}


func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var update tgbotapi.Update
	err := parseBody([]byte(request.Body), &update)
	if err != nil {
		return errResponse(ErrMessageIsEmpty)
	}

	rawText := update.Message.Text
	log.Printf("[%s] %s", update.Message.From.UserName, rawText)

	if strings.HasPrefix(rawText, "/start") && len(rawText) >= 8 {
		text := update.Message.Text[7:]
		username := userMap[text]
		message := fmt.Sprintf("Hello %s.\nOne more step. 😍😍😍Join group to get 1000 KeyMesh Tokens for free.👇👇👇\n还差一步，进群即可领 1000 KeyMeshToken，更有百万 token 空投活动在 telegram 群不定时赠送：🎁🎁🎁https://t.me/keymesh.", username)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)

		bot.Send(msg)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: "OK",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
