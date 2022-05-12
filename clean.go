package main

import (
	"fmt"
)

func Clean() {
	BlocksLinear()
	BlocksHaveTxs()
	TxsHaveLogs()
}

func SwapsFromLogs() {
	txs := GetAllTxs()

	q := len(txs)

	for o, z := range txs {

		fmt.Printf("\r %d/%d", o, q)

		ls := GetLogsForTxs(z)

		for _, l := range ls {

			_ = l

		}

	}
}

// BlocksLinear. Did we skip an index?
func BlocksLinear() {
	fmt.Println(" Checking linear sequence by parentHash.")
	a, e := db.Query(`select id, parentHash from blocks order by id asc`)
	if e != nil {
		return
	}

	var p int
	for a.Next() {
		var id int
		var phash string
		e := a.Scan(&id, &phash)
		if e != nil {
			return
		}
		_ = phash
		if id != p+1 && id != 0 {
			fmt.Printf("Jump Detected: %d\n", id)
		}
		if id != 0 {
			p++
		}
		fmt.Printf("\r %d", id)
	}

	a.Close()
}

// BlocksHaveTxs. Block.txCount == select count(*) from txs where bhash = ?
func BlocksHaveTxs() {

	a, e := db.Query(`select hash, txCount from blocks`)
	if e != nil {
		return
	}

	var q []btxsan
	for a.Next() {
		var hash string
		var txs int
		e := a.Scan(&hash, &txs)
		if e != nil {
			return
		}
		z := btxsan{hash, txs}
		q = append(q, z)
	}

	a.Close()

	qq := len(q)
	for u, r := range q {
		fmt.Printf("\r %d/%d", u, qq)
		var have int
		e := db.QueryRow(`select count(*) from txs where blockHash = ?`, r.hash).Scan(&have)
		if e != nil {
			return
		}
		if have != r.txcount {
			fmt.Printf("\nFound an error with: %s\n", r.hash)
		}
	}
}

// TxsHaveLogs. tx.logCount == select count(*) from logs where thash = ?
func TxsHaveLogs() {
	a, e := db.Query(`select hash, logCount from txs`)
	if e != nil {
		return
	}

	var tt []btxsan
	for a.Next() {
		var hash string
		var qq int
		e := a.Scan(&hash, &qq)
		if e != nil {
			return
		}
		t := btxsan{hash, qq}
		tt = append(tt, t)
	}

	a.Close()

	rr := len(tt)

	for z, i := range tt {
		fmt.Printf("\r %d/%d", z, rr)
		var y int
		e := db.QueryRow(`select count(*) from logs where thash = ?`, i.hash).Scan(&y)
		if e != nil {
			return
		}
		if y != i.txcount {
			fmt.Printf("Error with tx: %s\n", i.hash)
		}
	}
}

type btxsan struct {
	hash    string
	txcount int
}

// ReplenishPairs.
func ReplenishContracts() {

	i := GetAllContracts()
	ii := len(i)

	for u, a := range i {

		fmt.Printf(" Address: %s\t#%d/%d", a, u, ii)

	}
}

type contract struct {
	name, supply, symbol string
	decimals             int
}

func GetAllContracts() (contracts []string) {
	rows, e := db.Query(`select contractAddress from txs where contractAddress != ''`)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var a string
		e := rows.Scan(&a)
		if e != nil {
			return
		}
		contracts = append(contracts, a)
	}
	return
}

func GetAllTxs() (txs []string) {
	rows, e := db.Query(`select hash from txs`)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var a string
		e := rows.Scan(&a)
		if e != nil {
			return
		}
		txs = append(txs, a)
	}
	return
}

type Log struct {
	Data     string
	EventSig string
	One      string
	Two      string
	Three    string
}

func GetLogsForTxs(tx string) (logs []Log) {

	a, e := db.Query(`select data, eventsig, inone, intwo, inthree from logs where txhash = ?`, tx)
	if e != nil {
		return
	}

	for a.Next() {
		var data, eventsig, inone, intwo, inthree string
		e := a.Scan(&data, &eventsig, &inone, &intwo, &inthree)
		if e != nil {
			return
		}
		l := Log{data, eventsig, inone, intwo, inthree}
		logs = append(logs, l)
	}

	a.Close()

	return
}
