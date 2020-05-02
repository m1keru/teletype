package telegram

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/m1keru/teletype/internal/config"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
)

func downloadURL(url string) ([]byte, error) {
	log.Println("downloadURL:", url)
	fileBuffer := new(bytes.Buffer)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Download failed url: %s , error: %+v \n", url, err)
		return []byte{}, err
	}
	defer resp.Body.Close()
	io.Copy(fileBuffer, resp.Body)
	log.Debugf("%s downloaded successfull", url)
	return fileBuffer.Bytes(), nil
}

//Run - run
func Run(cfg *config.Config, wg *sync.WaitGroup, voiceChannel *chan []byte, textChannel *chan string) error {
	defer wg.Done()
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Errorf("Telegram: unable to connect, Error:\n %v", err)
		return err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Telegram.Interval
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Voice != nil {
			log.Println("Audio:", update.Message.Voice.FileID)
			url, err := bot.GetFileDirectURL(update.Message.Voice.FileID)
			if err != nil {
				log.Errorf("Error: %+v", err)
			}
			log.Println(url)
			voice, err := downloadURL(url)
			if err != nil {
				log.Errorf("Error: %+v", err)
			}
			*voiceChannel <- voice
			transcript := <-*textChannel
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, transcript)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
	return nil
}
