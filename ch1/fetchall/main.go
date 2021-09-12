package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
	f, err := os.Create("tmp")
	if err != nil {
		fmt.Print(err)
		return
	}
	w := bufio.NewWriter(f)
	for range os.Args[1:] {
		_, err = fmt.Fprintln(w, <-ch)
		if err != nil {
			fmt.Print(err)
			return
		}
	}

	fmt.Fprintf(w, "%.2fs elapsed\n", time.Since(start).Seconds())
	w.Flush()
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
