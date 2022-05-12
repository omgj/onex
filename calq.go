package main

import (
	"fmt"
	"strings"
	"time"
)

func GetCalYears() (row string) {
	a, e := db.Query(`select distinct year from cal`)
	if e != nil {
		fmt.Println(e)
		return
	}

	m1 := `<span id="%d" onclick="yearc(this.id)">%d</span>&nbsp;`

	for a.Next() {
		var q int
		e := a.Scan(&q)
		if e != nil {
			fmt.Println(e)
			return
		}
		row += fmt.Sprintf(m1, q, q)
	}

	a.Close()
	return
}

func GetYearRow(year int) (row string) {
	a, e := db.Query(`select distinct month from cal where year = ?`, year)
	if e != nil {
		fmt.Println(e)
		return
	}

	m1 := `<span id="%s" onclick="monthc(this.id)">%d</span>&nbsp;`

	for a.Next() {
		var m int
		e := a.Scan(&m)
		if e != nil {
			fmt.Println(e)
			return
		}
		row += fmt.Sprintf(m1, fmt.Sprintf("%d-%d", year, m), m)
	}

	a.Close()
	return
}

func GetMonthRow(year, month int) (row string) {
	a, e := db.Query(`select distinct day from cal where year = ? and month = ?`, year, month)
	if e != nil {
		fmt.Println(e)
		return
	}

	m1 := `<span id="%s" onclick="dayc(this.id)">%d</span>&nbsp;`

	for a.Next() {
		var d int
		e := a.Scan(&d)
		if e != nil {
			fmt.Println(e)
			return
		}
		row += fmt.Sprintf(m1, fmt.Sprintf("%d-%d-%d", year, month, d), d)
	}

	a.Close()
	return
}

func GetDayRow(year, month, day int) (row string) {
	a, e := db.Query(`select hour from cal where year = ? and month = ? and day = ?`, year, month, day)
	if e != nil {
		fmt.Println(e)
		return
	}

	m1 := `<span id="%s" onclick="pagec(this.id)">%d</span>&nbsp;`

	for a.Next() {
		var p int
		e := a.Scan(&p)
		if e != nil {
			fmt.Println(e)
			return
		}
		row += fmt.Sprintf(m1, fmt.Sprintf("%d-%d-%d-%d", year, month, day, p), p)
	}

	a.Close()
	return
}

func GetPageRow(year, month, day, page int) (row string) {
	e := db.QueryRow(`select arows from cal where year = ? and month = ? and day = ? and hour = ?`, year, month, day, page).Scan(&row)
	if e != nil {
		fmt.Println(e)
		return
	}
	return
}

func GetTxHour(year, month, day, hour int) (row string) {
	e := db.QueryRow(`select arows from tcal where year = ? and month = ? and day = ? and hour = ?`, year, month, day, hour).Scan(&row)
	if e != nil {
		fmt.Println(e)
		return
	}
	return
}

func MakeBlockPagingTable() {

	a, e := db.Query(`select fromshard, toshard, id, hash, txCount, stakingcount, signerscount, size, timestamp, nonce, epoch from blocks order by id asc limit 10000`)
	if e != nil {
		fmt.Println(e)
		return
	}

	BlockRow1 := fmt.Sprintf("<tr>%s</tr>", `<td scope="row">%s</td> <td>%s</td> <td id="%s" onclick="bhash(this.id)">%s</td> <td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td>`)
	var yp, mp, dp, hp int
	count := 0
	tt := ``
	for a.Next() {
		var fromshard, toshard, id, txCount, stakingcount, signercount, size, timestamp, nonce, epoch int64
		var hash string
		e := a.Scan(&fromshard, &toshard, &id, &hash, &txCount, &stakingcount, &signercount, &size, &timestamp, &nonce, &epoch)
		if e != nil {
			fmt.Println(e)
			return
		}

		shards := fmt.Sprintf("%d.%d", fromshard, toshard)
		idx := ifmt(id)
		sz := sizefmt(size)
		y := time.Unix(timestamp, 0).UTC()
		yy := y.Year()
		mm := Months[y.Month().String()]
		dd := y.Day()
		yu := strings.Split(y.String(), " ")[1]
		hh := y.Hour()
		count++
		if yp == 0 {
			yp = yy
			mp = mm
			dp = dd
			hp = hh
		}

		if yy == yp {
			if mm == mp {
				if dd == dp {
					if hh == hp {
						tt += fmt.Sprintf(BlockRow1, shards, idx, hash, bhash(hash), txCount, stakingcount, signercount, sz, ifmt(nonce), ifmt(epoch), yu)
						continue
					}

					// new page
					_, e := db.Exec(`insert into cal (year, month, day, hour, bcount, arows) values (?,?,?,?,?,?)`, yy, mm, dd, hh, count, tt)
					if e != nil {
						fmt.Println(e)
						return
					}
					tt = ``
					tt += fmt.Sprintf(BlockRow1, shards, idx, hash, bhash(hash), txCount, stakingcount, signercount, sz, ifmt(nonce), ifmt(epoch), yu)
					hp = hh
					count = 1
					continue

				}

				// new day
				_, e := db.Exec(`insert into cal (year, month, day, hour, bcount, arows) values (?,?,?,?,?,?)`, yy, mm, dd, hh, count, tt)
				if e != nil {
					fmt.Println(e)
					return
				}
				tt = ``
				tt += fmt.Sprintf(BlockRow1, shards, idx, hash, bhash(hash), txCount, stakingcount, signercount, sz, ifmt(nonce), ifmt(epoch), yu)
				dp, hp = dd, hh
				count = 1
				continue

			}

			// new month
			_, e := db.Exec(`insert into cal (year, month, day, hour, bcount, arows) values (?,?,?,?,?,?)`, yy, mm, dd, hh, count, tt)
			if e != nil {
				fmt.Println(e)
				return
			}
			tt = ``
			tt += fmt.Sprintf(BlockRow1, shards, idx, hash, bhash(hash), txCount, stakingcount, signercount, sz, ifmt(nonce), ifmt(epoch), yu)
			mp, dp, hp = mm, dd, hh
			count = 1
			continue

		}

		// new year
		_, e = db.Exec(`insert into cal (year, month, day, hour, bcount, arows) values (?,?,?,?,?,?)`, yy, mm, dd, hh, count, tt)
		if e != nil {
			fmt.Println(e)
			return
		}
		tt = ``
		tt += fmt.Sprintf(BlockRow1, shards, idx, hash, bhash(hash), txCount, stakingcount, signercount, sz, ifmt(nonce), ifmt(epoch), yu)
		yp, mp, dp, hp = yy, mm, dd, hh
		count = 1
		continue

	}

	a.Close()

}
