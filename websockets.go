package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/olekukonko/tablewriter"
)

type Wpack struct {
	Switch    int    `json:"switch"`
	Blocks    string `json:"blocks"`
	Txs       string `json:"txs"`
	Contracts string `json:"contracts"`
	Logs      string `json:"logs"`
	Circ      string `json:"circ"`
	Vals      string `json:"vals"`
	Epoch     string `json:"epoch"`
	Bcount    string `json:"bcount"`
	Tcount    string `json:"tcount"`
	Acount    string `json:"acount"`
	Toks      string `json:"toks"`
	Cal       string `json:"cal"`
}

type ConnConfig struct {
	id   int
	conn *websocket.Conn
	send chan Wpack
	Txs  bool
}

type Host struct {
	Active     map[int]*ConnConfig
	Broadcast  chan Wpack
	Disconnect chan int
	Connect    chan *ConnConfig
	HighBlock  string
	HighTx     string
	HighAcc    string
	HighNet    int64
}

func IsCurl(a string) bool {
	return strings.Contains(a, "curl")
}

// func BlockHandle(w http.ResponseWriter, r *http.Request) {
// 	a := strings.Split(r.URL.String(), "/")
// 	if len(a) == 1 {
// 		return
// 	}
// 	if len(a) == 2 {
// 		if IsCurl(r.UserAgent()) {
// 			ts := &strings.Builder{}
// 			t := tablewriter.NewWriter(ts)
// 			t.SetHeader([]string{"Shards", "Index", "Hash", "Txs", "Size", "Age"})
// 			d := l20()
// 			for _, a := range
// 		}
// 	}
// }

func webs() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", HomeHandle)

	// http.HandleFunc("/block/", BlockHandle)
	// http.HandleFunc("/tx/", TxHandle)
	// http.HandleFunc("/acc/", AccHandle)
	// http.HandleFunc("/tok/", TokHandle)
	// http.HandleFunc("/pair/", PairHandle)

	http.HandleFunc("/ws/", WSHandle)
	log.Println("Website on 5localhost:1234. Waiting 4u...")
	log.Fatal(http.ListenAndServe(":1234", nil))
}

func HomeHandle(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.RemoteAddr)
	fmt.Println(r.UserAgent())

	if strings.Contains(r.UserAgent(), "curl") {
		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)
		table.SetHeader([]string{"Shards", "Index", "Hash", "Txs", "StakingTxs", "Signers", "Size", "Miner", "Age"})

		data := Last20Curl()

		for _, v := range data {
			table.Append(v)
		}

		table.Render()

		w.Write([]byte(tableString.String()))
		return
	}

	q, e := ioutil.ReadFile("index.html")
	if e != nil {
		log.Panic(e)
	}
	_, e = w.Write(q)
	if e != nil {
		log.Panic(e)
	}
}

var lastindex int

func WSHandle(w http.ResponseWriter, r *http.Request) {
	c, e := upgrader.Upgrade(w, r, nil)
	if e != nil {
		log.Panic(e)
	}
	lastindex++
	cc := &ConnConfig{
		conn: c,
		send: make(chan Wpack),
		id:   lastindex,
	}
	ws.Connect <- cc

	go cc.read()
	go cc.write()
}

func (cc *ConnConfig) read() {
	fmt.Printf("Establishing Read Routing for %d\n", cc.id)
	defer cc.conn.Close()
	for {
		_, a, e := cc.conn.ReadMessage()
		if e != nil {
			ws.Disconnect <- cc.id
			return
		}
		b := string(a)
		q := strings.Split(b, "-")
		if len(q) == 1 {
			v := GetCalYears()
			i := Wpack{
				Switch: 5,
				Cal:    v,
			}
			cc.send <- i
		}
		if len(q) == 2 {
			qq, e := strconv.Atoi(q[1])
			if e != nil {
				fmt.Println(e)
				return
			}
			v := GetYearRow(qq)
			i := Wpack{
				Switch: 6,
				Cal:    v,
			}
			cc.send <- i
		}
		if len(q) == 3 {
			qq, e := strconv.Atoi(q[1])
			if e != nil {
				fmt.Println(e)
				return
			}
			qqq, e := strconv.Atoi(q[2])
			if e != nil {
				fmt.Println(e)
				return
			}
			v := GetMonthRow(qq, qqq)
			i := Wpack{
				Switch: 7,
				Cal:    v,
			}
			cc.send <- i
		}
		if len(q) == 4 {
			qq, e := strconv.Atoi(q[1])
			if e != nil {
				fmt.Println(e)
				return
			}
			qqq, e := strconv.Atoi(q[2])
			if e != nil {
				fmt.Println(e)
				return
			}
			qqq1, e := strconv.Atoi(q[3])
			if e != nil {
				fmt.Println(e)
				return
			}
			v := GetDayRow(qq, qqq, qqq1)
			i := Wpack{
				Switch: 8,
				Cal:    v,
			}
			cc.send <- i
		}
		if len(q) == 5 {
			qq, e := strconv.Atoi(q[1])
			if e != nil {
				fmt.Println(e)
				return
			}
			qqq, e := strconv.Atoi(q[2])
			if e != nil {
				fmt.Println(e)
				return
			}
			qqq1, e := strconv.Atoi(q[3])
			if e != nil {
				fmt.Println(e)
				return
			}
			q1, e := strconv.Atoi(q[4])
			if e != nil {
				fmt.Println(e)
				return
			}
			v := GetPageRow(qq, qqq, qqq1, q1)
			i := Wpack{
				Switch: 9,
				Cal:    v,
			}
			cc.send <- i
		}

	}
}

const (
	ContractRow = `<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`
)

func (cc *ConnConfig) write() {
	defer cc.conn.Close()

	rows, e := db.Query(`select address, symbol, name, supply, decimals, bhash, txhash from contracts`)
	if e != nil {
		return
	}

	var toks string
	for rows.Next() {
		var address, symbol, name, supply, decimals, bhash, txhash string
		rows.Scan(&address, &symbol, &name, &supply, &decimals, &bhash, &txhash)
		o := len(supply)
		if o > 40 {
			supply = supply[:40]
			o = 40
		}
		// TODO change decimals column to int and remove this.
		if decimals != `` {
			i, e := strconv.Atoi(decimals)
			if e != nil {
				fmt.Println(e)
			}
			if o > i {
				if o == 40 {
					supply = supply[:6] + `..`
				} else {
					supply = supply[:o-i]
				}
			}
		}
		toks += fmt.Sprintf(ContractRow, shash(address), symbol, name, supply, decimals, shash(bhash), shash(txhash))
	}

	rows.Close()

	p := Wpack{
		Switch: 2,
		Toks:   toks,
	}

	e = cc.conn.WriteJSON(p)
	if e != nil {
		return
	}

	for {
		select {
		case a := <-cc.send:
			e := cc.conn.WriteJSON(a)
			if e != nil {
				return
			}
		}
	}
}

// WEBSOCKETS

func NewHost() *Host {
	return &Host{
		Active:     make(map[int]*ConnConfig),
		Broadcast:  make(chan Wpack),
		Disconnect: make(chan int),
		Connect:    make(chan *ConnConfig),
	}
}

func (h *Host) run() {
	fmt.Printf("Running the Ws Host.")
	for {
		select {
		case cc := <-h.Connect:
			h.Active[cc.id] = cc
			fmt.Printf("Added a New Connection: %d\n", cc.id)
		case cc := <-h.Disconnect:
			h.Active[cc].conn.Close()
			close(h.Active[cc].send)
			delete(h.Active, cc)
			fmt.Printf("Deleted Connection: %d\n", cc)
		case m := <-h.Broadcast:
			fmt.Printf("Broadcasting Message to %d Peers.", len(h.Active))
			for _, a := range h.Active {
				a.send <- m
			}
		}
	}
}
