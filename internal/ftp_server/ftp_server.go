package ftp_server

import (
	"log"
	"os"

	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
)


func Connect() (*ftp.ServerConn,error) {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки .env файла", err)
		return nil, err
	}

	// Подключение к FTP хосту
	ftpClient, err := ftp.Dial(os.Getenv("host"))
	if err != nil {
		log.Println("Ошибка подключения к FTP:", err)
		return nil, err
	}

	// Аутентификация
	err = ftpClient.Login(os.Getenv("login"), os.Getenv("password"))
	if err != nil {
		log.Println("Ошибка авторизации:", err)
		return nil, err
	}
	
	return ftpClient, nil
}

func Disconnect(ftpClient *ftp.ServerConn) {
	ftpClient.Logout()
	ftpClient.Quit()
	log.Println("Подключение закрыто")
}