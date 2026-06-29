package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
	"github.com/tealeg/xlsx/v3"
)

const (
	YEAR = "2026"
)

func main() {
	err := connect()
	if err != nil {
		log.Println(err)
	}

	fmt.Scanf("\n")
	fmt.Scanf("\n")
}

func connect() error {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return err
	}

	// Подключение к FTP хосту
	ftpClient, err := ftp.Dial(os.Getenv("host"))
	if err != nil {
		log.Println("Ошибка подключения к FTP:", err)
		return err
	}
	defer ftpClient.Quit()

	// Аутентификация
	err = ftpClient.Login(os.Getenv("login"), os.Getenv("password"))
	if err != nil {
		log.Println("Ошибка авторизации:", err)
		return err
	}
	defer ftpClient.Logout()

	err = checkInstitute(ftpClient)
	if err != nil {
		return err
	}
	
	fmt.Println("Подключение закрыто")
	return nil
}

func checkInstitute(ftpClient *ftp.ServerConn) error {
	institutePath := chooseInstitute()

	err := ftpClient.ChangeDir(institutePath)
	if err != nil {
		log.Println("Ошибка перехода в директорию:", err)
		return err
	}
	checkGroups(ftpClient)
	
	return nil
}

func chooseInstitute() string {
	fmt.Println("Выберите институт:")
	for num, institute := range instituteNames {
		fmt.Println(num+1, institute)
	}

	var instituteNum int
	fmt.Scan(&instituteNum)
	instituteName := instituteNames[instituteNum-1]
	fmt.Println("Выбран институт:", instituteName)
	institutePath := instituteMap[instituteName]
	return institutePath

}

func checkVkrs(ftpClient *ftp.ServerConn) error {
	filesRaw, err := ftpClient.List(".")
	if err != nil {
		log.Println("Ошибка чтения: ", err)
		return err
	}
	path, err := ftpClient.CurrentDir()
	if err != nil {
		log.Println("Ошибка чтения: ", err)
		return err
	}
	files := []*ftp.Entry{}
	for _, file := range filesRaw {
		if file.Name != "." && file.Name != ".." {
			files = append(files, file)
		}
	}

	fileVkrs := getFileVkrs(files)
	excelVkrs, err := getExcelVkrs(ftpClient, files)
	if err != nil {
		return err
	}
	sort.Strings(fileVkrs)
	sort.Strings(excelVkrs)

	fmt.Println(path + "/")
	if len(fileVkrs) != len(excelVkrs) {
		fmt.Println("Количество файлов не совпадает с количеством Excel-файлов")
		fmt.Println("Количество файлов:", len(fileVkrs))
		fmt.Println("Количество Excel-файлов:", len(excelVkrs))
		return errors.New("количество файлов не совпадает с количеством Excel-файлов")
	} else if !slices.Equal(fileVkrs, excelVkrs) {
		fmt.Println("Названия файлов не совпадают")
		//fmt.Println(fileVkrs)
		//fmt.Println(excelVkrs)
		fmt.Println("Нет файлов работ:")
		for _, x := range excelVkrs {
			if !slices.Contains(fileVkrs, x) {
				fmt.Println(x)
			}
		}
		fmt.Println("Нет работ в Excel:")
		for _, x := range fileVkrs {
			if !slices.Contains(excelVkrs, x) {
				fmt.Println(x)
			}
		}
		return errors.New("названия файлов не совпадают")
	}

	fmt.Println("Количество файлов:", len(fileVkrs))
	return nil
}

func getFileVkrs(files []*ftp.Entry) (names []string) {
	names = []string{}
	for _, file := range files {
		if !strings.Contains(file.Name, ".xlsx") || !strings.Contains(file.Name, ".xls") {
			names = append(names, file.Name)
		}
	}
	return
}

func getExcelVkrs(ftpClient *ftp.ServerConn, files []*ftp.Entry) ([]string, error) {
	excel := []*ftp.Entry{}
	for _, file := range files {
		if strings.Contains(file.Name, ".xlsx") || strings.Contains(file.Name, ".xls") {
			excel = append(excel, file)
		}
	}

	if len(excel) == 0 {
		fmt.Println("Не найдено Excel-файлов")
		return nil, errors.New("не найдено Excel-файлов")
	} else if len(excel) != 1 {
		fmt.Println("Найдено несколько Excel-файлов")
		return nil, errors.New("найдено несколько Excel-файлов")
	}

	fileReader, err := ftpClient.Retr(excel[0].Name)
	if err != nil {
		log.Println("Ошибка открытия файла:", err)
		return nil, err
	}
	defer fileReader.Close()

	data, err := io.ReadAll(fileReader)
	if err != nil {
		log.Println("Ошибка чтения данных:", err)
		return nil, err
	}
	readerAt := bytes.NewReader(data)
	size := int64(len(data))

	xlFile, err := xlsx.OpenReaderAt(readerAt, size)
	if err != nil {
		log.Println("Ошибка чтения Excel:", err)
		return nil, err
	}
	names, err := getColumnByName(xlFile, "Имя pdf-файла ВКР")
	if err != nil {
		log.Println("Ошибка чтения столбца:", err)
		return nil, err
	}
	//fmt.Println("excel",names)
	return names, nil
}

// getColumnByName возвращает все значения столбца по его названию (из первой строки)
func getColumnByName(xlFile *xlsx.File, columnName string) ([]string, error) {
	var columnValues []string
	columnIndex := -1

	// Обрабатываем только первый лист
	if len(xlFile.Sheets) == 0 {
		return nil, fmt.Errorf("нет листов в файле")
	}

	sheet := xlFile.Sheets[0]

	// Ищем нужный столбец в первой строке
	headerRow, err := sheet.Row(0)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения заголовков: %v", err)
	}

	numHadders := 0
	// Перебираем ячейки первой строки для поиска столбца
	found := false
	for col := 0; col < sheet.MaxCol; col++ {
		cell := headerRow.GetCell(col)
		if cell == nil {
			continue
		}

		header, err := cell.FormattedValue()
		if err != nil {
			continue
		}

		numHadders++
		if header == columnName {
			columnIndex = col
			found = true
			break
		}
	}

	if len(tableStruct) > numHadders {
		return nil, fmt.Errorf("названия столбцов не совпадают %v", numHadders)
	}

	if !found {
		return nil, fmt.Errorf("столбец '%s' не найден", columnName)
	}

	// Собираем значения столбца (начиная со второй строки)
	for i := 1; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			columnValues = append(columnValues, "") // Пустая строка
			continue
		}

		cell := row.GetCell(columnIndex)
		if cell == nil {
			columnValues = append(columnValues, "")
			continue
		}

		value, err := cell.FormattedValue()
		if err != nil {
			columnValues = append(columnValues, "")
		} else {
			columnValues = append(columnValues, value)
		}
	}

	return columnValues, nil
}

func checkGroups(ftpClient *ftp.ServerConn) error {
	files, err := ftpClient.List(".")
	if err != nil {
		log.Println("Ошибка чтения: ", err)
		return err
	}

	groups := []string{}
	for _, file := range files {
		if file.Name != "." && file.Name != ".." {
			groups = append(groups, file.Name)
		}
	}
	if len(groups) == 0 {
		log.Println("Папка пуста")
		return errors.New("папка пуста")
	}
	for _, group := range groups {
		err = ftpClient.ChangeDir(group)
		if err != nil {
			log.Printf("Ошибка перехода в директорию группы %s: %s\n", group, err)
			return err
		}
		checkVkrs(ftpClient)
		err = ftpClient.ChangeDir("..")
		if err != nil {
			log.Println("Ошибка возвращения в директорию подразделения:", err)
			return err
		}
	}
	return nil
}

var instituteNames = []string{
	"Институт востоковедения",
	"Институт детства",
	"Институт иностранных языков",
	"Институт информационных технологий",
	"Институт дефектологического образования и реабилитации",
	"Институт истории и социальных наук",
	"Институт музыки, театра и хореографии",
	"Институт нородов севера",
	"Институт физической культуры и спорта",
	"Институт педагогики",
	"Институт психологии",
	"Институт русского языка как иностранного",
	"Институт физики",
	"Институт философии человека",
	"Институт художественного образования",
	"Институт экономики и управления",
	"Факультет безопасности жизнедеятельности",
	"Факультет биологии",
	"Факультет географии",
	"Факультет математики",
	"Факультет химии",
	"Филологический факультет ",
	"Юридический факультет",
	"Дагестанский филиал",
	"Ташкентский филиал",
	"Выборгский филиал",
	"Волховский филиал",
}

var instituteMap = map[string]string{
	"Институт востоковедения":                                "/Institut_inostrannykh_yazikov/Institut_vostokovedeniya/" + YEAR + "/",
	"Институт детства":                                       "/Institut_detstva/институт детства " + YEAR + "/",
	"Институт иностранных языков":                            "/Institut_inostrannykh_yazikov/" + YEAR + "/",
	"Институт информационных технологий":                     "/Institut_komputernykh_nauk_i_tehn_obr/" + YEAR + "/",
	"Институт дефектологического образования и реабилитации": "/Korrektsionnoy_pedagogiki/" + YEAR + "/",
	"Институт истории и социальных наук":                     "/Sotsialnykh_nauk/" + YEAR + "/",
	"Институт музыки, театра и хореографии":                  "/Institut_muziky,_teatra,_horeografii/" + YEAR + "/",
	"Институт нородов севера":                                "/Institut_narodov_Severa/" + YEAR + "/",
	"Институт физической культуры и спорта":                  "/Institut_fizicheskoy_kultury/ИФКиС/" + YEAR + "/",
	"Институт педагогики":                                    "/Institut_pedagogiky_I_psihologii/" + YEAR + "/Институт педагогики/",
	"Институт психологии":                                    "/Institut_pedagogiky_I_psihologii/" + YEAR + "/Институт психологии/",
	"Институт русского языка как иностранного":               "/Russkogo_yazyka_kak_inostrannogo/",
	"Институт физики":                                        "/Fiziki/" + YEAR + "/",
	"Институт философии человека":                            "/Filosofii_cheloveka/" + YEAR + "/",
	"Институт художественного образования":                   "/Izobrazitelnogo_iskusstva/" + YEAR + "/",
	"Институт экономики и управления":                        "/Institut_economiky_i_upravleniya/" + YEAR + "/",
	"Факультет безопасности жизнедеятельности":               "/Bezopasnosti_zhiznedeyatelnosti/",
	"Факультет биологии":                                     "/Biologii/ФБио ГИА " + YEAR + "/",
	"Факультет географии":                                    "/Geografii/" + YEAR + "/",
	"Факультет математики":                                   "/Matematiki/" + YEAR + "/",
	"Факультет химии":                                        "/Khimii/" + YEAR + "/",
	"Филологический факультет ":                              "/Filologichesky/" + YEAR + "/",
	"Юридический факультет":                                  "/Yuridichesky/" + YEAR + "/",
	"Дагестанский филиал":                                    "/Degestanskiy_filial/",
	"Ташкентский филиал":                                     "/Tahskentskiy_filial/",
	"Выборгский филиал":                                      "/Viborgskiy_filial/" + YEAR + "/",
	"Волховский филиал":                                      "/Volhovskiy_filial/" + YEAR + "/",
}

var tableStruct = []string{
	"Тема ВКР",
	"ФИО студента",
	"Ученая степень, должность руководителя",
	"ФИО руководителя",
	"ФИО консультанта, ученая степень, должность",
	"ФИО рецензента, ученая степень, должность",
	"Стандарт, по которому указаны УГСН и направления подготовки ВКР (ФГОС ВО (3+), ФГОС ВО (3++))",
	"УГСН (код - название УГСН)",
	"Направления подготовки (код - название направления подготовки)",
	"Профиль",
	"Магистерская программа",
	"Уровень образования",
	"Университет",
	"Факультет",
	"Институт",
	"Кафедра",
	"Год выпуска",
	"Форма обучения",
	"Аннотация",
	"Объем ВКР (кол-во стр.)",
	"Имя pdf-файла ВКР",
}
