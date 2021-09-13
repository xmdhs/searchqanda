package get

import (
	"database/sql"
	"encoding/json"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yanyiwu/gojieba"
	//数据库驱动
	_ "github.com/mattn/go-sqlite3"
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
		log.Println(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS config(id INT PRIMARY KEY NOT NULL,i INT NOT NULL)`)
	if err != nil {
		log.Println(err)
	}
	_, err = db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS qafts5 USING fts5(key,subject UNINDEXED, source)`)
	if err != nil {
		log.Println(err)
	}
	Db = db
}

func sqlset(t *thread) {
	stmt, err := db.Prepare(`INSERT INTO hidethread VALUES (?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
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

	_, err = stmt.Exec(tid, fid, authorid, author, views, dateline, lastpost, lastposter, subject)
	log.Println(tid, fid, authorid, author, views, dateline, lastpost, lastposter, subject)
	if err != nil {
		log.Println(err, t)
	}
}

var htmlreg = regexp.MustCompile(`\<[\S\s]+?\>`)

func qasave(t *thread) {
	stmt, err := db.Prepare(`INSERT INTO qafts5 VALUES (?,?,?)`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	tid := t.Variables.Thread["tid"].(string)
	subject := t.Variables.Thread["subject"].(string)
	temptxt := t.Variables.Postlist
	tt := make([]post, 0, len(temptxt))
	log.Println("start", tid)
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
	_, err = stmt.Exec(tid, subject, string(b))
	log.Println("end", tid)
	if err != nil {
		log.Println(err, t)
	}

}

type post struct {
	Message  string
	Authorid string
}

func sqlget(id int) int {
	stmt, err := db.Prepare(`SELECT i FROM config WHERE id = ?`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	rows.Next()
	var fid int
	rows.Scan(&fid)
	return fid
}

func sqlup(s, id int) {
	stmt, err := db.Prepare("UPDATE config SET i = ? WHERE id = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	stmt.Exec(s, id)
}
