package services

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/wallacemachado/microsservico-encoder-go/application/repositories"
	"github.com/wallacemachado/microsservico-encoder-go/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {

	// caso queira interromper o processo
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("video %v has been stored", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {

	// od.ModePerm dá permissão na pasta
	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"

	//mp4 foi instalado pelo dockerfile
	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

// encode para para outro formato
func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	// pegando o arquivo
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID+".frag")
	// pega o arquivo e divide em menores para não fazer o upload de uma vez
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	// output onde ele vai salvar
	cmdArgs = append(cmdArgs, "-o")
	// local onde ele vai salvar
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID)

	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	// pegar a saída do cmd
	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	// imprimir mensagem na tela
	printOutput(output)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n", string(out))
	}
}
