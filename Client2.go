package main

import (
	"fmt"
	"strings"
	"os"
	"strconv"
	"log"
	"net"
	"net/rpc/jsonrpc"
)

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

func main() {

	fmt.Println("Enter your option -")
	fmt.Println("1. Buy Stock \n2. View Portfolio")
	var option int
	fmt.Scanln(&option)

	if option == 1 {

		fmt.Print("Enter Stock details in format - GOOG:50%,YHOO:50%,... :")
		var stockString string
		fmt.Scanln(&stockString)

		/*
			write logic to handle percent of stocks = 100
			if less or more than 100 then exit from application
		*/

		////////////////////////////////////////////////////////
		var totalPercent int = 0
		temp1 := strings.Split(stockString, ",")
	  
		for i := 0; i < len(temp1); i++ {

			temp2 := strings.Split(temp1[i], ":")
			temp3 := strings.Split(temp2[1], "%")
			j, _ := strconv.Atoi(temp3[0])
			totalPercent = totalPercent + j
		}

		if totalPercent != 100 {

				fmt.Println("Inputs invalid, toaatal percent should be equal to 100.")
				os.Exit(1)
		}

		///////////////////////////////////////////////////////


		fmt.Print("Enter Total Budget :")
		var budget float64
		fmt.Scanf("%f", &budget)

		REQUEST1 := &StockRequestDetails{ budget, stockString }

		var RESPONSE1 StockResponseDetails

		client, err := net.Dial("tcp", "127.0.0.1:5555")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		c := jsonrpc.NewClient(client)
		err = c.Call("Computer.BuyingStocks", REQUEST1, &RESPONSE1 )
		if err != nil {
			log.Fatal("arith error:", err)
		}

		fmt.Println("Stocks Purchased, below are the details :-")
		fmt.Print("Trade ID :")
		fmt.Println(RESPONSE1.TradeID)
		fmt.Print("Stocks :")
		fmt.Println(RESPONSE1.Stocks)
		fmt.Print("Unvested Amount :")
		fmt.Println(RESPONSE1.UnvestedAmount)

	} else if option == 2{

		fmt.Print("Enter your TradeID: ")
		var tradeiD int
		fmt.Scanln(&tradeiD)

		REQUEST2 := &tradeiD

		var RESPONSE2 StockPortfolioDetails

		client, err := net.Dial("tcp", "127.0.0.1:5555")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		c := jsonrpc.NewClient(client)
		err = c.Call("Computer.CheckPortfolio", REQUEST2, &RESPONSE2 )
		if err != nil {
			log.Fatal("arith error:", err)
		}

		fmt.Println("Below are the Portfolio details :-")

		fmt.Print("Stocks :")
		fmt.Println(RESPONSE2.Stocks)
		fmt.Print("Current Market Value :")
		fmt.Println(RESPONSE2.CurrentMarketValue)
		fmt.Print("Unvested Amount :")
		fmt.Println(RESPONSE2.UnvestedAmount)

	} else {

		fmt.Println("Wrong option entered ...")
	}



}
