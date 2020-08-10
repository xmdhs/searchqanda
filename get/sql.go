package get

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"time"

	//数据库驱动
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var Db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./hidethread.db")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS hidethread(tid INT PRIMARY KEY NOT NULL,fid TEXT NOT NULL,authorid TEXT NOT NULL,author TEXT NOT NULL,views INT NOT NULL,dateline TEXT NOT NULL,lastpost TEXT NOT NULL,lastposter TEXT NOT NULL,subject TEXT NOT NULL)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS qa(tid INT PRIMARY KEY NOT NULL,fid TEXT NOT NULL,subject TEXT NOT NULL,txt TEXT NOT NULL)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS config(id INT PRIMARY KEY NOT NULL,i INT NOT NULL)`)
	_, err = db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS idx USING fts5(key, source, tokenize=icu)`)
	if err != nil {
		log.Println(err)
	}
	Db = db
}

func sqlset(t *thread) {
	stmt, err := db.Prepare(`INSERT INTO hidethread VALUES (?,?,?,?,?,?,?,?,?)`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	tid := t.Variables.Thread["tid"].(string)
	fid := t.Variables.Thread["fid"].(string)
	authorid := t.Variables.Thread["authorid"].(string)
	author := t.Variables.Thread["author"].(string)
	views := t.Variables.Thread["views"].(string)
	dateline := t.Variables.Thread["dateline"].(string)
	lastpost := t.Variables.Thread["lastpost"].(string)
	lastposter := t.Variables.Thread["lastposter"].(string)
	subject := t.Variables.Thread["subject"].(string)

	i, err := strconv.ParseInt(dateline, 10, 64)
	if err != nil {
		panic(err)
	}
	dateline = time.Unix(i, 0).Format("2006-01-02 15:04:05")

	_, err = stmt.Exec(tid, fid, authorid, author, views, dateline, lastpost, lastposter, subject)
	log.Println(tid, fid, authorid, author, views, dateline, lastpost, lastposter, subject)
	if err != nil {
		log.Println(err, t)
	}
}

func qasave(t *thread) {
	stmt, err := db.Prepare(`INSERT INTO qa VALUES (?,?,?,?)`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	tid := t.Variables.Thread["tid"].(string)
	fid := t.Variables.Thread["fid"].(string)
	subject := t.Variables.Thread["subject"].(string)
	temptxt := t.Variables.Postlist
	tt := make([]post, 0, len(temptxt))
	for _, v := range temptxt {
		p := post{}
		m := v.(map[string]interface{})
		k, ok := m["message"].(string)
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
	_, err = stmt.Exec(tid, fid, subject, string(b))
	if err != nil {
		log.Println(err, t)
	}

}

type post struct {
	Message  string
	Authorid string
}

var Sqlget = sqlget

func sqlget(id int) int {
	stmt, err := db.Prepare(`SELECT i FROM config WHERE id = ?`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query(id)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	rows.Next()
	var fid int
	rows.Scan(&fid)
	return fid
}

func sqlup(s, id int) {
	stmt, err := db.Prepare("UPDATE config SET i = ? WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	stmt.Exec(s, id)
}
