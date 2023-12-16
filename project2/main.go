package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

func main() {
	infile, err := os.Open("test.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer infile.Close()
	outfile, err := os.OpenFile("out.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outfile.Close()

	cr := csv.NewReader(infile)
	cw := csv.NewWriter(outfile)

	ch := make(chan *[][]string, 10)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		process(cw, ch)
		wg.Done()
	}()

	for {
		buf := make([][]string, 0, 1024)
		for i := 0; i < 1024; i++ {
			tmp, err := cr.Read()
			if err == io.EOF {
				ch <- &buf
				close(ch)
				goto wait
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			buf = append(buf, tmp)
		}
		ch <- &buf
	}
wait:
	wg.Wait()
	return
}

func process(writer *csv.Writer, ch chan *[][]string) {
	for {
		buf, ok := <-ch
		if !ok && buf == nil {
			return
		}
		for y, col := range *buf {
			for x, row := range col {
				if row != "" {
					num, err := strconv.Atoi(row)
					if err != nil {
						(*buf)[y][x] = "0"
						continue
					}
					(*buf)[y][x] = strconv.Itoa(num + 1)
				} else {
					(*buf)[y][x] = "0"
				}
			}
		}
		err := writer.WriteAll(*buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		writer.Flush()
	}
}
