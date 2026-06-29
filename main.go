package main

import (
	"log"
	"os"

	"github.com/solitmop/ftp-vkr-checker/internal/ftp_server"
	"github.com/solitmop/ftp-vkr-checker/internal/institute"
	"github.com/solitmop/ftp-vkr-checker/internal/vkr"
)

func main() {
	// Соединение с ftp сервером
	ftpClient, err := ftp_server.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Подключен к серверу")
	
	path := institute.ChooseInstitute()

	err = vkr.CheckPath(path, ftpClient)
	if err != nil {
		log.Println(err)
	}
	
	ftp_server.Disconnect(ftpClient)
	log.Println("Проверка завершена")
}