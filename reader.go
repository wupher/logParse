package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
)

func getFile(fileName string) (f *os.File, isGzip bool, err error) {
	f, err = os.Open(fileName)
	if strings.HasSuffix(fileName, ".gz") {
		isGzip = true
	}
	return
}

func process(f *os.File, isGzip bool) error {
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 25*1024)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	var reader *bufio.Reader

	if isGzip {
		fz, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		reader = bufio.NewReader(fz)
	} else {
		reader = bufio.NewReader(f)
	}

	var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)
		n, err := reader.Read(buf)
		buf = buf[:n]
		if n == 0 {
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				break
			}
			return err
		}
		nextUntilNewline, err := reader.ReadBytes('\n')

		if err != io.EOF {
			buf = append(buf, nextUntilNewline...)
		}
		wg.Add(1)
		go func() {
			processChunk(buf, &linesPool, &stringPool)
			wg.Done()
		}()

	}
	wg.Wait()
	return nil
}

func processChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool) {
	var wg2 sync.WaitGroup

	logs := stringPool.Get().(string)
	logs = string(chunk)

	linesPool.Put(chunk)
	logsSlice := strings.Split(logs, "\n")

	stringPool.Put(logs)
	chunkSize := 300
	n := len(logsSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}

	for i := 0; i < noOfThread; i++ {
		wg2.Add(1)

		go func(s int, e int) {
			defer wg2.Done()
			for i := s; i < e; i++ {
				text := logsSlice[i]
				if len(text) == 0 {
					continue
				}
				lineProcess(text)
			}
		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
	}
	wg2.Wait()
	logsSlice = nil
}
