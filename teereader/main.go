package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

func main() {
	out, err := os.Create("test.tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	resp, err := http.Get("https://dl.google.com/go/go1.17.1.src.tar.gz")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	teeReader := &Speeder{TotalSize: resp.ContentLength}
	// 打印读取速度
	go teeReader.Show()
	io.Copy(out, io.TeeReader(resp.Body, teeReader))
}

// Speeder 用于记录时间段内读取的字节数
type Speeder struct {
	count     int64
	size      int64
	TotalSize int64
}

// Write 实现Writer接口，记录读取的字节数
func (s *Speeder) Write(b []byte) (int, error) {
	c := len(b)
	atomic.AddInt64(&s.count, int64(c))
	atomic.AddInt64(&s.size, int64(c))
	fmt.Printf("\r%.2f\n", float64(atomic.LoadInt64(&s.size))/float64(atomic.LoadInt64(&s.TotalSize))*100)
	return c, nil
}

// Show 打印读取速度
func (s *Speeder) Show() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		fmt.Printf("\r%.2fkb/s\n", float64(atomic.LoadInt64(&s.count))/float64(1024))
		// fmt.Printf("\r%.2f", float64(atomic.LoadInt64(&s.size)))
		atomic.StoreInt64(&s.count, 0)
	}
}
