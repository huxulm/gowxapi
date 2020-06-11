package helper

import "fmt"

var MsgChannel = make(chan string)

func SendMsg(msg string) {
	go func() {
		MsgChannel <- msg
		fmt.Println("send one message " + msg)
	}()
}

func InitCallback(f func(string)) {
	for {
		msg := <-MsgChannel
		f(msg)
	}
}
