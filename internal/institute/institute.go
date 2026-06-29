package institute

import (
	"fmt"
	"strconv"
	"time"
)

var YEAR int = time.Now().Year()

func ChooseInstitute() string {
	fmt.Println("Выберите институт:")
	for num, institute := range instituteNames {
		fmt.Println(num+1, institute)
	}

	var pickedNum int
	fmt.Scan(&pickedNum)
	instituteName := instituteNames[pickedNum-1]
	fmt.Println("Выбран институт:", instituteName)
	institutePath := instituteMap[instituteName]
	return institutePath

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
	"Институт востоковедения":                                "/Institut_inostrannykh_yazikov/Institut_vostokovedeniya/" + strconv.Itoa(YEAR) + "/",
	"Институт детства":                                       "/Institut_detstva/институт детства " + strconv.Itoa(YEAR) + "/",
	"Институт иностранных языков":                            "/Institut_inostrannykh_yazikov/" + strconv.Itoa(YEAR) + "/",
	"Институт информационных технологий":                     "/Institut_komputernykh_nauk_i_tehn_obr/" + strconv.Itoa(YEAR) + "/",
	"Институт дефектологического образования и реабилитации": "/Korrektsionnoy_pedagogiki/" + strconv.Itoa(YEAR-1) + "/",
	"Институт истории и социальных наук":                     "/Sotsialnykh_nauk/" + strconv.Itoa(YEAR) + "/",
	"Институт музыки, театра и хореографии":                  "/Institut_muziky,_teatra,_horeografii/" + strconv.Itoa(YEAR) + "/",
	"Институт нородов севера":                                "/Institut_narodov_Severa/" + strconv.Itoa(YEAR) + "/",
	"Институт физической культуры и спорта":                  "/Institut_fizicheskoy_kultury/ИФКиС/" + strconv.Itoa(YEAR) + "/",
	"Институт педагогики":                                    "/Institut_pedagogiky_I_psihologii/" + strconv.Itoa(YEAR) + "/Институт педагогики/",
	"Институт психологии":                                    "/Institut_pedagogiky_I_psihologii/" + strconv.Itoa(YEAR) + "/Институт психологии/",
	"Институт русского языка как иностранного":               "/Russkogo_yazyka_kak_inostrannogo/",
	"Институт физики":                                        "/Fiziki/" + strconv.Itoa(YEAR) + "/",
	"Институт философии человека":                            "/Filosofii_cheloveka/" + strconv.Itoa(YEAR) + "/",
	"Институт художественного образования":                   "/Izobrazitelnogo_iskusstva/" + strconv.Itoa(YEAR) + "/",
	"Институт экономики и управления":                        "/Institut_economiky_i_upravleniya/" + strconv.Itoa(YEAR) + "/",
	"Факультет безопасности жизнедеятельности":               "/Bezopasnosti_zhiznedeyatelnosti/",
	"Факультет биологии":                                     "/Biologii/ФБио ГИА " + strconv.Itoa(YEAR) + "/",
	"Факультет географии":                                    "/Geografii/" + strconv.Itoa(YEAR) + "/",
	"Факультет математики":                                   "/Matematiki/" + strconv.Itoa(YEAR) + "/",
	"Факультет химии":                                        "/Khimii/" + strconv.Itoa(YEAR) + "/",
	"Филологический факультет ":                              "/Filologichesky/" + strconv.Itoa(YEAR) + "/",
	"Юридический факультет":                                  "/Yuridichesky/" + strconv.Itoa(YEAR) + "/",
	"Дагестанский филиал":                                    "/Degestanskiy_filial/",
	"Ташкентский филиал":                                     "/Tahskentskiy_filial/",
	"Выборгский филиал":                                      "/Viborgskiy_filial/" + strconv.Itoa(YEAR) + "/",
	"Волховский филиал":                                      "/Volhovskiy_filial/" + strconv.Itoa(YEAR) + "/",
}