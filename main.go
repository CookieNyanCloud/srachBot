package main

import (
	"bufio"
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
	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("open file: %v\n", err)
	}
	defer file.Close()
	token := os.Getenv("TG_TOKEN")
	conf, err := strconv.ParseInt(os.Getenv("TG_CONF"), 10, 64)
	if err != nil {
		fmt.Printf("conf: %v\n", err)
	}
	bot, updates := tg.StartSotaBot(token)
	last := time.Now().Add(-50 * time.Hour)
	all := make([]*durafmt.Durafmt, 0)
	err = readDateFile(all)
	if err != nil {
		fmt.Printf("read: %v\n", err)
	}
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
			err := saveDate(file, srach)
			if err != nil {
				fmt.Printf("write: %v\n", err)
			}
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
			err := saveDate(file, srach)
			if err != nil {
				fmt.Printf("write: %v\n", err)
			}
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
		}
	}
}

func saveDate(file *os.File, date *durafmt.Durafmt) error {
	_, err := file.WriteString(date.String())
	return err
}

func readDateFile(date []*durafmt.Durafmt) error {
	file, err := os.Open("data.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, err := durafmt.ParseString(scanner.Text())
		if err != nil {
			return err
		}
		date = append(date, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
