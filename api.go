package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/k0kubun/pp"
)

const (
	CircSupply                = `hmy_getCirculatingSupply`
	TotalSupplyAPI            = `hmy_getTotalSupply`
	GetEpoch                  = `hmyv2_getEpoch`
	PendingTxAPI              = `hmyv2_pendingTransactions`
	BlockBalanceAPI           = `hmyv2_getBalanceByBlockNumber`
	TxrApi                    = `hmyv2_getTransactionReceipt`
	TxApi                     = `hmyv2_getTransactionByHash`
	TxHistoryApi              = `hmyv2_getTransactionsHistory`
	GetBlocksApi              = `hmyv2_getBlocks`
	NetHighBlockApi           = `hmyv2_blockNumber`
	CallApi                   = `hmyv2_call`
	GetCode                   = `hmy_getCode`
	GetValidators             = `hmyv2_getValidators`
	BlockByHash               = `hmyv2_getBlockByHash`
	StakingNetInfo            = `hmy_getStakingNetworkInfo`
	GetAllValidatorAddresses  = `hmy_getAllValidatorAddresses`
	GetAllNValidatorAddresses = `hmy_getActiveValidatorAddresses`
	GetValidatorInfo          = `hmy_getValidatorInformation`
	GetTransactionCount       = `hmyv2_getTransactionCount`
)

func ValidatorsByEpoch(epoch int) {
	a := do(js(GetValidators, epoch)).(map[string]interface{})
	gf(a["shardID"])
	validators := a["validators"].([]interface{})
	for _, b := range validators {
		b := b.(map[string]interface{})
		gs(b["address"])
		gf(b["balance"]) // stake
	}
}

func GetEpochMeth() int64 {
	return int64(do(js(GetEpoch)).(float64))
}

func GetStakingNetworkInfo() {
	a := do(js(StakingNetInfo)).(map[string]interface{})
	gs(a["circulating-supply"])
	_ = int64(gf(a["epoch-last-block"]))
	gs(a["median-raw-stake"])
	float2atto(gf(a["total-staking"]))
	gs(a["total-supply"])
}

func GetAllValidatorAddressess() int {
	a := do(js(GetAllValidatorAddresses)).([]interface{})
	if len(a) == 0 {
		return 0
	}
	return len(a)
}
func GetAllNValidatorAddressess() {
	a := do(js(GetAllNValidatorAddresses)).([]interface{})
	if len(a) == 0 {
		return
	}
	pp.Println(len(a))
}

func GetValidatorInfos(vaddr string) {
	a := do(js(GetValidatorInfo, vaddr)).(map[string]interface{})
	gs(a["active-status"])
	gs(a["booted-status"])
	_ = a["current-epoch-performance"]
	_ = a["currently-in-committee"].(bool)
	gs(a["epos-status"])
	_ = a["epos-winning-stake"]
	lifetime := a["lifetime"].(map[string]interface{})
	gs(lifetime["apr"])
	lifetimeblocks := lifetime["blocks"].(map[string]interface{})
	gf(lifetimeblocks["signed"])
	gf(lifetimeblocks["to-sign"])
	_ = lifetime["epoch-apr"]
	_ = lifetime["epoch-blocks"]
	gf(lifetime["reward-accumulated"])
	_ = a["metrics"]
	gf(a["total-delegation"])
	validator := a["validator"].(map[string]interface{})
	gs(validator["address"])
	var blskeys []string
	vbls := validator["bls-public-keys"].([]interface{})
	if len(vbls) > 0 {
		for _, a := range vbls {
			blskeys = append(blskeys, a.(string))
		}
	}
	_ = blskeys
	gf(validator["creation-height"])
	delegations := validator["delegations"].([]interface{})
	if len(delegations) > 0 {
		for _, a := range delegations {
			aa := a.(map[string]interface{})
			gf(aa["amount"])
			gs(aa["delegator-address"])
			gf(aa["reward"])
			_ = len(aa["delegator-address"].([]interface{}))
		}
	}
	gs(a["details"])
	gs(a["identity"])
	gf(a["last-epoch-in-committee"])
	gs(a["max-change-rate"])
	gs(a["max-rate"])
	gf(a["max-total-delegation"])
	gf(a["max-self-delegation"])
	gs(a["name"])
	gs(a["rate"])
	gs(a["security-contract"])
	gf(a["update-height"])
	gs(a["website"])

}

func TxCountAtBlock(address string, block int) {
	gf(do(js(GetTransactionCount, address, block)))
}

func CirculatingSupply() string {
	a := strings.Split(do(js(CircSupply)).(string), ".")[0]
	q := len(a)
	if q < 8 {
		return a
	}
	if q == 9 {
		return fmt.Sprintf("%s,%s,%s", a[:3], a[3:6], a[6:9])
	}
	if q == 10 {
		return fmt.Sprintf("%s,%s,%s,%s", a[:1], a[1:4], a[4:7], a[7:10])
	}
	if q == 11 {
		return fmt.Sprintf("%s,%s,%s,%s", a[:2], a[2:5], a[5:8], a[8:11])
	}
	if q == 12 {
		return fmt.Sprintf("%s,%s,%s,%s", a[:3], a[3:6], a[6:9], a[9:12])
	}
	return a
}

func TotalSupply() string {
	return do(js(TotalSupplyAPI)).(string)
}

func PendingTx() {
	pp.Print(do(js(PendingTxAPI)))
}

func BalanceAtBlock(address string, block int) string {
	a := do(js(BlockBalanceAPI, address, block))
	if a == nil {
		return ``
	}
	return float2atto(a.(float64))
}

func Txr(tx string) map[string]interface{} {
	a := do(js(TxrApi, tx))
	if a != nil {
		aa := a.(map[string]interface{})
		return aa
	}
	return nil
}

func Tx(tx string) map[string]interface{} {
	a := do(js(TxApi, tx)).(map[string]interface{})
	return a
}

type TxHistoryArgs struct {
	Address   string `json:"address"`
	PageIndex uint32 `json:"pageIndex"`
	PageSize  uint32 `json:"pageSize"`
	FullTx    bool   `json:"fullTx"`
	TxType    string `json:"txType"`
	Order     string `json:"order"`
}

func GetTxHistory(address string) {
	a := TxHistoryArgs{
		Address:   address,
		PageIndex: 0,
		PageSize:  1000,
		FullTx:    true,
		TxType:    "ALL",
		Order:     "ASC",
	}
	for i := 0; i < 1000; i++ {
		a.PageIndex = uint32(i)
		b := do(js(TxHistoryApi, a))
		pp.Print(b)
		if b == nil {
			i = 1000
			continue
		}
		bb := b.(map[string]interface{})
		tx := bb["transactions"].([]interface{})
		l := len(tx)
		if l == 0 {
			i = 1000
			continue
		}
		for _, t := range tx {
			tt := t.(map[string]interface{})

			pp.Print(tt)
			i = 1000
			continue
			// blocks = append(blocks, int(tt["blockNumber"].(float64)))
			// addr = append(addr, onex(tt["from"].(string)))
			// u := do(js("hmyv2_getTransactionReceipt", tt["hash"].(string))).(map[string]interface{})
			// iii := onex(tt["from"].(string))
			// var yess int
			// err(db.QueryRow(`select count(*) from bb where address = ?`, iii).Scan(&yess))
			// if yess > 0 {
			// 	page++
			// 	fmt.Printf("%d\n", page)
			// 	continue
			// }
			// bbb := int(tt["blockNumber"].(float64))
			// _, e := db.Exec(`insert into bb (address, block) values (?,?)`, iii, bbb)
			// err(e)
			// page++
			// fmt.Printf("%d\n", page)
		}
	}
	return
}

type CallArgs struct {
	From     *common.Address `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Data     *hexutil.Bytes  `json:"data"`
}

func call(method, contract string) string {
	k := do(js(CallApi, callargs(method, contract, ""), "latest"))
	if k == nil {
		return ""
	}
	kk := k.(string)
	if len(kk) == 2 {
		return ""
	}
	if len(kk[2:])%2 != 0 {
		return ""
	}
	return kk[2:]
}

func callargs(method string, contract string, user string) CallArgs {
	bhash := keccak(method)
	var a CallArgs
	q := common.HexToAddress(contract)
	a.To = &q
	var w hexutil.Bytes
	if user == "" {
		w.UnmarshalText([]byte(bhash))
		a.Data = &w
		return a
	}
	w.UnmarshalText([]byte(bhash + `000000000000000000000000` + user[2:]))
	a.Data = &w
	return a
}

// TokenName. HRC20.name() 0x06fdde03
func cname(address string) string {
	a := call(NameMethod, address)
	v := len(a)
	if v == 0 {
		return ``
	}

	if len(trim(a)) == 0 {
		return ``
	}

	// qq, _ := strconv.ParseInt(a[:64], 16, 64)
	qw, _ := strconv.ParseInt(a[64:128], 16, 64)
	if qw > 100 {
		b, _ := hex.DecodeString(a[128:])
		return string(b)
	}
	b, _ := hex.DecodeString(a[128 : 128+qw*2])
	return string(b)
}

// TokenSymbol. HRC20.symbol() 0x95d89b41
func csymbol(address string) string {
	a := call(SymbolMethod, address)
	if len(a) == 0 {
		return ``
	}
	qw, _ := strconv.ParseInt(a[64:128], 16, 64)
	if qw > 100 {
		b, _ := hex.DecodeString(a[128:])
		return string(b)
	}
	b, _ := hex.DecodeString(a[128 : 128+qw*2])
	return string(b)
}

// TokenSupply. HRC20.totalSupply() 0x18160ddd.
func csupply(address string) string {
	a := call(TotalSupplyMethod, address)
	if len(a) == 0 {
		return ``
	}
	b, _ := new(big.Int).SetString(trim(a), 16)
	return b.String()
}

// TokenDecimals. HRC20.decimals() 0x313ce567
func cdec(address string) int64 {
	a := call(DecimalsMethod, address)
	if len(a) == 0 {
		return 0
	}
	b := trim(a)
	if b == `` {
		return 0
	}
	sb, _ := new(big.Int).SetString(b, 16)
	return sb.Int64()
}
