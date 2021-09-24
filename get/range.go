package get

import (
	"bufio"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func start(tid int, w *sync.WaitGroup) {
	defer w.Done()
	for {
		b, err := getjson(strconv.Itoa(tid))
		if err != nil {
			log.Println(err, "tid", tid)
			time.Sleep(3 * time.Second)
			continue
		}
		t, err := json2(b)
		if err != nil {
			log.Println(err, "tid", tid)
			time.Sleep(3 * time.Second)
			continue
		}
		if ishide(t) {
			sqlset(t)
		} else if isqa(t) {
			qasave(t)
		}
	}
}

func Start() {
	last := sqlget(-2)
	if last == 0 {
		_, err := db.Exec("INSERT INTO config VALUES (?,?)", -2, 1)
		if err != nil {
			panic(err)
		}
		last = 1
	}

	tid, err := getnewtid()
	if err != nil {
		panic(err)
	}
	if tid == "" {
		panic(`tid == ""`)
	}
	w := sync.WaitGroup{}

	itid, err := strconv.Atoi(tid)
	if err != nil {
		panic(err)
	}

	a := 0
	for i := last; i < itid; i++ {
		w.Add(1)
		go start(i, &w)
		a++
		if a > 7 {
			w.Wait()
			a = 0
			sqlup(i, -2)
			time.Sleep(1 * time.Second)
		}
	}

}

func getnewtid() (tid string, err error) {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	reqs, err := http.NewRequest("GET", "https://late-sound-313b.xmdhs.workers.dev/forum.php?mod=guide&view=newthread&page=3", nil)
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
		txt := bw.Text()
		if strings.Contains(txt, `target="_blank" class="xst" `) {
			tid, err = cut(txt, `<a href="thread-`, `-1-1.html" target`)
			if err != nil {
				log.Println(err)
				continue
			}
			return
		}
	}
	return
}

func cut(text, start, end string) (string, error) {
	a := strings.Index(text, start)
	if a == -1 {
		return "", ErrNotFound
	}
	temp := text[a+len(start):]
	b := strings.Index(temp, end)
	if b == -1 {
		return "", ErrNotFound
	}
	return temp[:b], nil
}

var ErrNotFound = errors.New(`not found`)
