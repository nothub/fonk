package main

import (
	"syscall"
)

var puid = -1
var pgid = -1

func setUidGid(uid int, gid int) {
	dlog.Printf("getuid() returned %v", syscall.Getuid())
	dlog.Printf("getgid() returned %v", syscall.Getgid())

	if err := syscall.Setgroups([]int{}); err != nil {
		elog.Fatalf("setgroups() failure (%d)", err)
	}

	// TODO: verify this also sets E and R ids

	if err := syscall.Setgid(gid); err != nil {
		elog.Fatalf("Setgid(%v) failure (%d)", gid, err)
	}

	if err := syscall.Setuid(uid); err != nil {
		elog.Fatalf("Setuid(%v) failure (%d)", uid, err)
	}

	dlog.Printf("getuid() returned %v", syscall.Getuid())
	dlog.Printf("getgid() returned %v", syscall.Getgid())
}

func init() {
	preservehooks = append(preservehooks, func() {
		if puid > 0 && pgid > 0 {
			setUidGid(puid, pgid)
		}
	})
}
