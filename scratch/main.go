package main

import "fmt"

func main() {

	ch := make(chan int, 1);
	for i := 0; i < 500; i++ {
		go process(i, ch);
	}

	for i := 0; i < 500; i++ {
		select {
		case result := <-ch:
			fmt.Println(result)
		}
	}
	fmt.Println("done")
}

func process(i int, ch chan int) {
	ch <- i
}