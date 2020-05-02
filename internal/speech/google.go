package speech

import (
	"context"
	"github.com/m1keru/teletype/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func recognise(client *speech.Client, data []byte, textChannel *chan string) error {
	ctx := context.Background()

	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_OGG_OPUS,
			SampleRateHertz: 48000,
			LanguageCode:    "ru-Ru",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}
	log.Debug("Executing request to Google Cloud API")
	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		*textChannel <- "Хер его знает че он там пизданул.\n"
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		*textChannel <- "Хер его знает че он там пизданул.\n"
		return err
	}

	log.Debug("Google Cloud API Request executed")
	log.Debug("Len Results:", len(resp.Results))
	var transcript string
	if len(resp.Results) == 0 {
		*textChannel <- "Хер его знает че он там пизданул.\n"
	}

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			log.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			transcript = transcript + alt.Transcript
		}
	}
	transcript = transcript + "\n"
	*textChannel <- transcript
	return nil
}

func configEnv(cfg *config.Config) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", cfg.Speech.SecretFile)
}

//Run - Run
func Run(cfg *config.Config, wg *sync.WaitGroup, voiceChannel *chan []byte, textChannel *chan string) error {
	configEnv(cfg)
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	for msg := range *voiceChannel {
		log.Debug("Recieved new msg from Telegram go-routine")
		err := recognise(client, msg, textChannel)
		if err != nil {
			log.Errorf("[RECOGNISE][FATAL]: %+v", err)
		}
	}
	return nil
}
