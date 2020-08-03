package get

import (
	"log"
	"strconv"
	"sync"
	"time"
)

var w sync.WaitGroup

func Start(start, end int, id int) {
	s := sqlget(id)
	if s == 0 {
		_, err := db.Exec("INSERT INTO config VALUES (?,?)", id, start)
		if err != nil {
			panic(err)
		}
	}
	if s < start || s > end {
		_, err := db.Exec("UPDATE config SET i = ? WHERE id = ?", start, id)
		if err != nil {
			panic(err)
		}
	}
	for s < end {
		time.Sleep(500 * time.Millisecond)
		s = sqlget(id)
		b, err := getjson(strconv.Itoa(s))
		if err != nil {
			log.Println(err, "tid", s)
			time.Sleep(5 * time.Second)
			continue
		}
		t, err := json2(b)
		if err != nil {
			s++
			sqlup(s, id)
			continue
		}
		if ishide(t) {
			sqlset(t)
		} else if isqa(t) {
			qasave(t)
		}
		s++
		sqlup(s, id)
	}
	w.Done()
}

func Range(mintid, maxtid, thread int) {
	a := (maxtid - mintid) / thread
	w.Add(1)
	go Start(a*thread+mintid, maxtid+1, thread)
	for i := 0; i < thread; i++ {
		b := a * i
		if b == 0 {
			b++
		}
		w.Add(1)
		go Start(b+mintid, a*(i+1)+mintid, i)
	}
	w.Wait()
}
