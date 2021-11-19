package services_test

import (
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/wallacemachado/microsservico-encoder-go/application/repositories"
	"github.com/wallacemachado/microsservico-encoder-go/application/services"
	"github.com/wallacemachado/microsservico-encoder-go/domain"
	"github.com/wallacemachado/microsservico-encoder-go/framework/database"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// cnpj

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}

	return video, repo
}

func TestVideoServiceDownload(t *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("fullcycletest")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

}
