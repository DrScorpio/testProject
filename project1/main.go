package main

import "fmt"

func main() {
	chin := make(chan int, 5)
	chout := make(chan int, 5)
	l := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < 5; i++ {
		go process(chin, chout)
	}
	for _, el := range l {
		chin <- el
	}
	close(chin)
	res := make([]int, 0, 10)
	for {
		select {
		case num := <-chout:
			res = append(res, num)
		default:
		}
		if len(res) == 10 {
			break
		}
	}
	fmt.Println(res)
}

func process(in, out chan int) {
	for {
		select {
		case num := <-in:
			out <- num + 1
		default:

		}
	}

}
