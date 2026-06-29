package excel

import (
	"bytes"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/jlaffaye/ftp"
	"github.com/extrame/xls"
	"github.com/tealeg/xlsx/v3"
)

func GetVkrs(table *ftp.Entry, ftpClient *ftp.ServerConn) ([]string, error) {
	fileReader, err := ftpClient.Retr(table.Name)
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

	if strings.Contains(table.Name, ".xlsx") {
		xlFile, err := xlsx.OpenReaderAt(readerAt, size)
		if err != nil {
			log.Println("Ошибка чтения Excel:", err)
			return nil, err
		}
		names, err := getColumnByNameXLSX(xlFile, "Имя pdf-файла ВКР")
		if err != nil {
			log.Println("Ошибка чтения столбца:", err)
			return nil, err
		}
		//fmt.Println("excel",names)
		return names, nil
	} else if strings.Contains(table.Name, ".xls") {
		// Открываем xls файл из буфера
		workbook, err := xls.OpenReader(bytes.NewReader(data), "utf-8")
		if err != nil {
			log.Fatal("Ошибка открытия xls файла:", err)
		}		

		names, err := getColumnByNameXLS(workbook, "Имя pdf-файла ВКР")
		if err != nil {
			log.Println("Ошибка чтения столбца:", err)
			return nil, err
		}

		return names, nil
	}
	return nil, errors.New("Неверное расширение файла таблицы")
}

// getColumnByNameXLSX возвращает все значения столбца по его названию (из первой строки)
func getColumnByNameXLSX(xlsxFile *xlsx.File, columnName string) ([]string, error) {
	var columnValues []string
	columnIndex := -1

	// Обрабатываем только первый лист
	if len(xlsxFile.Sheets) == 0 {
		return nil, errors.New("Нет листов в файле")
	}

	sheet := xlsxFile.Sheets[0]

	// Ищем нужный столбец в первой строке
	headerRow, err := sheet.Row(0)
	if err != nil {
		return nil, err
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
		return nil, errors.New("Названия столбцов не совпадают " + strconv.Itoa(numHadders))
	}

	if !found {
		return nil, errors.New("Не найден столбец " + columnName)
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


func getColumnByNameXLS(workbook *xls.WorkBook, columnName string) ([]string, error) {
	sheet := workbook.GetSheet(0)
	if sheet == nil {
		log.Fatal("Лист не найден")
	}

	var columnValues []string
	columnIndex := -1
	numHadders := 0

	headerRow := sheet.Row(0)
	
	found := false
	for col := 0; col <= headerRow.LastCol(); col++ {
		cell := headerRow.Col(col)
		if cell == "" {
			continue
		}
		numHadders++
		if cell == columnName {
			columnIndex = col
			found = true
			break
		}
	}

	if len(tableStruct) > numHadders {
		return nil, errors.New("Названия столбцов не совпадают " + strconv.Itoa(numHadders))
	}

	if !found {
		return nil, errors.New("Не найден столбец " + columnName)
	}

	for i := 1; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		
		cell := row.Col(columnIndex)
		columnValues = append(columnValues, cell)
	}

	return columnValues, nil
}

// TODO:

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
