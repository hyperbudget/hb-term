package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	ui "github.com/gizak/termui"
)

type transaction struct {
	Date         string
	DebitAmount  float64
	CreditAmount float64
	Description  string
}

type txnListResponse struct {
	Email        string        `json:"email"`
	ID           string        `json:"id"`
	Transactions []transaction `json:"transactions"`
}

func getTransactions(username string, password string) txnListResponse {
	values := map[string]string{"email": username, "password": password}

	jsonValue, _ := json.Marshal(values)

	resp, _ := http.Post(
		"http://api.hyperbudget.net/account/transactions/list",
		"application/json",
		bytes.NewBuffer(jsonValue),
	)

	var txn txnListResponse

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &txn)

	return txn
}

func floatToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

func getUserAndPass() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username\n")

	username, _ := reader.ReadString('\n')
	username = strings.TrimSuffix(username, "\n")

	fmt.Print("Enter password\n")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSuffix(password, "\n")

	return username, password
}

func main() {
	username, password := getUserAndPass()

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	rows1 := [][]string{
		[]string{"Date", "Descr", "Deb", "Cred"},
	}

	res := getTransactions(username, password)

	txns := res.Transactions

	// reverse
	for i := len(txns)/2 - 1; i >= 0; i-- {
		opp := len(txns) - 1 - i
		txns[i], txns[opp] = txns[opp], txns[i]
	}

	// txns := res.Transactions[:10]

	for _, txn := range txns {
		rows1 = append(rows1, []string{
			txn.Date,
			txn.Description,
			floatToString(txn.DebitAmount),
			floatToString(txn.CreditAmount),
		})
	}

	table1 := ui.NewTable()
	table1.Rows = rows1
	table1.FgColor = ui.ColorWhite
	table1.BgColor = ui.ColorDefault
	table1.Y = 0
	table1.X = 0
	table1.Width = 150
	table1.Height = len(rows1)*2 + 1

	ui.Render(table1)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
