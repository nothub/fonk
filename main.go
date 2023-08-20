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
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	golog "log"
	"log/syslog"
	notrand "math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"humungus.tedunangst.com/r/webs/log"
)

var softwareVersion = "develop"

func init() {
	notrand.Seed(time.Now().Unix())
}

var serverName string
var serverPrefix string
var masqName string
var dataDir = "./data"
var viewDir = "./views"
var iconName = "icon.png"
var serverMsg template.HTML
var aboutMsg template.HTML
var loginMsg template.HTML

func ElaborateUnitTests() {
}

func unplugserver(hostname string) {
	db := opendatabase()
	xid := fmt.Sprintf("%%https://%s/%%", hostname)
	db.Exec("delete from honkers where xid like ? and flavor = 'dub'", xid)
	db.Exec("delete from doovers where rcpt like ?", xid)
}

func reexecArgs(cmd string) []string {
	args := []string{"-datadir", dataDir}
	args = append(args, log.Args()...)
	args = append(args, cmd)
	return args
}

var elog, ilog, dlog *golog.Logger

func main() {
	log.Init(log.Options{Progname: "honk", Facility: syslog.LOG_UUCP})
	elog = log.E
	ilog = log.I
	dlog = log.D

	if os.Geteuid() == 0 {
		elog.Fatalf("do not run honk as root")
	}

	flag.StringVar(&dataDir, "datadir", dataDir, "data directory")
	flag.StringVar(&viewDir, "viewdir", viewDir, "view directory")
	flag.Parse()

	args := flag.Args()
	cmd := "run"
	if len(args) > 0 {
		cmd = args[0]
		args = args[1:]
	}

	switch cmd {
	case "genhash":
		os.Stderr.WriteString("This command generates bcrypt hashes for use with the init command.\n")
		os.Stderr.WriteString("To store the hash string output to a shell variable, do this:\n")
		os.Stderr.WriteString("hash=\"$(./honk genhash)\"\n")
		pass, err := askpassword()
		if err != nil {
			elog.Fatalf("error: %s\n", err.Error())
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
		if err != nil {
			elog.Fatalf("error: %s\n", err.Error())
		}
		os.Stderr.WriteString("For a flag option, copy the following line:\n")
		os.Stderr.WriteString("--hash \"" + strings.ReplaceAll(string(hash), "$", "\\$") + "\"\n")
		fmt.Println(string(hash))
		os.Exit(0)
	case "init":
		flags := flag.NewFlagSet(cmd, flag.ExitOnError)
		var hash string
		var listen string
		flags.StringVar(&hash, "hash", "", "password hash (bcrypt)")
		flags.StringVar(&listen, "listen", "0.0.0.0:8080", "listen address")
		err := flags.Parse(args)
		if err != nil {
			elog.Fatalf("failed parsing flags: %s\n", err.Error())
		}
		args = flags.Args()
		if len(args) < 2 {
			fmt.Printf("usage: init [--hash <hash>] [--listen <host:port>] <username> <fqdn>\n")
			return
		}
		initdb(args[0], hash, args[1], listen)
	case "upgrade":
		upgradedb()
	case "version":
		fmt.Println(softwareVersion)
		os.Exit(0)
	}

	db := opendatabase()
	dbversion := 0
	getconfig("dbversion", &dbversion)
	if dbversion != myVersion {
		elog.Fatal("incorrect database version. run upgrade.")
	}
	getconfig("servermsg", &serverMsg)
	getconfig("aboutmsg", &aboutMsg)
	getconfig("loginmsg", &loginMsg)
	getconfig("servername", &serverName)
	getconfig("masqname", &masqName)
	if masqName == "" {
		masqName = serverName
	}
	serverPrefix = fmt.Sprintf("https://%s/", serverName)
	getconfig("usersep", &userSep)
	getconfig("honksep", &honkSep)
	getconfig("devel", &develMode)
	getconfig("fasttimeout", &fastTimeout)
	getconfig("slowtimeout", &slowTimeout)
	getconfig("signgets", &signGets)
	prepareStatements(db)

	switch cmd {
	case "admin":
		adminscreen()
	case "import":
		if len(args) != 3 {
			elog.Fatal("usage: import <username> (honk|mastodon|twitter) <srcdir>\n")
		}
		importMain(args[0], args[1], args[2])
	case "export":
		if len(args) != 2 {
			elog.Fatal("export username destdir")
		}
		export(args[0], args[1])
	case "devel":
		if len(args) != 1 {
			elog.Fatal("usage: devel (on|off)")
		}
		switch args[0] {
		case "on":
			setconfig("devel", 1)
		case "off":
			setconfig("devel", 0)
		default:
			elog.Fatal("usage: devel (on|off)")
		}
	case "setconfig":
		if len(args) != 2 {
			elog.Fatal("usage: setconfig <key> <val>\n")
		}
		var val interface{}
		var err error
		if val, err = strconv.Atoi(args[1]); err != nil {
			val = args[1]
		}
		setconfig(args[0], val)
	case "adduser":
		flags := flag.NewFlagSet(cmd, flag.ExitOnError)
		var hash string
		flags.StringVar(&hash, "hash", "", "password hash (bcrypt)")
		err := flags.Parse(args)
		if err != nil {
			elog.Fatalf("failed parsing flags: %s\n", err.Error())
		}
		args = flags.Args()
		if len(args) < 1 {
			fmt.Printf("usage: honk adduser [--hash <hash>] <username>\n")
			return
		}
		adduser(args[0], hash)
	case "deluser":
		if len(args) < 1 {
			fmt.Printf("usage: honk deluser <username>\n")
			return
		}
		deluser(args[0])
	case "chpass":
		if len(args) < 1 {
			fmt.Printf("usage: honk chpass <username>\n")
			return
		}
		chpass(args[0])
	case "follow":
		if len(args) < 2 {
			fmt.Printf("usage: honk follow <username> <url>\n")
			return
		}
		user, err := butwhatabout(args[0])
		if err != nil {
			fmt.Printf("user not found\n")
			return
		}
		var meta HonkerMeta
		mj, _ := jsonify(&meta)
		honkerid, err := savehonker(user, args[1], "", "presub", "", mj)
		if err != nil {
			fmt.Printf("had some trouble with that: %s\n", err)
			return
		}
		followyou(user, honkerid, true)
	case "unfollow":
		if len(args) < 2 {
			fmt.Printf("usage: honk unfollow username url\n")
			return
		}
		user, err := butwhatabout(args[0])
		if err != nil {
			fmt.Printf("user not found\n")
			return
		}
		row := db.QueryRow("select honkerid from honkers where xid = ? and userid = ? and flavor in ('sub')", args[1], user.ID)
		var honkerid int64
		err = row.Scan(&honkerid)
		if err != nil {
			fmt.Printf("sorry couldn't find them\n")
			return
		}
		unfollowyou(user, honkerid, true)
	case "sendmsg":
		if len(args) < 3 {
			fmt.Printf("usage: honk sendmsg username filename rcpt\n")
			return
		}
		user, err := butwhatabout(args[0])
		if err != nil {
			fmt.Printf("user not found\n")
			return
		}
		data, err := os.ReadFile(args[1])
		if err != nil {
			fmt.Printf("can't read file\n")
			return
		}
		deliverate(user.ID, args[2], data)
	case "cleanup":
		arg := "30"
		if len(args) > 1 {
			arg = args[1]
		}
		cleanupdb(arg)
	case "unplug":
		if len(args) < 1 {
			fmt.Printf("usage: honk unplug <servername>\n")
			return
		}
		name := args[0]
		unplugserver(name)
	case "backup":
		if len(args) < 1 {
			fmt.Printf("usage: honk backup <dirname>\n")
			return
		}
		name := args[0]
		svalbard(name)
	case "ping":
		if len(args) < 2 {
			fmt.Printf("usage: honk ping (from username) (to username or url)\n")
			return
		}
		name := args[0]
		targ := args[1]
		user, err := butwhatabout(name)
		if err != nil {
			elog.Printf("unknown user")
			return
		}
		ping(user, targ)
	case "run":
		serve()
	case "backend":
		backendServer()
	case "test":
		ElaborateUnitTests()
	default:
		elog.Fatal("unknown command")
	}
}
