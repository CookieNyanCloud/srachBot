package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CookieNyanCloud/srachBot/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hako/durafmt"
	"github.com/joho/godotenv"
)

func main() {
	var local bool
	flag.BoolVar(&local, "local", false, "хост")
	flag.Parse()
	if local {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Printf("local env: %v\n", err)
		}
	}
	token := os.Getenv("TG_TOKEN")
	conf, err := strconv.ParseInt(os.Getenv("TG_CONF"), 10, 64)
	if err != nil {
		fmt.Printf("conf: %v\n", err)
	}
	bot, updates := tg.StartSotaBot(token)
	last := time.Now()
	all := make([]*durafmt.Durafmt, 0)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.ID != conf {
			continue
		}
		if strings.Contains(strings.ToLower(update.Message.Text), "опять срач") || strings.Contains(strings.ToLower(update.Message.Text), "hfzbn xgdz") {
			srach := durafmt.Parse(time.Since(last))
			text := fmt.Sprintf("поздравляю, новый срач, c прошлого прошло: %v", srach)
			msg := tgbotapi.NewMessage(conf, text)
			_, _ = bot.Send(msg)
			all = append(all, srach)
			last = time.Now()
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		switch update.Message.Command() {
		case "srach":
			srach := durafmt.Parse(time.Since(last))
			text := fmt.Sprintf("поздравляю, новый срач, c прошлого прошло: %v", srach)
			msg := tgbotapi.NewMessage(conf, text)
			_, _ = bot.Send(msg)
			all = append(all, srach)
			last = time.Now()
		case "last":
			text := fmt.Sprintf("с последнего срача %v", durafmt.Parse(time.Since(last)))
			msg := tgbotapi.NewMessage(conf, text)
			_, _ = bot.Send(msg)
		case "stat":
			var out string
			for _, srach := range all {
				out += fmt.Sprintln(srach)
			}
			text := fmt.Sprintf("статистика:\n%v", out)
			msg := tgbotapi.NewMessage(conf, text)
			_, _ = bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(conf, "/srach\n/last\n/stat")
			_, _ = bot.Send(msg)
		}
	}
}
