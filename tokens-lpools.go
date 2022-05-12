package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/go-ethereum/common"
	"github.com/go-sql-driver/mysql"
	"github.com/k0kubun/pp"
)

// PairTokens. MochiPair.token0/1()
func MakePair(address string, an int64) {
	t0id := maketoken(trim(call(`token0()`, address)), ``)
	t1id := maketoken(trim(call(`token1()`, address)), ``)

	fmt.Println("put pair ", an)
	_, e := db.Exec(`insert into pairs (pair, t0, t1, thash) values (?,?,?,?)`, an, t0id, t1id, ``)
	if e != nil {
		log.Println(e)
		return
	}
}

const swaps = `0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822`

// PairSwaps. MochiPair emits Swap(msg.sender address,amount0In unint256, amount1In unint256, amount0Out unint256, amount1Out uint256, to address)

// GetSwapLogs. 2/3 topics. data has 4 numbers in 256 bits. so far 3 log tables.
func Swaps() {
	a, e := db.Query(`select address, inone, intwo, data from logs2 where eventsig = ? limit 1000`, swaps)
	if e != nil {
		log.Println(e)
		return
	}
	defer a.Close()
	for a.Next() {
		var address, inone, intwo, data string
		e := a.Scan(&address, &inone, &intwo, &data)
		if e != nil {
			fmt.Println(e)
			return
		}
		a0in := trim(data[2:][:64])
		a1in := trim(data[2:][64:128])
		a0out := trim(data[2:][128 : 128+64])
		a1out := trim(data[2:][128+64 : 128+128])
		var in, out string
		var intok, outtok, pid int64
		e = db.QueryRow(`select id from tokens where zerox = ?`, address[2:]).Scan(&pid)
		if e == sql.ErrNoRows {
			fmt.Println("no token, ", address)
			pid = maketoken(address[2:], ``)
			MakePair(address, pid)
			e = nil
		}
		if e != nil {
			e = db.QueryRow(`select id from tokens where zerox = ?`, address[2:]).Scan(&pid)
			if e != nil {
				log.Println(e)
				return
			}
			return
		}
		if a0in != `` {
			in = a0in
			e := db.QueryRow(`select id from tokens t inner join pairs p on p.t0 = t.id where pair = ?`, pid).Scan(&intok)
			if e != nil {
				log.Panic(e)
				return
			}
		}
		if a1in != `` {
			in = a1in
			e := db.QueryRow(`select id from tokens t inner join pairs p on p.t1 = t.id where pair = ?`, pid).Scan(&intok)
			if e != nil {
				log.Panic(e)
				return
			}
		}
		if a0out != `` {
			out = a0out
			e := db.QueryRow(`select id from tokens t inner join pairs p on p.t0 = t.id where pair = ?`, pid).Scan(&outtok)
			if e != nil {
				log.Panic(e)
				return
			}
		}
		if a1out != `` {
			out = a1out
			e := db.QueryRow(`select id from tokens t inner join pairs p on p.t1 = t.id where pair = ?`, pid).Scan(&outtok)
			if e != nil {
				log.Panic(e)
				return
			}
		}

		fmt.Println("inserting for pair ", pid)
		_, e = db.Exec(`insert into swaps (acc, afrom, ato, ain, aout) values (?,?,?,?,?)`, trim(intwo), intok, outtok, in, out)
		if e != nil {
			fmt.Println(e)
			return
		}
	}
}

// PairsAll. Returns a list of [40]string
func PairsAll() (pairs []string) {
	a, e := db.Query(`select zerox from tokens t inner join pairs p on p.pair = t.id`)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer a.Close()
	for a.Next() {
		var q string
		e := a.Scan(&q)
		if e != nil {
			fmt.Println(e)
			return
		}
		pairs = append(pairs, q)
	}
	fmt.Println("Pairs Found: ", len(pairs))
	return
}

// TokensAll.
func TokensAll() (toks []string) {
	a, e := db.Query(`select zerox from tokens`)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer a.Close()
	for a.Next() {
		var q string
		e := a.Scan(&q)
		if e != nil {
			fmt.Println(e)
			return
		}
		toks = append(toks, q)
	}
	return
}

// TokenHistory.
func TokenHistory(address string) {
	targs := TxHistoryArgs{
		Address:   `0x` + address,
		PageIndex: 0,
		PageSize:  1000,
		FullTx:    true,
		TxType:    "ALL",
		Order:     "ASC",
	}
	var Txs []map[string]interface{}
	for i := 0; i < 1000; i++ {
		targs.PageIndex = uint32(i)
		t := do(js(`hmyv2_getTransactionsHistory`, targs))
		if t == nil {
			i = 1000
			continue
		}
		tt := t.(map[string]interface{})
		if tt["error"] != nil {
			fmt.Println(tt["error"].(string))
			i = 1000
			continue
		}
		txs := tt["transactions"].([]interface{})
		yy := len(txs)
		if yy == 0 {
			i = 1000
			continue
		}
		for _, q := range txs {
			qq := q.(map[string]interface{})
			Txs = append(Txs, qq)
		}
		fmt.Printf("\rToken: %s\t Txs: %d", address, len(Txs))
	}
	pp.Print(Txs)
}

// GetReserves. MochiPair.getReserves() uint112 _reserve0, uint112 _reserve1, uint32 _blockTimestampLast
// 0x0902f1ac5dbaeedd3217f11b3cbaf929216c9c5abc2d69da89d54964bead575d
func getreserves(pair string, an int64) {

	a := call(`getReserves()`, pair)
	b := len(a)
	if b == 192 {
		qq := trim(a[:64])
		ww := trim(a[64:128])
		ee := trim(a[128:])
		q, _ := new(big.Int).SetString(qq, 16)
		w, _ := new(big.Int).SetString(ww, 16)
		e, _ := new(big.Int).SetString(ee, 16)
		fmt.Printf("Reserve0:  %s\t hex: %s\n", q.String(), qq)
		fmt.Printf("Reserve1:  %s\t hex: %s\n", w.String(), ww)
		fmt.Printf("Timestamp: %s\t hex: %s\n", e.String(), ee)

		_, r := db.Exec(`insert into reserves (pair, t0, t1, btime) values (?,?,?,?)`, an, q.String(), w.String(), e.Int64())
		if r != nil {
			log.Println(e)
			return
		}
	}
}

// TotalPairs. MochiFactory.allPairsLength() uint
// 0x574f2ba3b1a925e685747ccc4a1a928462c32f9e78a901d91c2ad75498800f78
func TotalPairs() {}

// PairCreations. MochiFactory.createPair(address token0, address token1) emits PairCreated(address,address,address,uint256)
// It will make the tokens if they dont exist. maketoken()
func PairCreations() {

	r, e := db.Query(`select txhash, inone, intwo, data from logs2 where eventsig = ?`, sigPairCreated)
	if e != nil {
		log.Println(e)
		return
	}
	defer r.Close()
	for r.Next() {
		var txhash, inone, intwo, data string
		e := r.Scan(&txhash, &inone, &intwo, &data)
		if e != nil {
			log.Println(e)
			return
		}
		t0 := trim(inone)
		t1 := trim(intwo)
		ad := trim(data)
		if len(ad) > 40 {
			ad = ad[:40]
		}
		fmt.Printf("Token0: %s\n", t0)
		fmt.Printf("Token1: %s\n", t1)
		fmt.Printf("Pair:   %s\n", ad)

		t0id := maketoken(t0, ``)
		t1id := maketoken(t1, ``)
		adid := maketoken(ad, txhash)

		log.Println("Let's put the new token pair address in.")

		_, e = db.Exec(`insert into pairs (pair, t0, t1, thash) values (?,?,?,?)`, adid, t0id, t1id, txhash)
		if e != nil {
			log.Println(e)
			return
		}

		fmt.Println("Let's put the initial reserve measure for the pool.")
	}
}

// earlier func
func maketoken(a string, thash string) int64 {
	pp.Println(a)
	var t1count int64
	new := false
	e := db.QueryRow(`select id from tokens where zerox = ?`, a).Scan(&t1count)
	if driverErr, ok := e.(*mysql.MySQLError); ok {
		if driverErr.Number == 1062 {
			new = true
			e = nil
		}
	}
	if e == sql.ErrNoRows {
		new = true
		e = nil
	}
	if e != nil {
		log.Panic(e)
		return 0
	}
	if !new {
		log.Println("already in ", t1count)
		return t1count
	}

	done := true
	if thash == `` {
		e := db.QueryRow(`select hash from txs where contractAddress = ?`, `0x`+a).Scan(&thash)
		if e == sql.ErrNoRows {
			e = nil
			done = false
		}
		if e != nil {
			fmt.Println(a)
			log.Panic(e)
			return 0
		}
		fmt.Printf("Getting Tx Hash: %s", thash)
	}

	var afrom, bhash string
	var gas, gasPrice int64
	if done {
		e = db.QueryRow(`select blockHash, afrom, gas, gasPrice from txs where hash = ?`, thash).Scan(&bhash, &afrom, &gas, &gasPrice)
		if e == sql.ErrNoRows {
			e = nil
		}
		if e != nil {
			log.Println(e)
			return 0
		}
	}

	fee := gas
	if gasPrice != 0 {
		fee = gas * gasPrice
	}

	var id int64
	if done {
		e = db.QueryRow(`select id from blocks where hash = ?`, bhash).Scan(&id)
		if e == sql.ErrNoRows {
			e = nil
		}
		if e != nil {
			log.Println(e)
			return 0
		}
	}

	log.Println(afrom)

	name := cname(a)
	symbol := csymbol(a)
	decimals := cdec(a)
	supplyhex, supply := TokenSupply(a)
	one1 := xone(common.Hex2Bytes(a))[4:]
	var maker1 string
	if done {
		if len(thash) > 0 {
			thash = thash[2:]
		}
		if len(bhash) > 0 {
			bhash = bhash[2:]
		}
		if afrom != `` {
			afrom = afrom[2:]
			maker1 = xone(common.Hex2Bytes(afrom[2:]))[4:]
		}
	}

	log.Println("Inserting")
	q, e := db.Exec(`insert into tokens (zerox, one1, bhash, bnum, thash, fee, makerx, maker1, name, symbol, decimals, supply, supplyhex) values (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		a, one1, bhash, id, thash, fee, afrom, maker1, strings.TrimLeft(name, " "), strings.Trim(symbol, " "), decimals, supply, supplyhex)
	if e != nil {
		log.Println(e)
		return 0
	}

	i, e := q.LastInsertId()
	if e != nil {
		log.Println(e)
		return 0
	}
	return i
}

func TokenSupply(address string) (string, string) {
	a := call(`totalSupply()`, `0x`+address)
	hex := trim(a)
	ints, _ := new(big.Int).SetString(hex, 16)
	return hex, ints.String()
}
