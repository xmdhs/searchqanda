package get

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var c = http.Client{Timeout: 5 * time.Second}

var cookie string

func init() {
	b, err := ioutil.ReadFile("cookie.txt")
	if err != nil {
		panic(err)
	}
	cookie = string(b)
}

func h(tid string) (b []byte, err error) {
	reqs, err := http.NewRequest("GET", `https://www.mcbbs.net/api/mobile/index.php?version=4&module=viewthread&tid=`+tid, nil)
	if err != nil {
		return
	}
	reqs.Header.Set("Accept", "*/*")
	reqs.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	reqs.Header.Set("Cookie", cookie)
	rep, err := c.Do(reqs)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return
	}
	if rep.StatusCode != 200 {
		err = errors.New(rep.Status)
		return
	}
	b, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		return
	}
	return
}

func getjson(tid string) (b []byte, err error) {
	for i := 0; i < 3; i++ {
		b, err = h(tid)
		if err != nil {
			log.Println(err, tid)
			continue
		}
		break
	}
	return
}

func ishide(t *thread) bool {
	b := false
	_, ok := t.Variables.Thread["tid"].(string)
	if ok {
		b = true
	}
	if len(t.Variables.Postlist) != 0 {
		b = false
	}
	return b
}

func isqa(t *thread) bool {
	fid, ok := t.Variables.Thread["fid"].(string)
	if ok && (fid == "265" || fid == "110" || fid == "431" || fid == "1566" || fid == "266") {
		return true
	}
	return false
}

type thread struct {
	Variables variables
}

type variables struct {
	Thread   map[string]interface{} `json:"thread"`
	Postlist []interface{}          `json:"postlist"`
}

func json2(b []byte) (t *thread, err error) {
	t = &thread{}
	t.Variables.Postlist = make([]interface{}, 0)
	t.Variables.Thread = make(map[string]interface{})
	err = json.Unmarshal(b, t)
	return
}
