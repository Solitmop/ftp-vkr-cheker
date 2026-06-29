package vkr

import (
	"errors"
	"ftpVkrChecker/internal/excel"
	"log"
	"slices"
	"sort"
	"strings"
	"fmt"

	"github.com/jlaffaye/ftp"
)

func CheckPath(path string, ftpClient *ftp.ServerConn) error {
	err := ftpClient.ChangeDir(path)
	if err != nil {
		log.Println("Ошибка перехода в директорию:", err)
		return err
	}

	return checkFolder(ftpClient)
}

func checkFolder(ftpClient *ftp.ServerConn) error {
	path, err := ftpClient.CurrentDir()
	if err != nil {
		log.Println("Ошибка чтения: ", err)
		return err
	}
	fmt.Println(path)
	filesRaw, err := ftpClient.List(path)
	if err != nil {
		log.Println("Ошибка чтения: ", err)
		return err
	}

	pdfVkrs := []*ftp.Entry{}
	tables := []*ftp.Entry{}
	folders := []*ftp.Entry{}
	for _, file := range filesRaw {
		if strings.Contains(file.Name, ".xlsx") || strings.Contains(file.Name, ".xls") {
			tables = append(tables, file)
		} else if strings.Contains(file.Name, ".pdf") {
			pdfVkrs = append(pdfVkrs, file)
		} else if !(file.Name == ".") && !(file.Name == "..") {
			folders = append(folders, file)
		}
	}

	//log.Println(len(pdfVkrs), len(tables), len(folders))

	sortFiles(pdfVkrs)
	//sortFiles(tables)
	sortFiles(folders)

	errorsList := []error{}

	if len(folders) != 0 {
		for _, folder := range(folders) {
			//log.Println("New", "./" + folder.Name)
			err = ftpClient.ChangeDir("./" + folder.Name)
			if err != nil {
				log.Printf("Ошибка перехода в директорию %s: %s\n", folder.Name, err)
				errorsList = append(errorsList, err)
				continue
			}
			err = checkFolder(ftpClient)
			if err != nil {
				errorsList = append(errorsList, err)
			}
			err = ftpClient.ChangeDirToParent()
			if err != nil {
				log.Printf("Ошибка возвращения в директорию %s: %s\n", path, err)
				return err
			}
		}
	}

	if len(pdfVkrs) != 0 {
		if len(tables) == 0 {
			err = errors.New("Нет таблицы")
			log.Println(err)
			return err
		}
		err = checkVkrs(filesToNames(pdfVkrs), tables, ftpClient)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if len(folders) >0 {
		fmt.Printf("Проверенно групп %v/%v\n", len(folders)-len(errorsList), len(folders))
	}

	return nil
}

func checkVkrs(pdfVkrs []string, tables []*ftp.Entry, ftpClient *ftp.ServerConn) error {
	if len(tables) != 1 {
		return errors.New("Каталог содержит больше 1 таблицы")
	}
	table := tables[0]
	excelVkrs, err := excel.GetVkrs(table, ftpClient)
	if err != nil {
		return err
	}
	sort.Strings(excelVkrs)

	if slices.Equal(pdfVkrs, excelVkrs){
		fmt.Println("Файлы совпадают. Количество файлов:", len(pdfVkrs))
		fmt.Println()
		return nil
	}

	fmt.Println("Количество файлов:", len(pdfVkrs))
	fmt.Println("Количество Excel-файлов:", len(excelVkrs))

	fmt.Println("Нет pdf файлов работ:")
	for _, x := range excelVkrs {
		if !slices.Contains(pdfVkrs, x) {
			fmt.Println(x)
		}
	}
	fmt.Println("Нет работ в Excel:")
	for _, x := range pdfVkrs {
		if !slices.Contains(excelVkrs, x) {
			fmt.Println(x)
		}
	}
	return errors.New("Количество файлов не совпадает с количеством Excel записей\n")

}



func sortFiles(files []*ftp.Entry) {
	sort.Slice(files, func(i, j int) bool {
    	return files[i].Name < files[j].Name
	})
}

func filesToNames(files []*ftp.Entry) []string {
	names := []string{}
	for _, file := range(files){
		names = append(names, file.Name)
	}
	return names
}

// TODO: x/y без ошибок