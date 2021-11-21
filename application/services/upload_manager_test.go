package services_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/wallacemachado/microsservico-encoder-go/application/services"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("fullcycletest")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "fullcycletest"
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(10, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "upload completed")

	err = videoService.Finish()
	require.Nil(t, err)

}
