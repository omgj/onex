package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func newHub() *Ws {
	return &Ws{
		con:  make(chan *wscon),
		on:   make(map[*wscon]bool),
		dcon: make(chan *wscon),
		send: make(chan wpack),
	}
}

type Ws struct {
	con  chan *wscon
	on   map[*wscon]bool
	dcon chan *wscon
	send chan wpack
}

type wpack struct {
	Switch int `json:"switch"`
	Lat0   int `json:"lat0"`
	Lat1   int `json:"lat1"`
	Lat2   int `json:"lat2"`
	Lat3   int `json:"lat3"`
	Lat4   int `json:"lat4"`
	Lat5   int `json:"lat5"`
}

type wscon struct {
	con  *websocket.Conn
	send chan wpack
}

func (w *wscon) read() {
	defer func() {
		// hub.dcon <- w
	}()
	for {
		_, m, e := w.con.ReadMessage()
		if e != nil {
			log.Println(e)
			return
		}
		_ = m
	}
}
func (w *wscon) write() {
	t := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-t.C:
			qq := wpack{
				Switch: 0,
				Lat0:   1,
				Lat1:   1,
				Lat2:   1,
				Lat3:   1,
				Lat4:   1,
				Lat5:   1,
			}
			log.Println("sending ", qq)
			e := w.con.WriteJSON(qq)
			if e != nil {
				log.Println(e)
				t.Stop()
				return
			}
		}
	}
}

func (w *Ws) run() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case c := <-w.con:
			w.on[c] = true
		case c := <-w.dcon:
			close(c.send)
			c.con.Close()
			delete(w.on, c)
			fmt.Println("deleted")
		case b := <-w.send:
			for aa := range w.on {
				aa.send <- b
			}
			fmt.Printf("sent new block to %d peers\n", len(w.on))
		}
	}
}

func websock(w http.ResponseWriter, r *http.Request) {
	a, e := upgrader.Upgrade(w, r, nil)
	if e != nil {
		log.Print(e)
		return
	}
	q := &wscon{
		con:  a,
		send: make(chan wpack),
	}
	// hub.con <- q
	go q.read()
	go q.write()
	log.Println("Reading & Writing Websockets for: ", r.RemoteAddr)
}
