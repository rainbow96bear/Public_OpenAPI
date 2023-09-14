package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
)

func main(){

	err := godotenv.Load(".env")
	if err != nil{
		fmt.Println("env 읽기 실패 :", err)
	}

	port := 8080
	r := pat.New()

	r.Get("/data",GetData)

	http.Handle("/",r)
	log.Println("포트번호", port,"서버 생성")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",port),nil))
}

func GetData(w http.ResponseWriter, r *http.Request){
	queryParams := r.URL.Query()
	strDate := queryParams.Get("date")
	numDate,err := strconv.Atoi(strDate)
	if err != nil {
		fmt.Println("문자열 변환 실패", err)
	}


	id := os.Getenv("id")
	passwd := os.Getenv("passwd")

	url := fmt.Sprintf("http://www.garak.co.kr/publicdata/dataJsonOpen.do?id=%v&passwd=%v&dataid=data35&pagesize=10&pageidx=1&portal.templet=false&s_date=%d&s_pummok=", id, passwd, numDate)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("GET 요청 오류", err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Body 읽기 실패", err)
		return
	}

	newData:= strings.ReplaceAll(string(body), "{ LIST_COUNT", "{ \"LIST_COUNT\"")
	newData = strings.ReplaceAll(newData, ",:", ",\"items\":")
	newData = strings.ReplaceAll(newData, "],[", ",")

	var list List

	if err := json.Unmarshal([]byte(newData), &list); err != nil{
		fmt.Println("Json Parsing Error")
	}

	jsonData, err := json.MarshalIndent(list, "", "    ") // 빈 문자열과 스페이스 4개로 들여쓰기
	if err != nil {
		fmt.Println("JSON 변환 오류:", err)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

type List struct {
	ListCount int `json:"LIST_COUNT"`
	Items []Item `json:"items"`
}

type Item struct {
	Rowno int	`json:"ROWNO"`
	ItemName string	`json:"ITM_NM"`
	FavitmCdnm string	`json:"FAVITM_CDNM"`
	UnitName string	`json:"UNIT_NM"`
	GrdName string	`json:"GRD_NM"`
	SaleUnitPrice int	`json:"SALE_UNIT_PRICE"`
	ZipNoName string	`json:"ZIP_NO_NM"`
	Kind string	`json:"KIND"`
}