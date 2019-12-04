package system

import (
	"syscall"

	log "github.com/sirupsen/logrus"
)

func GetLimits() (int64, int64) {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		log.Fatalf("Getrlimit failed with: %v", err)
	}

	return int64(limit.Cur), int64(limit.Max)
}
