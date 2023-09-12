//
// Copyright (c) 2019 Ted Unangst <tedu@tedunangst.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"golang.org/x/term"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"humungus.tedunangst.com/r/webs/httpsig"
	"humungus.tedunangst.com/r/webs/login"
)

var re_plainname = regexp.MustCompile("^[[:alnum:]_-]+$")

var dbtimeformat = "2006-01-02 15:04:05"

var alreadyopendb *sql.DB
var stmtConfig *sql.Stmt

func initdb(username string, hash string, hostname string, listen string) {
	dbname := dataDir + "/honk.db"
	_, err := os.Stat(dbname)
	if err == nil {
		elog.Fatalf("%s already exists", dbname)
	}
	db, err := sql.Open(sqlDriver, dbname)
	if err != nil {
		elog.Fatal(err)
	}
	alreadyopendb = db
	defer func() {
		os.Remove(dbname)
		os.Exit(1)
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Remove(dbname)
		os.Exit(1)
	}()

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		elog.Print(err)
		return
	}
	for _, line := range strings.Split(sqlSchema, ";") {
		_, err = db.Exec(line)
		if err != nil {
			elog.Print(err)
			return
		}
	}

	initblobdb()

	prepareStatements(db)

	err = createuser(db, username, hash)
	if err != nil {
		elog.Fatalf("error: %s\n", err.Error())
	}

	// must came later or user above will have negative id
	err = createserveruser(db)
	if err != nil {
		elog.Print(err)
		return
	}

	// TODO: better validation
	if len(hostname) < 1 {
		elog.Print("that's way too short")
		return
	}
	setconfig("servername", hostname)

	// TODO: better validation
	if len(listen) < 1 {
		elog.Print("that's way too short")
		return
	}
	setconfig("listenaddr", listen)

	var randbytes [16]byte
	rand.Read(randbytes[:])
	key := fmt.Sprintf("%x", randbytes)
	setconfig("csrfkey", key)
	setconfig("dbversion", myVersion)

	setconfig("servermsg", "<h2>Things happen.</h2>")
	setconfig("aboutmsg", "<h3>What is honk?</h3><p>Honk is amazing!</p>")
	setconfig("loginmsg", "<h2>login</h2>")
	setconfig("devel", 0)

	db.Close()
	fmt.Printf("done.\n")
	os.Exit(0)
}

func initblobdb() {
	blobdbname := dataDir + "/blob.db"
	_, err := os.Stat(blobdbname)
	if err == nil {
		elog.Fatalf("%s already exists", blobdbname)
	}
	blobdb, err := sql.Open(sqlDriver, blobdbname)
	if err != nil {
		elog.Print(err)
		return
	}
	_, err = blobdb.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		elog.Print(err)
		return
	}
	_, err = blobdb.Exec("create table filedata (xid text, media text, hash text, content blob)")
	if err != nil {
		elog.Print(err)
		return
	}
	_, err = blobdb.Exec("create index idx_filexid on filedata(xid)")
	if err != nil {
		elog.Print(err)
		return
	}
	_, err = blobdb.Exec("create index idx_filehash on filedata(hash)")
	if err != nil {
		elog.Print(err)
		return
	}
	blobdb.Close()
}

func adduser(username string, hash string) {
	db := opendatabase()
	err := createuser(db, username, hash)
	if err != nil {
		elog.Fatalf("error: %s\n", err.Error())
	}
	os.Exit(0)
}

func deluser(username string) {
	user, _ := butwhatabout(username)
	if user == nil {
		elog.Printf("no userfound")
		return
	}
	userid := user.ID
	db := opendatabase()

	where := " where honkid in (select honkid from honks where userid = ?)"
	doordie(db, "delete from donks"+where, userid)
	doordie(db, "delete from onts"+where, userid)
	doordie(db, "delete from honkmeta"+where, userid)
	where = " where chonkid in (select chonkid from chonks where userid = ?)"
	doordie(db, "delete from donks"+where, userid)

	doordie(db, "delete from honks where userid = ?", userid)
	doordie(db, "delete from chonks where userid = ?", userid)
	doordie(db, "delete from honkers where userid = ?", userid)
	doordie(db, "delete from zonkers where userid = ?", userid)
	doordie(db, "delete from doovers where userid = ?", userid)
	doordie(db, "delete from hfcs where userid = ?", userid)
	doordie(db, "delete from auth where userid = ?", userid)
	doordie(db, "delete from users where userid = ?", userid)
}

func chpass(username string) {
	user, err := butwhatabout(username)
	if err != nil {
		elog.Fatalf("error: %s\n", err.Error())
	}

	pass, err := askpassword()
	if err != nil {
		elog.Fatalf("error: %s\n", err.Error())
	}

	db := opendatabase()
	login.Init(login.InitArgs{Db: db, Logger: ilog})

	err = login.SetPassword(user.ID, pass)
	if err != nil {
		elog.Fatalf("error: %s\n", err.Error())
	}

	fmt.Printf("done\n")
	os.Exit(0)
}

func askpassword() (string, error) {
	os.Stderr.WriteString("enter password: ")
	pass, err := term.ReadPassword(syscall.Stdin)
	os.Stderr.WriteString("\n")
	if err != nil {
		return "", err
	}

	os.Stderr.WriteString("enter again: ")
	repeat, err := term.ReadPassword(syscall.Stdin)
	os.Stderr.WriteString("\n")
	if err != nil {
		return "", err
	}

	if string(pass) != string(repeat) {
		return "", fmt.Errorf("inputs were not the same")
	}
	if len(pass) < 8 {
		return "", fmt.Errorf("that's way too short")
	}

	return string(pass), nil
}

func createuser(db *sql.DB, name string, hash string) error {
	if len(name) < 1 {
		return fmt.Errorf("that's way too short")
	}
	if !re_plainname.MatchString(name) {
		return fmt.Errorf("alphanumeric only please")
	}
	if _, err := butwhatabout(name); err == nil {
		return fmt.Errorf("user already exists")
	}

	if len(hash) == 0 {
		pass, err := askpassword()
		if err != nil {
			elog.Fatalf("error: %s\n", err.Error())
		}
		b, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
		if err != nil {
			return err
		}
		hash = string(b)
	}

	// validate hash
	if _, err := bcrypt.Cost([]byte(hash)); err != nil {
		return fmt.Errorf("invalid bcrypt hash")
	}

	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	pubkey, err := httpsig.EncodeKey(&k.PublicKey)
	if err != nil {
		return err
	}
	seckey, err := httpsig.EncodeKey(k)
	if err != nil {
		return err
	}
	about := "what about me?"
	_, err = db.Exec("insert into users (username, displayname, about, hash, pubkey, seckey, options) values (?, ?, ?, ?, ?, ?, ?)", name, name, about, hash, pubkey, seckey, "{}")
	if err != nil {
		return err
	}

	return nil
}

func createserveruser(db *sql.DB) error {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	pubkey, err := httpsig.EncodeKey(&k.PublicKey)
	if err != nil {
		return err
	}
	seckey, err := httpsig.EncodeKey(k)
	if err != nil {
		return err
	}
	name := "server"
	about := "server"
	hash := "*"
	_, err = db.Exec("insert into users (userid, username, displayname, about, hash, pubkey, seckey, options) values (?, ?, ?, ?, ?, ?, ?, ?)", serverUID, name, name, about, hash, pubkey, seckey, "")
	if err != nil {
		return err
	}
	return nil
}

func opendatabase() *sql.DB {
	if alreadyopendb != nil {
		return alreadyopendb
	}
	dbname := dataDir + "/honk.db"
	_, err := os.Stat(dbname)
	if err != nil {
		elog.Fatalf("unable to open database: %s", err)
	}
	db, err := sql.Open(sqlDriver, dbname)
	if err != nil {
		elog.Fatalf("unable to open database: %s", err)
	}
	stmtConfig, err = db.Prepare("select value from config where key = ?")
	if err != nil {
		elog.Fatal(err)
	}
	alreadyopendb = db
	return db
}

func openblobdb() *sql.DB {
	blobdbname := dataDir + "/blob.db"
	_, err := os.Stat(blobdbname)
	if err != nil {
		elog.Fatalf("unable to open database: %s", err)
	}
	db, err := sql.Open(sqlDriver, blobdbname)
	if err != nil {
		elog.Fatalf("unable to open database: %s", err)
	}
	return db
}

func getconfig(key string, value interface{}) error {
	m, ok := value.(*map[string]bool)
	if ok {
		rows, err := stmtConfig.Query(key)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var s string
			err = rows.Scan(&s)
			if err != nil {
				return err
			}
			(*m)[s] = true
		}
		return nil
	}
	row := stmtConfig.QueryRow(key)
	err := row.Scan(value)
	if err == sql.ErrNoRows {
		err = nil
	}
	return err
}

func setconfig(key string, val interface{}) error {
	db := opendatabase()
	db.Exec("delete from config where key = ?", key)
	_, err := db.Exec("insert into config (key, value) values (?, ?)", key, val)
	return err
}

func openListener() (net.Listener, error) {
	var listenAddr string
	err := getconfig("listenaddr", &listenAddr)
	if err != nil {
		return nil, err
	}
	if listenAddr == "" {
		return nil, fmt.Errorf("must have listenaddr")
	}
	proto := "tcp"
	if listenAddr[0] == '/' {
		proto = "unix"
		err := os.Remove(listenAddr)
		if err != nil && !os.IsNotExist(err) {
			elog.Printf("unable to unlink socket: %s", err)
		}
	}
	listener, err := net.Listen(proto, listenAddr)
	if err != nil {
		return nil, err
	}
	if proto == "unix" {
		os.Chmod(listenAddr, 0777)
	}
	return listener, nil
}
