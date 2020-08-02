package simple

import (
	"fmt"
	"time"
)

func printCount(c chan int) {
	num := 0
	for num >= 0 {
		num = <-c
		fmt.Println(num, " ")
	}
}

func poll(s chan string) {
	msg := "initial"

	select {
	case msg = <-s:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received", msg)
	}
}

func pollLoop(s chan string) {
	msg := "initial"

	for msg != "stop" {
		select {
		case msg = <-s:
			fmt.Println("received message", msg)
		default:
			fmt.Println("no message received", msg)
		}
		time.Sleep(time.Microsecond * 1)
	}
	fmt.Println("Stop signal received!!!")
}

func node(http chan string, grpc chan string) {
	httpRequest := "initial"
	grpcMessage := "initial"

	for grpcMessage != "stop" {
		select {
		case grpcMessage = <-grpc:
			fmt.Println("received gRPC message", grpcMessage)
		default:
			fmt.Println("no gRPC message received")
		}

		select {
		case httpRequest = <-http:
			fmt.Println("received HTTP message", httpRequest)
		default:
			fmt.Println("no message received")
		}

		time.Sleep(time.Microsecond * 1)
	}
	fmt.Println("gRPC stop signal received!!!")
}

func main() {
	c := make(chan int)
	a := []int{8, 3, 7, 1, 4}

	go printCount(c)

	for _, v := range a {
		c <- v
	}

	s := make(chan string)

	go pollLoop(s)

	//time.Sleep(time.Microsecond * 1)

	s <- "hello1"

	time.Sleep(time.Microsecond * 100)

	s <- "hello2"
	s <- "stop"

	time.Sleep(time.Microsecond * 100)
	fmt.Println("End of Main")
}
