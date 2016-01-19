package pump

import (
	"os"

	"github.com/hpcloud/tail"
)

func TailFile(cs chan string, filename string) {
	file, _ := tail.TailFile(filename, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_CUR}})
	for line := range file.Lines {
		cs <- string(line.Text)
	}
}
