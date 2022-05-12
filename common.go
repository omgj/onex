package main

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-sql-driver/mysql"
	"github.com/k0kubun/pp"
)

func float2atto(a float64) string {
	return strconv.FormatFloat(a, 'f', -1, 64)
}

func dfmt(decimals string) string {
	var a string
	var zeros string
	for _, b := range decimals {
		if string(b) == `0` {
			zeros += `0`
			continue
		}
		if len(zeros) > 0 {
			a += zeros + string(b)
			zeros = ``
			continue
		}
		a += string(b)
	}
	return a
}

func do(a []byte) interface{} {
	q, e := http.NewRequest("POST", "https://api.s0.t.hmny.io", bytes.NewBuffer(a))
	if e != nil {
		log.Fatal(e)
	}
	q.Header.Add("Content-Type", "application/json")
	w := http.Client{}
	w.Timeout = time.Second * 100
	r, e := w.Do(q)
	if e != nil {
		log.Panic(e)
	}
	u, e := ioutil.ReadAll(r.Body)

	if e != nil {
		log.Println(e)
		return nil
	}
	defer r.Body.Close()
	s := make(map[string]interface{})

	if len(u) == 0 {
		return nil
	}

	if u[0] == 0x7b {
		e = json.Unmarshal(u, &s)
		if e != nil {
			log.Panic(e)
		}
		if s["error"] != nil {
			fmt.Println(s["error"])
			fmt.Println("ERROR ERROR ERROR")
			return nil
		}
		return s["result"]
	}

	if u[0] == 0x3c {
		pp.Println(string(u))
		return nil
	}
	return nil
}

// misordering a fails the call
func js(q string, a ...interface{}) []byte {
	if len(a) == 0 {
		b, e := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  q,
			"params":  []interface{}{},
		})
		if e != nil {
			log.Panic(e)
		}
		return b
	}
	b, e := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  q,
		"params":  a,
	})
	if e != nil {
		log.Panic(e)
	}
	return b
}

func trim(a string) string {
	yes := true
	trimmed := 0
	if len(a) < 2 {
		return ""
	}
	a = a[2:]
	for yes {
		if len(a) == 0 {
			return ""
		}
		if string(a[0]) != "0" {
			yes = false
			continue
		}
		if string(a[0]) == "0" {
			trimmed++
			a = a[1:]
		}
	}
	return a
}

func keccakfull(a string) string {
	return crypto.Keccak256Hash([]byte(a)).String()
}

func keccak(a string) string {
	return crypto.Keccak256Hash([]byte(a)).String()[:10]
}

func xone(address []byte) string {
	b, e := bech32.ConvertBits(address, 8, 5, true)
	if e != nil {
		log.Println(e)
		return ``
	}
	a, e := bech32.Encode("one", b)
	if e != nil {
		log.Println(e)
		return ``
	}
	return a
}

func onex(address string) string {
	_, a, _ := bech32.Decode(address)
	b, e := bech32.ConvertBits(a, 5, 8, false)
	if e != nil {
		log.Panic(e)
	}
	return common.BytesToAddress(b).String()
}

func err(e error) {
	if driverErr, ok := e.(*mysql.MySQLError); ok {
		if driverErr.Number == 1062 {
			e = nil
		}
	}
	if e != nil {
		log.Fatal(e)
	}
}

func hex2atto(hexnum string) string {
	a := new(big.Int)
	a, _ = a.SetString(trim(hexnum), 16)
	return a.String()
}

func gf(a interface{}) float64 {
	if a == nil {
		return 0
	}
	return a.(float64)
}

func gs(a interface{}) string {
	if a == nil {
		return ``
	}
	b := trim(a.(string))
	if b == `` {
		return ``
	}
	return a.(string)
}

func gb(a interface{}) bool {
	if a == nil {
		return false
	}
	return a.(bool)
}

func sifmt(a string) string {
	b := len(a)
	if b < 4 {
		return a
	}
	if b == 4 {
		return fmt.Sprintf("%s,%s", a[:1], a[1:])
	}
	if b == 5 {
		return fmt.Sprintf("%s,%s", a[:2], a[2:])
	}
	if b == 6 {
		return fmt.Sprintf("%s,%s", a[:3], a[3:])
	}
	if b == 7 {
		return fmt.Sprintf("%s,%s,%s", a[:1], a[1:4], a[4:])
	}
	if b == 8 {
		return fmt.Sprintf("%s,%s,%s", a[:2], a[2:5], a[5:])
	}
	if b == 9 {
		return fmt.Sprintf("%s,%s,%s", a[:3], a[3:6], a[6:])
	}
	return a
}

func ifmt(bnum int64) string {
	if bnum < 1000 {
		return fmt.Sprintf("%d", bnum)
	}
	bs := fmt.Sprintf("%d", bnum)
	if bnum < 10000 && bnum > 999 {
		return fmt.Sprintf("%s,%s", bs[:1], bs[1:])
	}
	if bnum < 100000 && bnum > 9999 {
		return fmt.Sprintf("%s,%s", bs[:2], bs[2:])
	}
	if bnum < 1000000 && bnum > 99999 {
		return fmt.Sprintf("%s,%s", bs[:3], bs[3:])
	}
	if bnum < 10000000 && bnum > 999999 {
		return fmt.Sprintf("%s,%s,%s", bs[:1], bs[1:4], bs[4:7])
	}
	if bnum < 100000000 && bnum > 9999999 {
		return fmt.Sprintf("%s,%s,%s", bs[:2], bs[2:5], bs[5:8])
	}
	return ``
}

func init() {
	var e error
	db, e = sql.Open("mysql", "me:testing@/one")
	if e != nil {
		log.Panic(e)
	}
}

func createtables() {

	_, e := db.Exec(sqlBlock)
	if e != nil {
		log.Panic(e)
	}
	_, e = db.Exec(sqlTx)
	if e != nil {
		log.Panic(e)
	}
	_, e = db.Exec(sqlLog)
	if e != nil {
		log.Panic(e)
	}
	_, e = db.Exec(sqlContracts)
	if e != nil {
		log.Panic(e)
	}
	_, e = db.Exec(sqlPair)
	if e != nil {
		log.Panic(e)
	}
	_, e = db.Exec(sqlAccs)
	if e != nil {
		log.Panic(e)
	}
	_, e = db.Exec(sqlStx)
	if e != nil {
		log.Panic(e)
	}
	_, e = db.Exec(sqlSigners)
	if e != nil {
		log.Panic(e)
	}
	_, e = db.Exec(sqlTokens)
	if e != nil {
		log.Panic(e)
	}
	_, e = db.Exec(sqlCal)
	if e != nil {
		log.Panic(e)
	}
}

func low() int64 {
	var id int64
	e := db.QueryRow(`select id from blocks2 order by id desc limit 1`).Scan(&id)
	if e == sql.ErrNoRows {
		return 0
	}
	fmt.Printf("Fetching Highest Block from Database: %d\n", id)
	return id
}

func tlow() int64 {
	var id int64
	e := db.QueryRow(`select id from txs order by id desc limit 1`).Scan(&id)
	if e == sql.ErrNoRows {
		return 0
	}
	return id
}

func alow() int64 {
	var id int
	e := db.QueryRow(`select count(*) from accs`).Scan(&id)
	if e == sql.ErrNoRows {
		return 0
	}
	return int64(id)
}

func high() int64 {
	i := int64(gf(do(js(NetHighBlockApi))))
	fmt.Printf("Fetcing highet block fomr the Net: %d\n", i)
	return i
}

func pname(a string) string {

	if len(trim(a)) == 0 {
		return ``
	}

	b := a[8:]
	off := b[:64]
	l := b[64:128]
	ll := b[128:]

	i, _ := new(big.Int).SetString(trim(off), 16)
	o, _ := new(big.Int).SetString(trim(l), 16)

	fmt.Println(i)
	fmt.Println(o)

	e, _ := hex.DecodeString(ll[:i.Int64()+(o.Int64()*2)])
	return string(e)

}

func shash(hash string) string {
	if len(hash) == 0 {
		return ``
	}
	return fmt.Sprintf(`<a href="%s">%s</a>`, hash, hash[2:8]+`.`+hash[len(hash)-6:])
}

func bhash(hash string) string {
	return hash[2:8] + `.` + hash[len(hash)-6:]
}

func rtrim(a string) string {
	b := []rune(a)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	c := trim(string(b))
	b = []rune(c)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)

}

func Last20Curl() (blocks [][]string) {
	a, e := db.Query(`select fromshard, toshard, id, hash, txCount, stakingcount, signerscount, size, miner, timestamp from blocks order by id desc limit 20`)
	if e != nil {
		fmt.Println(e)
		return
	}
	for a.Next() {
		var fromshard, toshard, id, txCount, stakingcount, signercount, size, timestamp int64
		var hash, miner string
		e := a.Scan(&fromshard, &toshard, &id, &hash, &txCount, &stakingcount, &signercount, &size, &miner, &timestamp)
		if e != nil {
			fmt.Println(e)
			return
		}
		shards := fmt.Sprintf("%d.%d", fromshard, toshard)
		blocks = append(blocks, []string{shards, ifmt(id), hash, strconv.Itoa(int(txCount)), strconv.Itoa(int(stakingcount)), strconv.Itoa(int(signercount)), sizefmt(size), miner, fmt.Sprintf("%d", time.Now().Unix()-timestamp)})
	}

	a.Close()
	return
}

func readableseconds(input int) (result string) {
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	if years > 0 {
		result = pl(int(years), "y") + pl(int(months), "m") + pl(int(weeks), "w") + pl(int(days), "d") + pl(int(hours), "h") + pl(int(minutes), "m") + pl(int(seconds), "s")
	} else if months > 0 {
		result = pl(int(months), "m") + pl(int(weeks), "w") + pl(int(days), "d") + pl(int(hours), "h") + pl(int(minutes), "m") + pl(int(seconds), "s")
	} else if weeks > 0 {
		result = pl(int(weeks), "w") + pl(int(days), "d") + pl(int(hours), "h") + pl(int(minutes), "m") + pl(int(seconds), "s")
	} else if days > 0 {
		result = pl(int(days), "d") + pl(int(hours), "h") + pl(int(minutes), "m") + pl(int(seconds), "s")
	} else if hours > 0 {
		result = pl(int(hours), "h") + pl(int(minutes), "m") + pl(int(seconds), "s")
	} else if minutes > 0 {
		result = pl(int(minutes), "m") + pl(int(seconds), "s")
	} else {
		result = pl(int(seconds), "s")
	}

	return
}

func pl(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
		result = strconv.Itoa(count) + singular + " "
	} else {
		result = strconv.Itoa(count) + singular + " "
	}
	return
}
