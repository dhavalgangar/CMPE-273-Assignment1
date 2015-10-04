package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

//Global Variables
//var test string = "test"
var stocksList = []string{}
var stocksPercent = []int{}
var stocksBudget = []float64{}
var stocksPrice = []float64{}
var currentStocksPrice = []float64{}
var noOfStocks = []int{}
var budget = 0.00
var tradingID = 11
var uninvestedAmount float64 = 0.00

type StockResponse struct {
	List struct {
		Meta struct {
			Type  string `json:"type"`
			Start int    `json:"start"`
			Count int    `json:"count"`
		} `json:"meta"`
		Resources []struct {
			Resource struct {
				Classname string `json:"classname"`
				Fields    struct {
					Name    string `json:"name"`
					Price   string `json:"price"`
					Symbol  string `json:"symbol"`
					Ts      string `json:"ts"`
					Type    string `json:"type"`
					Utctime string `json:"utctime"`
					Volume  string `json:"volume"`
				} `json:"fields"`
			} `json:"resource"`
		} `json:"resources"`
	} `json:"list"`
}

type StockRequestDetails struct {

	Budget float64
	Stocks string
}

type StockResponseDetails struct {

	TradeID int
	Stocks string
	UnvestedAmount float64
}

type StockPortfolioDetails struct {

	Stocks string
	CurrentMarketValue float64
	UnvestedAmount float64
}

type Computer struct{}

func (t *Computer) BuyingStocks(REQUEST1 *StockRequestDetails, RESPONSE1 *StockResponseDetails) error {

	//fmt.Println(REQUEST1.Stocks)
	temp1 := strings.Split(REQUEST1.Stocks, ",")
  fmt.Println(temp1)

	for i := 0; i < len(temp1); i++ {

		temp2 := strings.Split(temp1[i], ":")
		fmt.Println(temp2)
		stocksList = append(stocksList, temp2[0])
		temp3 := strings.Split(temp2[1], "%")
		j, _ := strconv.Atoi(temp3[0])
		fmt.Println(temp3)
	  stocksPercent = append(stocksPercent, j)
	}
	fmt.Println("end 1st for")
	fmt.Println(stocksList)
	fmt.Println(stocksPercent)
	//get price of each stock
	getStocksPrice()

	var b1 int =1
	var b2 float64 = 1.00

	//compute number of stocks
	for i := 0; i < len(stocksList); i++ {

		b1 = 100 / stocksPercent[i]
		b2 = (REQUEST1.Budget)/float64 (b1)
		stocksBudget = append(stocksBudget, b2)
		b1 = int (b2) / int (stocksPrice[i])
		noOfStocks = append(noOfStocks, b1)

		uninvestedAmount = uninvestedAmount + ( b2 - (stocksPrice[i] * float64(b1)) )
	}

	fmt.Println(stocksBudget)
	fmt.Println(noOfStocks)
	var buffer bytes.Buffer

	for i:= 0; i < len(stocksList); i++ {

		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString("{")
		buffer.WriteString(stocksList[i])
		buffer.WriteString(":")
		b2 = float64(noOfStocks[i])
		buffer.WriteString(strconv.FormatFloat(b2, 'f', 0, 64))
		buffer.WriteString(":$")
		buffer.WriteString(strconv.FormatFloat(stocksPrice[i], 'f', 2, 64))
		buffer.WriteString("}")

	}

	(*RESPONSE1).TradeID = tradingID
	tradingID = tradingID + 11
	(*RESPONSE1).Stocks = buffer.String()
	(*RESPONSE1).UnvestedAmount = uninvestedAmount

	return nil
}

func (t *Computer) CheckPortfolio(REQUEST2 *int, RESPONSE2 *StockPortfolioDetails) error {

		getCurrentStockPrice()

		var buffer bytes.Buffer
		var currentPrice float64 = 0.0
		for i:= 0; i < len(stocksList); i++ {
			if i > 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString("{")
			buffer.WriteString(stocksList[i])
			buffer.WriteString(":")
			b2 := float64(noOfStocks[i])
			buffer.WriteString(strconv.FormatFloat(b2, 'f', 0, 64))
			buffer.WriteString(":")

			if currentStocksPrice[i] > stocksPrice[i] {
				buffer.WriteString("+$")
			} else if currentStocksPrice[i] < stocksPrice[i] {
								buffer.WriteString("-$")
							} else {
								buffer.WriteString("$")
							}
			currentPrice = currentPrice + (currentStocksPrice[i] * noOfStocks[i])
			fmt.Println(currentPrice)
			b2 = float64(currentStocksPrice[i])
			buffer.WriteString(strconv.FormatFloat(b2, 'f', 2, 64))
			buffer.WriteString("}")

		}

		(*RESPONSE2).Stocks = buffer.String()
		(*RESPONSE2).CurrentMarketValue = currentPrice
		(*RESPONSE2).UnvestedAmount = uninvestedAmount

		return nil
}


func getCurrentStockPrice(){

	var s StockResponse
	var urlStocks string
	var buffer bytes.Buffer

	buffer.WriteString("http://finance.yahoo.com/webservice/v1/symbols/")
	for i := 0; i < len(stocksList); i++ {

		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(stocksList[i])

	}
	buffer.WriteString("/quote?format=json")
	urlStocks = buffer.String()

	response, err := http.Get(urlStocks)
	if err != nil {
		fmt.Printf("error occured")
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
			defer response.Body.Close()

			contents, err := ioutil.ReadAll(response.Body)

			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}

			json.Unmarshal([]byte(contents), &s)

			for i := 0; i < s.List.Meta.Count; i++ {
				f, err1 := strconv.ParseFloat(s.List.Resources[i].Resource.Fields.Price, 64)
				currentStocksPrice = append(currentStocksPrice, f)

				if err1 != nil {
					fmt.Printf("%s", err1)
					os.Exit(1)
				}

			}

		}


}

func getStocksPrice(){

	var s StockResponse
	var urlStocks string
	var buffer bytes.Buffer

	buffer.WriteString("http://finance.yahoo.com/webservice/v1/symbols/")
	for i := 0; i < len(stocksList); i++ {

		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(stocksList[i])

	}
	buffer.WriteString("/quote?format=json")
	urlStocks = buffer.String()

	response, err := http.Get(urlStocks)
	if err != nil {
		fmt.Printf("error occured")
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
			defer response.Body.Close()

			contents, err := ioutil.ReadAll(response.Body)

			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}

			json.Unmarshal([]byte(contents), &s)
			//fmt.Println(s)
			//fmt.Println(s.List.Meta.Count)
			for i := 0; i < s.List.Meta.Count; i++ {
				f, err1 := strconv.ParseFloat(s.List.Resources[i].Resource.Fields.Price, 64)
				stocksPrice = append(stocksPrice, f)

				if err1 != nil {
					fmt.Printf("%s", err1)
					os.Exit(1)
				}
				//fmt.Println("stock:", s.List.Resources[i].Resource.Fields.Name, "price:", s.List.Resources[i].Resource.Fields.Price)
			}
			//fmt.Println("====================== ReturnStockValue STOCKS EXITED==================")

		}
}


func main() {

	  cal := new(Computer)
		server := rpc.NewServer()
		server.Register(cal)
		server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
		listener, e := net.Listen("tcp", ":5555")
		if e != nil {
			log.Fatal("listen error:", e)
		}
		for {
			if conn, err := listener.Accept(); err != nil {
				log.Fatal("accept error: " + err.Error())
			} else {
				log.Printf("new connection established\n")
				go server.ServeCodec(jsonrpc.NewServerCodec(conn))
			}
		}
}
