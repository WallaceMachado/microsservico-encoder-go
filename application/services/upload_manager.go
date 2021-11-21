package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {

	//canal que armazena um numero inteiro
	// runtime.NumCPU() => retorna o numero de cpu da máquina, se tivermos numa maquina octocore, ele retorna um buth de 8
	// e o canal vai poder armazenar até 8 dados
	in := make(chan int, runtime.NumCPU()) // qual o arquivo baseado na posicao do slice Paths
	returnChannel := make(chan string)     // vai avisar quando cada upload terminar, pode ser um erro ou um ok

	// carrega todos os caminhos dos arquivos
	err := vu.loadPaths()
	if err != nil {
		return err
	}

	// gera um storage client
	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	// loop que cria a quantidade de processos (goroutine) de upload que serão gerados em concorrência/simultaneos
	// essas goroutines ficam rodando em bacground de forma assincrona
	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	// goroutine/processo assincrono que vai ler vu.Paths e enviar cada path no canal de entrada in
	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
			// o canal só armazendo um dado de cada vez, caso no retorno do loop o canal ainda esetja cheio
			// o processo fica parado esperadno canal ser esvaziado
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil

}

// executa/faz o upload
func (vu *VideoUpload) uploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {

	// le/esvazia o canal in
	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)

		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("error during the upload: %v. Error: %v", vu.Paths[x], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- "upload completed"

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
