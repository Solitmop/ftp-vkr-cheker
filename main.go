package main

import (
	"log"
	"os"

	"ftpVkrChecker/internal/ftp_server"
	"ftpVkrChecker/internal/institute"
	"ftpVkrChecker/internal/vkr"
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