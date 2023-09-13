package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main(){

	err := godotenv.Load(".env")
	if err != nil{
		fmt.Println("env 읽기 실패 :", err)
	}

	port := 8080

	r:= mux.NewRouter()

	r.HandleFunc("/data",GetData).Methods("GET")

	
	log.Println("서울 열린데이터 광장 OpenAPI 받아보기")
	log.Println("Server Starting on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",port),r))

}

func GetData(w http.ResponseWriter, r *http.Request){
	queryParams := r.URL.Query()
	strIndex := queryParams.Get("index")
	numIndex, err := strconv.Atoi(strIndex)
	if err != nil {
		fmt.Println("문자열 숫자로 변환 실패 :", err)
	}
	key := os.Getenv("key")
	url := fmt.Sprintf("http://openapi.seoul.go.kr:8088/%s/json/fsiRadioactivityInfo/%d/%d/", key, numIndex, numIndex+4)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("GET 요청 오류 :",err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("body 읽기 오류 :",err)
		return
	}

	var foodInfo fsiRadioactInfo
	if err := json.Unmarshal(body, &foodInfo); err != nil {
		fmt.Println("JSON 파싱 오류:", err)
		http.Error(w, "JSON 파싱 오류", http.StatusInternalServerError)
		return
	}
	
	jsonData, err := json.MarshalIndent(foodInfo, "", "    ") // 빈 문자열과 스페이스 4개로 들여쓰기
	if err != nil {
		fmt.Println("JSON 변환 오류:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	
}

type fsiRadioactInfo struct {
	Info info `json:"fsiRadioactivityInfo"`
}
type info struct {
	ListTotalCount int	`json:"list_total_count"`
	Result status		`json:"RESULT"`
	Row []data			`json:"row"`
}
type status struct {
	Code string			`json:"CODE"`
	Message string		`json:"MESSAGE"`
}

type data struct {
	CheckNo float64				`json:"CHECK_NO"`
	CheckType string			`json:"CHECK_TYPE"`
	PrdtNM string				`json:"PRDT_NM"`
	PrdtOrigin string			`json:"PRDT_ORIGIN"`
	ColYmd string				`json:"COL_YMD"`
	ColPlace string				`json:"COL_PLACE"`
	CheckResultYmd string		`json:"CHECK_RESULT_YMD"`
	CheckCesiumStandard string	`json:"CHECK_CESIUM_STANDARD"`
	CheckCesiumResult string	`json:"CHECK_CESIUM_RESULT"`
	CheckIodinStandard string	`json:"CHECK_IODIN_STANDARD"`
	CheckIodinResult string		`json:"CHECK_IODIN_RESULT"`
	Desision string				`json:"DESISION"`
}
