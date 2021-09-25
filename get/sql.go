package get

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/yanyiwu/gojieba"
)

var db *sql.DB
var Db *sql.DB

var X = gojieba.NewJieba(`dict/jieba.dict.utf8`, `dict/hmm_model.utf8`, `dict/user.dict.utf8`, `dict/idf.utf8`, `dict/stop_words.utf8`)

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./hidethread.db?_txlock=IMMEDIATE&_journal_mode=WAL")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS hidethread(tid INT PRIMARY KEY NOT NULL,fid TEXT NOT NULL,authorid TEXT NOT NULL,author TEXT NOT NULL,views INT NOT NULL,dateline TEXT NOT NULL,lastpost TEXT NOT NULL,lastposter TEXT NOT NULL,subject TEXT NOT NULL)`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS config(id INT PRIMARY KEY NOT NULL,i INT NOT NULL)`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS qafts5 USING fts5(key,subject UNINDEXED, source)`)
	if err != nil {
		panic(err)
	}
	Db = db
}

func sqlset(t *thread) {
	tid := t.Variables.Thread["tid"].(string)
	fid := t.Variables.Thread["fid"].(string)
	authorid := t.Variables.Thread["authorid"].(string)
	author := t.Variables.Thread["author"].(string)
	views := t.Variables.Thread["views"].(string)
	dateline := t.Variables.Thread["dateline"].(string)
	lastpost := t.Variables.Thread["lastpost"].(string)
	lastposter := t.Variables.Thread["lastposter"].(string)
	subject := t.Variables.Thread["subject"].(string)

	subject = html.UnescapeString(html.UnescapeString(subject))

	i, err := strconv.ParseInt(dateline, 10, 64)
	if err != nil {
		panic(err)
	}
	l, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	dateline = time.Unix(i, 0).In(l).Format("2006-01-02 15:04:05")

	_, err = db.Exec(`INSERT INTO hidethread VALUES (?,?,?,?,?,?,?,?,?)`, tid, fid, authorid, author, views, dateline, lastpost, lastposter, subject)
	log.Println(tid, fid, authorid, author, views, dateline, lastpost, lastposter, subject)
	if err != nil {
		e := sqlite3.Error{}
		if errors.As(err, &e) {
			if e.Code == sqlite3.ErrBusy || e.Code == sqlite3.ErrLocked {
				log.Println(err)
				time.Sleep(1 * time.Second)
				sqlset(t)
				return
			}
		}
		panic(err)
	}
}

var htmlreg = regexp.MustCompile(`\<[\S\s]+?\>`)

func qasave(t *thread) {
	tid := t.Variables.Thread["tid"].(string)
	subject := t.Variables.Thread["subject"].(string)
	temptxt := t.Variables.Postlist
	tt := make([]post, 0, len(temptxt))
	log.Println(tid)
	for _, v := range temptxt {
		p := post{}
		m := v.(map[string]interface{})
		k, ok := m["message"].(string)
		if !ok {
			continue
		}
		k = htmlreg.ReplaceAllString(k, "")
		k = html.UnescapeString(html.UnescapeString(k))

		ks := X.CutForSearch(k, true)
		k = strings.Join(ks, "/")
		p.Message = k
		k, ok = m["authorid"].(string)
		p.Authorid = k
		if ok && k != "" {
			tt = append(tt, p)
		}
	}
	b, err := json.Marshal(tt)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO qafts5 VALUES (?,?,?)`, tid, subject, string(b))
	if err != nil {
		e := sqlite3.Error{}
		if errors.As(err, &e) {
			if e.Code == sqlite3.ErrBusy || e.Code == sqlite3.ErrLocked {
				log.Println(err)
				time.Sleep(1 * time.Second)
				qasave(t)
				return
			}
		}
		panic(err)
	}
}

type post struct {
	Message  string
	Authorid string
}

func sqlget(id int) int {
	row := db.QueryRow(`SELECT i FROM config WHERE id = ?`, id)
	var fid int
	err := row.Scan(&fid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	return fid
}

func sqlup(s, id int) {
	_, err := db.Exec("UPDATE config SET i = ? WHERE id = ?", s, id)
	if err != nil {
		e := sqlite3.Error{}
		if errors.As(err, &e) {
			if e.Code == sqlite3.ErrBusy || e.Code == sqlite3.ErrLocked {
				log.Println(err)
				time.Sleep(1 * time.Second)
				sqlup(s, id)
				return
			}
		}
		panic(err)
	}
}

func hasPost(tid int) bool {
	row := db.QueryRow(`SELECT COUNT(*) FROM qafts5 WHERE key MATCH ?`, tid)
	var c int
	err := row.Scan(&c)
	if err != nil {
		panic(err)
	}
	return c > 0
}
