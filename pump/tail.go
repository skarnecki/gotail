package pump

import (
	"log"
	"os"

	"github.com/hpcloud/tail"
)

const maxBytesLoaded = 10000

func TailFile(cs chan string, filename string) {
	cfg := tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: -maxBytesLoaded, Whence: os.SEEK_END}}
	if getFileSize(filename) < maxBytesLoaded {
		cfg.Location.Whence = os.SEEK_CUR
		cfg.Location.Offset = 0
	}

	fileTail, _ := tail.TailFile(filename, cfg)
	for line := range fileTail.Lines {
		if line.Text != "" {
			cs <- string(line.Text)
		}
	}
}

func getFileSize(filename string) int64 {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return fi.Size()
}
