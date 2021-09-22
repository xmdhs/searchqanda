package get

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	return httpget(`https://late-sound-313b.xmdhs.workers.dev/api/mobile/index.php?version=4&module=viewthread&tid=`+tid, cookie)
}

func httpget(url string, cookie string) ([]byte, error) {
	reqs, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpget: %w", err)
	}
	reqs.Header.Set("Accept", "*/*")
	reqs.Header.Add("accept-encoding", "gzip")
	reqs.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
	reqs.Header.Set("Cookie", cookie)
	rep, err := c.Do(reqs)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("httpget: %w", err)
	}
	if rep.StatusCode != 200 {
		return nil, fmt.Errorf("httpget: %w", &ErrHttpCode{Code: rep.StatusCode})
	}
	var reader io.ReadCloser
	switch rep.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(rep.Body)
		if err != nil {
			return nil, fmt.Errorf("httpget: %w", err)
		}
		defer reader.Close()
	default:
		reader = rep.Body
	}
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("httpget: %w", err)
	}
	return b, nil
}

type ErrHttpCode struct {
	Code int
}

func (e *ErrHttpCode) Error() string {
	return "http code: " + strconv.Itoa(e.Code)
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
