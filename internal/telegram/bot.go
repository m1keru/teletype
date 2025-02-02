package telegram

import (
	"bytes"
	"encoding/binary"
	//"encoding/hex"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/m1keru/teletype/internal/config"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
)

//DownloadURL - DownloadURL
func DownloadURL(url string) ([]byte, error) {
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

//DetectRate - DetectRate
func DetectRate(file *[]byte) (uint32, error) {
	if len(*file) < 40 {
		return 0, errors.New("No data in downloaded from Telegram audio")
	}
	rate := (*file)[40:42]
	rate = append(rate, byte(0), byte(0))
	data := binary.LittleEndian.Uint32(rate)
	log.Printf("BITRATE: %d", data)
	return data, nil
}

//Run - run
func Run(cfg *config.Config, wg *sync.WaitGroup, voiceChannel *chan []byte, textChannel *chan string) error {
	defer wg.Done()
	if cfg.Daemon.Debug {
		log.SetLevel(log.DebugLevel)
	}
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Errorf("Telegram: unable to connect, Error:\n %v", err)
		return err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Telegram.Interval
	updates, err := bot.GetUpdatesChan(u)
	serviceEnabled := true

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Voice != nil && serviceEnabled {
			log.Println("Audio:", update.Message.Voice.FileID)
			url, err := bot.GetFileDirectURL(update.Message.Voice.FileID)
			if err != nil {
				log.Errorf("Error: %+v", err)
			}
			log.Println(url)
			voice, err := DownloadURL(url)
			if err != nil {
				log.Errorf("Error: %+v", err)
			}
			rate, _ := DetectRate(&voice)
			log.Debugf("HZ: %d", rate)
			*voiceChannel <- voice
			transcript := <-*textChannel
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, transcript)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "Включть распознавание: /start_voice Выключить: /stop_voice"
			case "start_voice":
				msg.Text = "включился"
				serviceEnabled = true
				log.Println("Service Enabled from Chat Command")
			case "stop_voice":
				msg.Text = "выключился"
				serviceEnabled = false
				log.Println("Service Disabled from Chat Command")
			default:
				msg.Text = "Че тупой? Есть только /start_voice и /stop_voice. Окай?"
			}
			bot.Send(msg)
		}
	}
	return nil
}
