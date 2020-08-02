package get

import (
	"database/sql"

	//数据库驱动
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./hidethread.db")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(`SELECT json FROM hidethread`)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	var j string
	for rows.Next() {
		rows.Scan(&j)
		break
	}
	if j == "" {
		_, err := db.Exec(`CREATE TABLE hidethread(tid INT PRIMARY KEY NOT NULL,fid TEXT NOT NULL,authorid TEXT NOT NULL,author TEXT NOT NULL,views INT NOT NULL,lastpost TEXT NOT NULL,lastposter TEXT NOT NULL)`)
		if err != nil {
			panic(err)
		}
	}
}

func sqlset(t *thread) {
	stmt, err := db.Prepare(`INSERT INTO hidethread VALUES (?,?,?,?,?,?,?)`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	tid := t.Variables.Thread["tid"].(string)
	fid := t.Variables.Thread["fid"].(string)
	authorid := t.Variables.Thread["authorid"].(string)
	author := t.Variables.Thread["author"].(string)
	views := t.Variables.Thread["views"].(string)
	lastpost := t.Variables.Thread["lastpost"].(string)
	lastposter := t.Variables.Thread["lastposter"].(string)
	_, err = stmt.Exec(tid, fid, authorid, author, views, lastpost, lastposter)
	if err != nil {
		panic(err)
	}
}

func sqlget(id int) int {
	stmt, err := db.Prepare(`SELECT fid FROM hidethread WHERE tid = ?`)
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
	stmt, err := db.Prepare("UPDATE hidethread SET fid = ? WHERE tid = ?")
	if err != nil {
		panic(err)
	}
	stmt.Exec(s, id)
}
