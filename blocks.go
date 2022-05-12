package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/k0kubun/pp"
)

const wsapi = `wss://ws.s0.t.hmny.io`

func wsgo() {
	http.HandleFunc("/", wslisten)
}

func wslisten(w http.ResponseWriter, r *http.Request) {
	c, e := upgrader.Upgrade(w, r, nil)
	fmt.Println(e)
	for {
		var an interface{}
		e := c.ReadJSON(&an)
		fmt.Println(e)
		fmt.Println(an)
	}
}

// json2sql. we can subscribe to new blocks instead of polling every 2.2 secs or something. go websockets
func wsjson2sql() {
	c, _, e := websocket.DefaultDialer.Dial(wsapi, nil)
	if e != nil {
		log.Println(e)
		return
	}
	defer c.Close()

	b, e := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "hmyv2_subscribe",
		"params": []interface{}{
			"newBlocks",
		},
	})
	if e != nil {
		log.Println(e)
		return
	}
	e = c.WriteJSON(b)
	w := &wsc{c}
	go w.read()
	fmt.Println(e)
	select {}
}

type wsc struct {
	*websocket.Conn
}

func (c *wsc) read() {
	for {
		var an interface{}
		fmt.Println("reading")
		e := c.ReadJSON(&an)
		fmt.Println(e)
		fmt.Println(an)
	}
}

func Blocks(from, to int) {

	aa := do(js(GetBlocksApi, from, to, BlockArgs))
	if aa == nil {
		return
	}
	a := aa.([]interface{})
	we := len(a)

	for ib, bb := range a {

		var yu Wpack
		yu.Switch = 0

		b := bb.(map[string]interface{})
		bnum := int64(gf(b["number"]))
		epoch := int(gf(b["epoch"]))
		bhash := gs(b["hash"])
		parentHash := gs(b["parentHash"])
		nonce := int(gf(b["nonce"]))
		size := int64(gf(b["size"]))
		mixHash := gs(b["mixHash"])
		logsBloom := gs(b["logsBloom"])
		stateRoot := gs(b["stateRoot"])
		miner := onex(gs(b["miner"]))
		difficulty := int(gf(b["difficulty"]))
		extraData := gs(b["extraData"])
		gasLimit := int64(gf(b["gasLimit"]))
		gasUsed := int64(gf(b["gasUsed"]))
		times := int64(gf(b["timestamp"]))
		troot := gs(b["tranactionsRoot"])
		rroot := gs(b["receiptsRoot"])
		transactions := b["transactions"].([]interface{})
		txCount := len(transactions)
		fromshard := int(gf(b["shardID"]))
		toshard := int(gf(b["toShardID"]))
		uncles := b["uncles"].([]interface{})
		unclecount := len(uncles)
		stakingTransactions := b["stakingTransactions"].([]interface{})
		stakingCount := len(stakingTransactions)
		var signers []string
		if b["signers"] != nil {
			signs := b["signers"].([]interface{})
			for _, s := range signs {
				signers = append(signers, s.(string))
			}
		}
		signersCount := len(signers)

		fmt.Printf("\rProcessed %d/%d. Txs: %d", ib, we, txCount)

		diff := int64(0)
		if bnum > 0 {
			var ptime int64
			e := db.QueryRow(`select timestamp from blocks2 where id = ?`, bnum-1).Scan(&ptime)
			if e == sql.ErrNoRows {
				e = nil
			}
			if e != nil {
				log.Panic(e)
			}
			if ptime != 0 {
				diff = times - ptime
			}
		}

		_, e := db.Exec(`insert into blocks2 (
		id, epoch, hash, parentHash, nonce, size, mixHash, logsBloom, stateRoot, miner, difficulty, extraData, gasLimit, gasUsed, 
		timestamp, tdiff, transactionsRoot, receiptsRoot, fromshard, toshard, txcount, unclecount, stakingcount, signerscount) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			bnum, epoch, bhash, parentHash, nonce, size, mixHash, logsBloom, stateRoot, miner, difficulty, extraData, gasLimit, gasUsed, times, diff, troot, rroot, fromshard, toshard, txCount, unclecount, stakingCount, signersCount)

		fmt.Println("putting block. ", e)
		if driverErr, ok := e.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 {
				e = nil
			}
		}

		if e != nil {
			log.Panic(e)
		}

		_, e = db.Exec(`insert into accs (address, bnum, thash) values (?,?,?)`, miner, bnum, ``)
		if driverErr, ok := e.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 {
				e = nil
			}
		}
		if e != nil {
			log.Panic(e)
		}

		// Make HTML

		now := time.Now().Unix()
		secspan := fmt.Sprintf(`<span class="timing">%d</span>`, now-times)

		iu := fmt.Sprintf(BlockRow, fmt.Sprintf("%d.%d", fromshard, toshard), ifmt(bnum), shash(bhash), txCount, stakingCount, signersCount, sizefmt(size), shash(miner), secspan)

		yu.Blocks = iu

		if txCount != 0 {

			for _, a := range transactions {

				tx := a.(map[string]interface{})

				thash := gs(tx["hash"])
				time.Sleep(time.Second / 8)
				rec := Txr(thash)
				if rec == nil {
					fmt.Println("skipping ", thash)
					continue
				}
				bhash, bnum := gs(tx["blockHash"]), gf(tx["blockNumber"])
				afrom, ato := onex(gs(tx["from"])), onex(gs(tx["to"]))
				vone := float2atto(gf(tx["value"]))
				timestamp := gf(tx["timestamp"])
				status := gf(rec["status"]) == 1.0
				gas, gasPrice, gasUsed := gf(tx["gas"]), gf(tx["gasPrice"]), gf(rec["gasUsed"])
				cumulativeGasUsed := gf(rec["cumulativeGasUsed"])
				nonce := int64(gf(tx["nonce"]))
				caddr := gs(rec["contractAddress"])
				ainput := gs(tx["input"])
				txIndex := gf(tx["transactionIndex"])
				ethHash := gs(tx["ethHash"])
				logsBloom := gs(rec["logsBloom"])
				shardID, toShardID := int(gf(tx["shardID"])), int(gf(tx["toShardID"]))
				logs := rec["logs"].([]interface{})
				logCount := len(logs)
				root := gs(rec["root"])
				v := gs(tx["v"])
				r := gs(tx["r"])
				s := gs(tx["s"])

				if gasPrice > 1<<63-1 {
					gasPrice = 1<<63 - 1
				}

				// 13826104 was 1e19
				// you will now need to check if all those with max values have more than max.

				_, e := db.Exec(`insert into txs2 (
				hash, blockHash, blockNumber, afrom, ato, vone, timestamp, status, gas, gasPrice, gasUsed, cumulativeGasUsed,
				nonce, contractAddress, input, transactionIndex, ethHash, logsBloom, shardID, toShardID, logCount, root, v, r, s
				) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
					thash, bhash, bnum, afrom, ato, vone, timestamp, status, gas, gasPrice, gasUsed, cumulativeGasUsed,
					nonce, caddr, ainput, txIndex, ethHash, logsBloom, shardID, toShardID, logCount, root, v, r, s)

				if e != nil {
					log.Panic(e)
				}

				age := time.Now().Unix() - int64(timestamp)
				po := fmt.Sprintf(TxRow, fmt.Sprintf("%d.%d", shardID, toShardID), shash(thash), shash(afrom), shash(ato), age, logCount, 0)

				yu.Txs += po

				// Contracts

				if caddr != `` {
					cq := cname(caddr)
					if cq != `` {
						symbol := csymbol(caddr)
						supply := csupply(caddr)
						decimals := cdec(caddr)
						if len(supply) > 100 {
							supply = supply[:100]
						}
						_, e := db.Exec(`insert into contracts (address, txhash, bhash, name, supply, decimals, symbol) values (?,?,?,?,?,?,?)`,
							caddr, thash, bhash, cq, supply, decimals, symbol)
						if e != nil {
							log.Panic(e)
						}

						t0 := call(Token0, caddr)
						if t0 != `` {
							t0 = pname(t0)
							t1 := pname(call(Token1, caddr))
							_, e := db.Exec(`insert into pairs (tokenone, tokentwo, address, txhash, bhash, name, supply, decimals, symbol) values (?,?,?,?,?,?,?,?,?)`,
								t0, t1, caddr, thash, bhash, cq, supply, decimals, symbol)
							if e != nil {
								log.Panic()
							}
						}

					}
				}

				_, e = db.Exec(`insert into accs (address, bnum, thash) values (?,?,?)`, afrom, bnum, thash)
				if driverErr, ok := e.(*mysql.MySQLError); ok {
					if driverErr.Number == 1062 {
						e = nil
					}
				}
				if e != nil {
					log.Panic(e)
				}
				_, e = db.Exec(`insert into accs (address, bnum, thash) values (?,?,?)`, ato, bnum, thash)
				if driverErr, ok := e.(*mysql.MySQLError); ok {
					if driverErr.Number == 1062 {
						e = nil
					}
				}
				if e != nil {
					log.Panic(e)
				}

				if logCount != 0 {

					// This is almost where the actual thinking begins.
					// Logs are events emitted by Solidity. So they span contracts and functions.
					// Things we care about like liquidity adds/subs token swaps are list of events.
					// Because some tokens are bridged and the nature of DEX cross chain etc
					// the sequences will differ when one swap hmochi for one or some other token for one.
					// Just because its X -> WONE != Y -> Wone or x -> y != y -> x etc..
					// they only difference within these types through time is the version changes of code
					// how do we build these patterns without hard coding everything?

					// Whats hardcpdomg vwso;dsklsd;kfl;sdkf;lsdkf

					for _, a := range logs {
						log := a.(map[string]interface{})

						address := gs(log["address"])

						logIndex := gs(log["logIndex"])
						var logIdx int64
						if logIndex != `` {
							a, _ := new(big.Int).SetString(logIndex[2:], 16)
							logIdx = a.Int64()
							pp.Print(logIdx)
						}
						pp.Print(logIndex)
						valid := !gb(log["removed"])
						data := gs(log["data"])
						topics := log["topics"].([]interface{})
						topicCount := len(topics)

						var eventsig, in1, in2, in3 string
						if topicCount > 0 {
							for i, a := range topics {
								cc := a.(string)
								if i == 0 {
									eventsig = cc
								}
								if i == 1 {
									in1 = cc
								}
								if i == 2 {
									in2 = cc
								}
								if i == 3 {
									in3 = cc
								}
							}
						}

						_, e := db.Exec(`insert into logs2 (address, aindex, txhash, valid, data, topicCount, eventsig, inone, intwo, inthree) values (?,?,?,?,?,?,?,?,?,?)`, address, logIdx, thash, valid, data, topicCount, eventsig, in1, in2, in3)
						if e != nil {
							fmt.Print(e)
							os.Exit(1)
						}

					}
				}

			}

		}

		if stakingCount != 0 {
			for _, a := range stakingTransactions {

				stx := a.(map[string]interface{})
				shash := gs(stx["hash"])
				stype := gs(stx["type"])
				msg := stx["msg"].(map[string]interface{})
				gas := int64(gf(stx["gas"]))
				gasPrice := int64(gf(stx["gasPrice"]))
				daddr := gs(msg["delegatorAddress"])
				var atto, vaddr string
				if stype == `Delegate` {
					atto = float2atto(gf(msg["amount"]))
					vaddr = gs(msg["validatorAddress"])
				}

				_, e := db.Exec(`insert into stxs (thash, gas, gasPrice, stype, vaddr, daddr, atto) values (?,?,?,?,?,?,?)`, shash, gas, gasPrice, stype, vaddr, daddr, atto)
				if e != nil {
					log.Panic(e)
				}
			}
		}

		// ws.Broadcast <- yu

		_ = yu

	}

}
