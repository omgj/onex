package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var lastinweb int
var ws *Host
var upgrader = websocket.Upgrader{}
var db *sql.DB
var booting bool

type br struct {
	brow string `json:"brow"`
	unix int64  `json:"unix"`
}

func MakeAllContracts() {
	// a, e := db.Query(`select distinct(afrom) from txs order by bnum desc`)
	// if e != nil {
	// 	fmt.Println(e)
	// 	return
	// }
	// for a.Next() {
	// 	var b string
	// 	e := a.Scan(&b)
	// 	if e != nil {
	// 		fmt.Println(e)
	// 		return
	// 	}
	// 	checker(b)
	// }

	// a.Close()

	// a1, e := db.Query(`select distinct(ato) from txs order by bnum desc`)
	// if e != nil {
	// 	fmt.Println(e)
	// 	return
	// }
	// defer a1.Close()
	// for a1.Next() {
	// 	var b string
	// 	e := a1.Scan(&b)
	// 	if e != nil {
	// 		fmt.Println(e)
	// 		return
	// 	}
	// 	checker(b)
	// }

	a2, e := db.Query(`select contractAddress from txs where contractAddress != '' order by blockNumber desc`)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer a2.Close()
	for a2.Next() {
		var b string
		e := a2.Scan(&b)
		if e != nil {
			fmt.Println(e)
			return
		}
		checker(b)
	}
}

func checker(b string) {
	xx := cname(b)
	xz := csymbol(b)
	if xx == `` || xz == `` {
		return
	}
	fmt.Println("Inserting: ", b)
	_, e := db.Exec(`insert into contracts (name, symbol) values (?,?)`, xx, xz)
	if e != nil {
		fmt.Println(e)
		return
	}
}

func EpochCalculator() {
	nextlastblock := int64(do(js(`hmy_epochLastBlock`, GetEpochMeth())).(float64))
	currentblock := high()
	fmt.Printf("Blocks until next epoch: %d\n", nextlastblock-currentblock)
	fmt.Printf("Latency %s s\t So epoch ends in T minus %d mins\n", lat(100), (nextlastblock-currentblock)/60)
}

func lat(an int) string {
	var l float64
	e := db.QueryRow(`select avg(tdiff) from blocks order by id desc limit 100`).Scan(&l)
	if e != nil {
		fmt.Println(e)
		return ``
	}

	return fmt.Sprintf("%.2f", l)
}

func main() {
	l := low()
	h := high() - 1
	ws = NewHost()
	go ws.run()
	go stats()
	go webs()
	for i := l; i <= h; i = i + 500 {
		Blocks(int(i), int(i)+500)
	}
}

var Months = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

const TxCal = `select shardID, toShardID, hash, afrom, ato, timestamp, logCount, gas, gasPrice, vone, input, blockHash, blockNumber, status, contractAddress from txs order by timestamp asc limit 1000`
const Tcalrow = `<td>%s</td> <td>%s</td> <td>%s</td> <td>%s</td> <td>%s</td> <td>%d</td> <td>%d</td> <td>%s</td>`

// func MakeTxTable() {
// 	a, e := db.Query(TxCal)
// 	if e != nil {
// 		fmt.Println(e)
// 		return
// 	}
// 	var yp, mp, dp, hp int
// 	tt := ``
// 	count := 0
// 	for a.Next() {
// 		var shardID, toShardID, timestamp, logCount, gas, gasPrice, blockNumber int64
// 		var hash, afrom, ato, vone, input, blockHash, contractAddress string
// 		var stat bool
// 		e := a.Scan(&shardID, &toShardID, &hash, &afrom, &ato, &timestamp, &logCount, &gas, &gasPrice, &vone, &input, &blockHash, &blockNumber, &stat, &contractAddress)
// 		if e != nil {
// 			fmt.Println(e)
// 			return
// 		}

// 		b := time.Unix(timestamp, 0)
// 		year := b.Year()
// 		month := Months[b.Month().String()]
// 		day := b.Day()
// 		hour := b.Hour()
// 		count++
// 		shards := fmt.Sprintf("%d.%d", fromshard, toshard)
// 		if gasPrice == 0 {
// 			gas = gas*gasPrice
// 		}
// 		if input == `` {

// 		}

// 		if yp == 0 {
// 			yp, mp, dp, hp = year, month, day, hour
// 		}

// 		if yp == year {
// 			if mp == month {
// 				if dp == day {
// 					if hp == hour {
// 						tt += fmt.Sprintf(Tcalrow, shards, shash(hash), shash(afrom), shash(ato), vone, gas, )
// 						continue
// 					}
// 					// new hour
// 					e := db.Exec(`insert into tcal (year, month, day, hour, tcount, arows) values (?,?,?,?,?,?)`, year, month, day, hour, count, tt)
// 					if e != nil {
// 						fmt.Println(e)
// 						return
// 					}
// 					tt = ``
// 					count = 1
// 					hp = hour
// 					tt +=
// 					continue
// 				}
// 				// new day
// 				e := db.Exec(`insert into tcal (year, month, day, hour, tcount, arows) values (?,?,?,?,?,?)`, year, month, day, hour, count, tt)
// 					if e != nil {
// 						fmt.Println(e)
// 						return
// 					}
// 					tt = ''
// 					count = 1
// 					hp, dp = hour, day
// 					tt +=
// 					continue
// 			}

// 			// new month
// 			e := db.Exec(`insert into tcal (year, month, day, hour, tcount, arows) values (?,?,?,?,?,?)`, year, month, day, hour, count, tt)
// 					if e != nil {
// 						fmt.Println(e)
// 						return
// 					}
// 					tt = ''
// 					count = 1
// 					hp, dp, mp = hour, day, month
// 					tt +=
// 					continue

// 		}

// 		// new year
// 		e := db.Exec(`insert into tcal (year, month, day, hour, tcount, arows) values (?,?,?,?,?,?)`, year, month, day, hour, count, tt)
// 					if e != nil {
// 						fmt.Println(e)
// 						return
// 					}
// 					tt = ''
// 					count = 1
// 					hp, dp, mp, yp = hour, day, month, year
// 					tt +=
// 					continue

// 	}
// 	defer a.Close()
// }

func MakeTokenHolders() {
	a, e := db.Query(`select contractAddress from txs where contractAddress != '';`)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer a.Close()

	var cons []string
	for a.Next() {
		var as string
		e := a.Scan(&as)
		if e != nil {
			fmt.Println(e)
			return
		}
		cons = append(cons, as)
	}

	for _, c := range cons {

		GetTxHistory(c)

	}

}

// HTTP HANDLES

func MostRecentContractID() (id int64) {
	e := db.QueryRow(`select id from contracts order by id desc limit 1`).Scan(&id)
	if e != nil {
		fmt.Println(e)
		return
	}
	return
}

func MostRecentLogID() (id int64) {
	e := db.QueryRow(`select id from logs order by id desc limit 1`).Scan(&id)
	if e != nil {
		fmt.Println(e)
		return
	}
	return
}

func CurrentEpoch() (epoch int64) {
	e := db.QueryRow(`select epoch from blocks order by id desc limit 1`).Scan(&epoch)
	if e != nil {
		fmt.Println(e)
		return
	}
	return
}

func stats() {
	t := time.NewTicker(time.Second * 3)
	defer t.Stop()

	for {
		select {
		case <-t.C:

			l := low()
			lt := tlow()
			al := alow()

			last20 := GetLast20()
			l12 := Last20Txs()
			rec := ifmt(MostRecentContractID())
			re1c := ifmt(MostRecentLogID())
			qp := ifmt(CurrentEpoch())
			// ii := CirculatingSupply()
			oo := ifmt(int64(GetAllValidatorAddressess()))

			p := Wpack{
				Switch:    1,
				Bcount:    ifmt(l),
				Tcount:    ifmt(lt),
				Acount:    ifmt(al),
				Blocks:    last20,
				Txs:       l12,
				Contracts: rec,
				Logs:      re1c,
				Epoch:     qp,
				// Circ:      ii,
				Vals: oo,
			}

			for _, r := range ws.Active {
				r.send <- p
			}

		}
	}
}

func GetLast20() (blocks string) {

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
		rb := fmt.Sprintf("<tr>%s</tr>", BlockRow)
		shards := fmt.Sprintf("%d.%d", fromshard, toshard)
		times := fmt.Sprintf(`<span class="timing">%d</span>`, time.Now().Unix()-timestamp)
		blocks += fmt.Sprintf(rb, shards, ifmt(id), shash(hash), txCount, stakingcount, signercount, sizefmt(size), shash(miner), times)
	}

	a.Close()
	return
}

func sizefmt(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d b", size)
	}
	if size >= 1024 {
		return fmt.Sprintf("%.2f kB", (float64(size) / 1024))
	}
	return fmt.Sprintf("%d b", size)
}

func Last20Txs() (txs string) {

	a, e := db.Query(`select shardID, toShardID, hash, afrom, ato, timestamp, logCount, gas, gasPrice, vone, input from txs order by id desc limit 10`)
	if e != nil {
		fmt.Println(e)
		return
	}

	for a.Next() {
		var shardID, toShardID, timestamp, logCount, gas, gasPrice int64
		var hash, afrom, ato, vone, inputs string
		e := a.Scan(&shardID, &toShardID, &hash, &afrom, &ato, &timestamp, &logCount, &gas, &gasPrice, &vone, &inputs)
		if e != nil {
			fmt.Println(e)
			return
		}
		var gf string
		if gasPrice == 0 {
			gf = fmt.Sprintf("%s Atto", ifmt(gas))
		} else {
			y := fmt.Sprintf("%d", gas*gasPrice)
			u := len(y)
			for i := (18 - u); i > 0; i-- {
				y = `0` + y
			}
			gf = fmt.Sprintf(".%s", rtrim(y))
		}
		q := fmt.Sprintf("<tr>%s</tr>", TxRow)
		qq := len(vone)
		if qq > 18 {
			vone = rtrim(sifmt(vone[:qq-18]) + `.` + vone[qq-18:])
		}
		if qq == 18 {
			vone = `.` + rtrim(vone)
		}
		txs += fmt.Sprintf(q, fmt.Sprintf("%d.%d", shardID, toShardID), shash(hash), shash(afrom), shash(ato), vone, time.Now().Unix()-timestamp, logCount, gf)
	}

	a.Close()

	return
}

const BlockRow = `<td scope="row">%s</td> <td>%s</td> <td>%s</td> <td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td>`
const TxRow = `<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%s</td>`

type BlockWithFullTxArgs struct {
	WithSigners bool     `json:"withSigners"`
	InclTx      bool     `json:"inclTx"`
	FullTx      bool     `json:"fullTx"`
	Signers     []string `json:"signers"`
	InclStaking bool     `json:"inclStaking"`
}

var BlockArgs = BlockWithFullTxArgs{
	FullTx:      true,
	InclTx:      true,
	InclStaking: true,
	WithSigners: true,
}
