package services

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {

	// caminho/x/b/arquivo.mp4
	// split: caminho/x/b/
	// [0] caminho/x/b/
	// [1] arquivo.mp4
	path := strings.Split(objectPath, os.Getenv("localStoragePath")+"/")

	// abre o arquivo
	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	// abre uma coneção para gravar algo
	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)

	//permissão para todos os usuários que pordem ler
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	// copia o arquivo para a conexão writer
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil

}

// carregar todos os caminhos que preciamos para fazer o upload
func (vu *VideoUpload) loadPaths() error {

	// entra no diretorio filepath
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {

		// para cada arquivo do diretorio filepath, tudo que não for um diretório dentro de filepath
		if !info.IsDir() {
			// salva o caminho de cada arquivo em vu.paths
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// retorna um storage client
func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	return client, ctx, nil
}
