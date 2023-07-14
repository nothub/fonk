package main

import (
	"syscall"
)

func setUidGid(uid int, gid int) {
	err := syscall.Setgroups([]int{})
	if err != nil {
		elog.Fatalf("setgroups() failure (%d)", err)
	}

	err = syscall.Setgid(gid)
	if err != nil {
		elog.Fatalf("Setgid(%s) failure (%d)", gid, err)
	}

	err = syscall.Setuid(uid)
	if err != nil {
		elog.Fatalf("Setuid(%s) failure (%d)", uid, err)
	}
}

func init() {
	preservehooks = append(preservehooks, func() {
		setUidGid(10042, 10042)
	})
}
