package get

import (
	"bufio"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var w sync.WaitGroup
var M sync.RWMutex

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
		M.RLock()
		s = sqlget(id)
		M.RUnlock()
		b, err := getjson(strconv.Itoa(s))
		if err != nil {
			log.Println(err, "tid", s)
			time.Sleep(5 * time.Second)
			continue
		}
		t, err := json2(b)
		M.Lock()
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
		M.Unlock()
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

func Startrange() {
	start := sqlget(-2)
	end := sqlget(-1)
	if start == 0 {
		_, err := db.Exec("INSERT INTO config VALUES (?,?)", -2, 0)
		if err != nil {
			log.Println(err)
		}
	}
	if end == 0 {
		_, err := db.Exec("INSERT INTO config VALUES (?,?)", -1, 1092244)
		if err != nil {
			panic(err)
		}
		end = 1092244
	}
	Range(start, end, 5)
	tid, err := getnewtid()
	if err != nil {
		log.Println(err)
		return
	}
	if tid == "" {
		log.Println(`tid == ""`)
		return
	}

	_, err = db.Exec("UPDATE config SET i = ? WHERE id = ?", end, -2)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("UPDATE config SET i = ? WHERE id = ?", tid, -1)
	if err != nil {
		panic(err)
	}
}

func getnewtid() (tid string, err error) {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	reqs, err := http.NewRequest("GET", "https://www.mcbbs.net/forum.php?mod=guide&view=newthread&page=2", nil)
	if err != nil {
		return
	}
	reqs.Header.Set("Accept", "*/*")
	reqs.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	rep, err := c.Do(reqs)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return
	}
	bw := bufio.NewScanner(rep.Body)
	for bw.Scan() {
		if strings.Contains(bw.Text(), `target="_blank" class="xst" `) {
			t := bw.Text()
			a := strings.Index(t, `<a href="thread-`) + 16
			b := strings.Index(t[a:], `-1-1.html" target`)
			tid = t[a:][:b]
			return
		}
	}
	return
}
