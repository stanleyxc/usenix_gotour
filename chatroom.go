//session 3 - chatroom
package main

import (
	//"io"
	//"strings"
	"fmt"
	"net/http"
	//"time"
)

type Ajax struct {
	Chan chan string
}


func (ajax Ajax) ShowMessage(resp http.ResponseWriter, req *http.Request) {
	//resp.Write([]byte("hello\n"))
	//c := Chatter{"stanley", "who's there"}
	//resp.Write([]byte(fmt.Sprintf("%v", c)))

	 chat, more := <-ajax.Chan
	 if more {
		 msg := fmt.Sprintf("%v", chat)
		 resp.Write([]byte(msg))
	 } else {
		 resp.Write([]byte("idle\n"))
	 }

}




/* Message - a generic interface so channel is flexible to accept any type of messages */
type Message interface {
}
type Join struct {
	Who string
}
type Exit struct {
	Who string
}
type Chatter struct {
	Who string
	Msg string
}

func (chatter Chatter) String() string {
	return chatter.Who + ": " + chatter.Msg
}

type Chatroom struct {
	Channel chan Message
	//View chan string
	ajax Ajax
}

func (chatroom Chatroom) run() {
	fmt.Println("chatroom initialized\n")
	for {
		switch m := (<-chatroom.Channel).(type) {
		case Join:
			chatroom.ajax.Chan <- Chatter{m.Who, m.Who + " has joined." }.String()
		case Exit:
			chatroom.ajax.Chan <- Chatter{m.Who, m.Who + " has lefted the room." }.String()
		case Chatter:
			chatroom.ajax.Chan <- m.String()
		default:
			fmt.Println("Unknown chat message type")
		}
	}
}
func (chatroom Chatroom) join(resp http.ResponseWriter, req *http.Request) {
	chatroom.Channel <- Join{req.FormValue("id")}
	resp.Write([]byte("ok"))
}
func (chatroom Chatroom) exit(resp http.ResponseWriter, req *http.Request) {
	chatroom.Channel <- Exit{req.FormValue("id")}
	resp.Write([]byte("ok"))
}

func (chatroom Chatroom) say(resp http.ResponseWriter, req *http.Request) {
	chatroom.Channel <- Chatter{req.FormValue("id"), req.FormValue("msg")}
	resp.Write([]byte("ok"))
}




func main() {
	ajax := Ajax{Chan: make(chan string)}
	chatroom := Chatroom{Channel: make(chan Message), ajax: ajax}
	go chatroom.run()

	http.HandleFunc("/join", chatroom.join)
	http.HandleFunc("/say", chatroom.say)
	http.HandleFunc("/exit", chatroom.exit)
	http.HandleFunc("/chatroom", ajax.ShowMessage)
	http.ListenAndServe(":4000", nil)



}
