package speech

import (
	"context"
	"github.com/m1keru/teletype/internal/config"
	"log"
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

	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	var transcript string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			//fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			log.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			transcript = transcript + alt.Transcript
		}
	}
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
		recognise(client, msg, textChannel)
	}
	return nil
}
