package main

import (
	"database/sql"
	"log"
)

// Get size of tbls SELECT TABLE_SCHEMA AS `Database`, TABLE_NAME AS `Table`, ROUND((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024) AS `Size (MB)` FROM information_schema.TABLES ORDER BY (DATA_LENGTH + INDEX_LENGTH) ASC;

const sqlBlock = `create table if not exists blocks (
	id BIGINT PRIMARY KEY,
	epoch INT,
	hash VARCHAR(66),
	parentHash VARCHAR(66),
	nonce BIGINT,
	size BIGINT,
	mixHash VARCHAR(100),
	logsBloom VARCHAR(1000),
	stateRoot VARCHAR(66),
	miner VARCHAR(44),
	difficulty INT,
	extraData VARCHAR(1000),
	gasLimit BIGINT,
	gasUsed BIGINT,
	timestamp BIGINT,
	tdiff BIGINT,
	transactionsRoot VARCHAR(100),
	receiptsRoot VARCHAR(100),
	fromshard INT,
	toshard INT,
	txcount INT,
	unclecount INT,
	stakingcount INT,
	signerscount INT
)`

const sqlLog = `create table if not exists logs (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	address VARCHAR(44),
	aindex INT,
	txhash VARCHAR(66),
	valid BOOL,
	data TEXT,
	topicCount INT,
	eventsig VARCHAR(66),
	inone VARCHAR(66),
	intwo VARCHAR(66),
	inthree VARCHAR(66)
)`

const sqlTx = `create table if not exists txs (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	hash VARCHAR(66),
	blockHash VARCHAR(66),
	blockNumber INT,
	afrom VARCHAR(44),
	ato VARCHAR(44),
	vone VARCHAR(100),
	timestamp BIGINT,
	status BOOL,
	gas BIGINT,
	gasPrice BIGINT,
	gasUsed BIGINT,
	cumulativeGasUsed BIGINT,
	nonce BIGINT,
	contractAddress VARCHAR(100),
	input TEXT,
	transactionIndex INT,
	ethHash VARCHAR(100),
	logsBloom TEXT,
	shardID INT,
	toShardID INT,
	logCount INT,
	root TEXT,
	v VARCHAR(100),
	r VARCHAR(100),
	s VARCHAR(100)
)`

const sqlAccs = `create table if not exists accs (
	address VARCHAR(44) PRIMARY KEY,
	bnum BIGINT,
	thash VARCHAR(66)
)`

const sqlContracts = `create table if not exists contracts (
	idx BIGINT AUTO_INCREMENT,
	address VARCHAR(44) PRIMARY KEY,
	creator VARCHAR(44),
	bhash VARCHAR(66),
	txhash VARCHAR(66),
	name VARCHAR(100),
	symbol VARCHAR(100),
	supply VARCHAR(200),
	decimals INT,
	tone VARCHAR(44),
	ttwo VARCHAR(44),
	ccode TEXT
)`

const sqlPair = `create table if not exists pairs (
	address VARCHAR(44),
	tokenone VARCHAR(100),
	tokentwo VARCHAR(100),
	creator VARCHAR(44),
	bhash VARCHAR(66),
	txhash VARCHAR(66),
	name VARCHAR(100),
	symbol VARCHAR(100),
	supply VARCHAR(100),
	decimals VARCHAR(100),
	ccode TEXT
)`

// -----------------------------------------------
// ----------- TRANSACTION & RECEIPT -------------
// -----------------------------------------------

const sqlStx = `create table if not exists stxs (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	thash VARCHAR(66),
	gas BIGINT,
	gasPrice BIGINT,
	stype VARCHAR(20),
	vaddr VARCHAR(44),
	daddr VARCHAR(44),
	atto VARCHAR(100)
)`

// -----------------------------------------------
// ---------------------- LOG --------------------
// -----------------------------------------------

// -----------------------------------------------
// ----------------- ACCOUNTS --------------------
// -----------------------------------------------

const sqlSigners = `create table if not exists signers (
   address VARCHAR(44),
   block BIGINT
)`

// -----------------------------------------------
// ----------------- CONTRACT --------------------
// -----------------------------------------------

const sqlTokens = `create table if not exists tokens (
	id INT PRIMARY KEY AUTO_INCREMENT,
	address VARCHAR(44),
	creator VARCHAR(44),
	bhash VARCHAR(66),
	txhash VARCHAR(66),
	name VARCHAR(100),
	symbol VARCHAR(100),
	supply VARCHAR(200),
	decimals VARCHAR(100),
	ccode TEXT
)`

const sqlCal = `create table if not exists cal (
	year INT,
	month INT,
	day INT,
	hour INT,
	bcount INT,
	arows MEDIUMTEXT
)`

// Prior Version

const sqlBalances = `CREATE TABLE IF NOT EXISTS balances (
	acc BIGINT,
	bnum BIGINT,
	token BIGINT,
	bal VARCHAR(100),

)`

const sqlSwaps = `CREATE TABLE IF NOT EXISTS swaps (
	acc VARCHAR(40),
	afrom BIGINT,
	ato BIGINT,
	ain VARCHAR(100),
	aout VARCHAR(100)
)`

const sqlReserves = `CREATE TABLE IF NOT EXISTS reserves (
	pair BIGINT,
	t0 VARCHAR(100),
	t1 VARCHAR(100),
	btime BIGINT
)`

// InsertReserves. If btime exists nothing will be inserted. Btime is updated on chain
// via _update whereafter a Sync Event is emitted.
func InsertReserves(pair int, t0, t1 string, btime int64) {
	var q0, q1 string
	e := db.QueryRow(`select t0, t1 from reserves where btime = ?`, btime).Scan(&q0, &q1)
	if e == sql.ErrNoRows {
		e = nil
	}
	if e != nil {
		log.Println(e)
		return
	}
	_, e = db.Exec(`insert into reserves (pair, t0, t1, btime) values (?,?,?,?)`, pair, t0, t1, btime)
	if e != nil {
		log.Print(e)
		return
	}
}

const sqlPairs = `CREATE TABLE IF NOT EXISTS pairs (
	pair BIGINT,
	t0 BIGINT,
	t1 BIGINT,
	thash VARCHAR(66)
)`

const sqlIP = `CREATE TABLE IF NOT EXISTS ips (
	aip VARCHAR(30),
	atime BIGINT,
	uagent VARCHAR(200)
)`
const sqlAcc = `CREATE TABLE IF NOT EXISTS accs (
	bech VARCHAR(38),
	ox VARCHAR(40)
)`

func Maketable() {
	if _, e := db.Exec(sqlPairs); e != nil {
		log.Fatal(e)
	}
}
