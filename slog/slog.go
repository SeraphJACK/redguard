package slog

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const LogPath = "./slog.log"

var ch = make(chan Entry, 10)

type Entry struct {
	IP        string    `json:"ip"`
	Operation string    `json:"operation"`
	Content   string    `json:"content"`
	Result    string    `json:"result"`
	Time      time.Time `json:"time"`
}

func Log(ip, op, content, res string) {
	e := Entry{
		IP:        ip,
		Operation: op,
		Content:   content,
		Result:    res,
		Time:      time.Now(),
	}

	ch <- e
}

func init() {
	f, err := os.OpenFile(LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			entry := <-ch
			b, _ := json.Marshal(entry)
			_, err = f.Write(append(b, []byte("\n")...))
			if err != nil {
				log.Println("Failed to write log: ", err)
			}
			_ = f.Sync()
		}
	}()
}
