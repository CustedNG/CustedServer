package logger

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/LollipopKit/custed-server/config"
	"github.com/LollipopKit/custed-server/consts"
)

func init() {
	if err := os.MkdirAll(consts.LogDir, consts.Permission); err != nil {
		panic(err)
	}
	go setup()
}

func W(format string, args ...any) {
	log.Printf("[WARN] "+format, args...)
}

func I(format string, args ...any) {
	log.Printf("[INFO] "+format, args...)
}

func E(format string, args ...any) {
	log.Printf("[ERROR] "+format, args...)
}

func D(format string, args ...any) {
	if !config.Debug {
		return
	}
	log.Printf("[DEBUG] "+format, args...)
}

// Must call this func using:
// `go setup()`
func setup() {
	for {
		file := consts.LogDir + time.Now().Format("2006-01-02") + ".txt"
		logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, consts.Permission)
		if err != nil {
			panic(err)
		}
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
		time.Sleep(time.Hour)
	}
}
