package main

import (
	"fmt"
	"time"
)

const numNodes = 5

type Message struct {
	fromNode    int
	toNode      int
	messageType string
	payload     string
}

func Node(id int, httpChans [numNodes]chan string, grpcChans [numNodes]chan Message) {

	fmt.Println("Node:", id, "starting...")
	httpRequest := "initial"
	currentMessage := Message{-1, -1, "initial", "initial"}
	//chain := []string{}
	onboardingPool := []string{}
	sendBuffer := []Message{}

	for currentMessage.messageType != "stop" {
		select {
		case currentMessage = <-grpcChans[id]:
			fmt.Println("Node:", id, "received gRPC message", currentMessage)
		default:
			fmt.Println("Node:", id, "no gRPC message received")
		}

		select {
		case httpRequest = <-httpChans[id]:
			fmt.Println("Node:", id, "received HTTP request", httpRequest)
			for i := 0; i < numNodes; i++ {
				if i == id {
					onboardingPool = append(onboardingPool, httpRequest)
					continue
				}
				sendBuffer = append(sendBuffer, Message{id, i, "onboarding", httpRequest})
			}
		default:
			fmt.Println("Node:", id, "no HTTP request received")
		}

		if len(sendBuffer) > 0 {

			select {
			case grpcChans[sendBuffer[0].toNode] <- sendBuffer[0]:
				fmt.Println("Node:", id, "sent GRPC message", sendBuffer[0], "to node", sendBuffer[0].toNode)
				sendBuffer = sendBuffer[1:]
			case <-time.After(10 * time.Millisecond):
				fmt.Println("Node:", id, "failed to send GRPC message", sendBuffer[0], "to node", sendBuffer[0].toNode)
			}
		}

		time.Sleep(time.Millisecond * 1)
	}
	fmt.Println("Node:", id, "gRPC stop signal received. Shutting down...")
}

func main() {
	fmt.Println("main start")

	var httpChans [numNodes]chan string
	for i := range httpChans {
		httpChans[i] = make(chan string)
	}

	var grpcChans [numNodes]chan Message
	for i := range grpcChans {
		grpcChans[i] = make(chan Message)
	}

	for i := 0; i < numNodes; i++ {
		go Node(i, httpChans, grpcChans)
	}

	time.Sleep(time.Millisecond * 1)

	httpChans[0] <- "tx1"

	time.Sleep(time.Millisecond * 10)

	grpcChans[0] <- Message{-1, -1, "stop", ""}
	grpcChans[1] <- Message{-1, -1, "stop", ""}

	time.Sleep(time.Millisecond * 10)
	fmt.Println("End of Main")
}
