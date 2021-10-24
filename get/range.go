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
		b, err := h(strconv.Itoa(tid))
		if err != nil {
			log.Println(err, "tid", tid)
			time.Sleep(3 * time.Second)
			continue
		}
		t, err := json2(b)
		if err != nil {
			log.Println(err, "tid", tid)
			return
		}
		if ishide(t) {
			sqlset(t)
		} else if isqa(t) {
			qasave(t)
		}
		break
	}
}

func Start() {
	sqlup(1270428, -1)
	last := sqlget(-1)
	if last == 0 {
		_, err := db.Exec("INSERT INTO config VALUES (?,?)", -1, 1)
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
	llast := 0
	for i := last; i < itid; i++ {
		if hasPost(i) {
			continue
		}
		w.Add(1)
		go start(i, &w)
		a++
		llast = i
		if a > 7 {
			w.Wait()
			a = 0
			sqlup(i, -1)
			time.Sleep(1 * time.Second)
		}
	}
	w.Wait()
	if llast != 0 {
		sqlup(llast, -1)
	}
}

const root = "https://late-sound-313b.xmdhs.workers.dev"

func getnewtid() (tid string, err error) {
	reqs, err := http.NewRequest("GET", root+"/forum.php?mod=guide&view=newthread&page=3", nil)
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
	if rep.StatusCode != 200 {
		return "", &ErrHttpCode{Code: rep.StatusCode}
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
