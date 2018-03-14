//session 3 - chatroom advanced - with userlist
package main

import (
	//"io"
	"strings"
	"fmt"
	"net/http"
		"io/ioutil"
		//"stanleyxcnx/usenixgotour/ajax"
//	"time"
)

//Chatroom users
type Users struct {
	users map[string]bool
}

func (us Users) Add(name string) {
	us.users[name] = true
}
func (us Users) Remove(name string) {
	delete(us.users, name)
	//us.users[name] = false, false
}
func (us Users) String() string {
	names := make([]string, len(us.users))
	for n := range us.users {
		names = append(names, n)
	}
	return strings.Join(names, ", ")
}
func (us Users) addAndPrint(name string) string {

	x := "[+" + name + ", " + us.String() + "]"
	us.Add(name)
	return x
}
func (us Users) removeAndPrint(name string) string {
	us.Remove(name)
	x := "[-" + name + ", " + us.String() + "]"
	return x
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
type Status struct {
	Info string
}



func (chatter Chatter) String() string {
	return chatter.Who + ": " + chatter.Msg
}

type Chatroom struct {
	Ch chan Message
	Participants Users
}


func (chatroom Chatroom) run(ajax Ajax) {
	fmt.Println("chatroom active...\n")
	for {

		switch m := (<-chatroom.Ch).(type) {
		case Join:
			s := Status{Info: chatroom.Participants.addAndPrint(m.Who)}
			//chatroom.Participants.Add(m.Who)
			ajax.Chan <- Chatter{m.Who, m.Who + " has joined." }
			// broadcast status
			//ajax.Chan <- Chatter{"*>>>", s.Info}
			go func() { chatroom.Ch <- s }()
		case Exit:
			s := Status{Info: chatroom.Participants.removeAndPrint(m.Who)}
			//chatroom.Participants.Remove(m.Who)
			ajax.Chan <- Chatter{m.Who, m.Who + " has lefted the room." }

			//broadcast status
			//ajax.Chan <- Chatter{"*>>>", s.Info}
			go func() { chatroom.Ch <- s }()
		case Chatter:
			ajax.Chan <- m
		case Status:
			//print Status info
			ajax.Chan <- Chatter{"*>>>", m.Info}
		default:
			//time.Sleep(time.Millisecond *10)
			fmt.Println("Unknown chat message type")
		}
	}
}
func (chatroom Chatroom) join(resp http.ResponseWriter, req *http.Request) {
	d :=req.FormValue("id")
	chatroom.Ch <- Join{req.FormValue("id")}
	//fmt.Printf("form: %v\n", req.Form)
	//fmt.Printf("%v joined.\n", d);
	//req.ParseForm()
	//fmt.Printf("Hello, %s!", req.PostFormValue("id"))
	resp.Write([]byte("ok " + d))
}
func (chatroom Chatroom) exit(resp http.ResponseWriter, req *http.Request) {
	chatroom.Ch <- Exit{req.FormValue("id")}
	fmt.Printf("%v exited.\n", req.FormValue("id"));
	resp.Write([]byte("ok"))
}

func (chatroom Chatroom) say(resp http.ResponseWriter, req *http.Request) {
	chatroom.Ch <- Chatter{req.FormValue("id"), req.FormValue("msg")}
	fmt.Printf("%v said: %v\n", req.FormValue("id"), req.FormValue("msg"));
	resp.Write([]byte("ok"))
}


//Ajax's purpose is to display chat messages
type Ajax struct {
	Chan chan Chatter
}

func (ajax Ajax) ShowMessage(resp http.ResponseWriter, req *http.Request) {
	 if len(ajax.Chan) == 0 {
		 fmt.Println("no data in queue")
		 return
		}
	 chat, more := <-ajax.Chan
	 if more {
		 msg := fmt.Sprintf("%v\n", chat)
		 resp.Write([]byte(msg))
	 } else {
		 resp.Write([]byte("idle\n"))
	 }


}


func WebIndex(resp http.ResponseWriter, req *http.Request) {

	 HTML, _ := ioutil.ReadFile("./chatroom.html")
	 resp.Header().Set("Content-Type", "text/html")
	 resp.Write(HTML)
}

func  ServJS(resp http.ResponseWriter, req *http.Request) {
	JS, _ := ioutil.ReadFile("./ajax.js")
	resp.Write(JS)
}

func main() {
	ajax := Ajax{Chan: make(chan Chatter, 100)}

	chatroom := Chatroom{Ch: make(chan Message), Participants: Users{users: make(map[string]bool)} }
	go chatroom.run(ajax)

	http.HandleFunc("/join", chatroom.join)
	http.HandleFunc("/say", chatroom.say)
	http.HandleFunc("/exit", chatroom.exit)
	http.HandleFunc("/chatroom", ajax.ShowMessage)
		http.HandleFunc("/ajax.js", ServJS);
	http.HandleFunc("/", WebIndex);

	fmt.Println("web chat running on ':4000/'")

	http.ListenAndServe(":4000", nil)



}
