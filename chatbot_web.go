//session 3 - chatroom advanced - with userlist
package main

import (
	//"io"
	"strings"
	"fmt"
	"net/http"
//	"time"
	"io/ioutil"
	"stanleyxcnx/usenixgotour/ajax"
)

/*
//Ajax's purpose is to display chat messages
type Ajax struct {
	Chan chan string
}
*/

//Chatroom users
type Users map[string]bool

func (us Users) Add(name string) {
	us[name] = true
}
func (us Users) Remove(name string) {
	delete(us, name)
	//us.users[name] = false, false
}
func (us Users) String() string {
	names := make([]string, len(us))
	for n := range us {
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
/// end chatroom users


// Chatroom lobby receives all chat messages, and out put them to Ajax.
type ChatroomLobby ajax.Ajax

// Message - a generic interface so channel is flexible to accept any type of messages
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
	return "[" + chatter.Who + "] " + chatter.Msg + "\n"
}

type Status string



type Chatbot struct {
	Receive, Speak chan Message		//parrot's io
	Terms map[string]string			//dictionary of responses
	//Speak chan Message
}
func (parrot Chatbot) updateTerm(key, value string) {
	parrot.Terms[key] = value
}
func (parrot Chatbot) getTerm(key string) string {
	if v, ok := parrot.Terms[key]; ok {
			return v
	}
	return "plese tell parrot about " + key
}
func (parrot Chatbot) parseInstruction(inst string) (msg string, ok bool){
	msg, ok = "", false
	ts := strings.Fields(inst)
	if len(ts) < 4 {
		return
	}
	if tell, user, about, term := ts[0], ts[1], ts[2], ts[3]; tell == "tell" && about == "about" {
		msg = user + ": " + parrot.getTerm(term)
		ok = true
		return
	}
	return
}
func (parrot Chatbot) parseTerm(msg string) (term, definition string, ok bool) {
	term, definition, ok = "", "", false
	l := strings.Split(msg, "=");
	if len(l) == 1 {
		return
	}
	term = strings.TrimSpace(l[0])
	definition = strings.TrimSpace(l[1])
	ok = true
	return
}

func (parrot Chatbot) parseChat(chat Chatter) (recipient, msg string, ok bool) {
	recipient, msg, ok = "", "", false
	l := strings.Split(chat.Msg, ":");
	if len(l) == 1 {
		return
	}
	//
	if recipient = strings.TrimSpace(l[0]); recipient == "parrot" {
		msg = strings.Join(l[1:], ":")			//stich broken pieces back together, in case Msg contains multiple separator characters
		ok = true
		return
	}
	return

}



func (parrot Chatbot) Run() {
	fmt.Println("parrot has entered the room...\n")
	parrot.Speak <- Status("An invisible parrot has just ventured in.")
	for {

		switch m := (<-parrot.Receive).(type) {

		case Chatter:
			if _, msg, ok := parrot.parseChat(m); ok {
				if t, def, ok := parrot.parseTerm(msg); ok {
					fmt.Println("Debug: parrot received a training message " + m.Msg)
					parrot.updateTerm(t, def)

					ack := m.Who + ":" + msg
					chat := Chatter{Who: "parrot", Msg: ack}
					parrot.Speak <- chat
					break
				}
				if inst, ok := parrot.parseInstruction(msg); ok {
					fmt.Println("Debug: parrot received an instruction message " + msg)
					chat := Chatter{Who: "parrot", Msg: inst}
					parrot.Speak <- chat
					break
				}
			}

		default:
			fmt.Println("Debug: ignored by parrot")
		}
	}
}

type Chatroom struct {
	Ch chan Message
	Participants Users
	Lobby ChatroomLobby
	Parrot Chatbot
}


func (chatroom Chatroom) MainPage(resp http.ResponseWriter, req *http.Request) {
	if len(chatroom.Lobby.Chan) == 0 {
		fmt.Println("Debug: No msg in lobby channel.")
		return
	 }

	 chat, more := <-chatroom.Lobby.Chan
	 if more {
		 msg := fmt.Sprintf("%v", chat)
		 resp.Write([]byte(msg))
	 } else {
		 resp.Write([]byte("idle\n"))
	 }

}



func (chatroom Chatroom) Run() {
	fmt.Println("chatroom active...\n")
	//go chatroom.Parrot.Run()
	for {

		switch m := (<-chatroom.Ch).(type) {
		case Join:
			s := Status(chatroom.Participants.addAndPrint(m.Who))
			//chatroom.Participants.Add(m.Who)
			chatroom.Lobby.Chan <- Chatter{m.Who, m.Who + " has joined." }.String()
			// broadcast status
			//ajax.Chan <- Chatter{"*>>>", s.Info}
			go func() { chatroom.Ch <- s }()
		case Exit:
			s := Status(chatroom.Participants.removeAndPrint(m.Who))
			//chatroom.Participants.Remove(m.Who)
			chatroom.Lobby.Chan <- Chatter{m.Who, m.Who + " has lefted the room." }.String()

			//broadcast status
			//ajax.Chan <- Chatter{"*>>>", s.Info}
			go func() { chatroom.Ch <- s }()
		case Chatter:
			chatroom.Lobby.Chan <- m.String()
		case Status:
			//print Status info
			//chatroom.Lobby.Chan <- Chatter{"*>>>", string(m)}.String()
			chatroom.Lobby.Chan <- "*>>>" + string(m) + "\n"
		default:
			//time.Sleep(time.Millisecond *10)
			fmt.Println("ERROR: unknown chat message type")
		}
	}
}
func (chatroom Chatroom) join(resp http.ResponseWriter, req *http.Request) {
	d :=req.FormValue("id")
	chatroom.Ch <- Join{req.FormValue("id")}
	resp.Write([]byte("ok " + d))
}
func (chatroom Chatroom) exit(resp http.ResponseWriter, req *http.Request) {
	chatroom.Ch <- Exit{req.FormValue("id")}
	resp.Write([]byte("ok"))
}

func (chatroom Chatroom) say(resp http.ResponseWriter, req *http.Request) {
	chatroom.Ch <- Chatter{req.FormValue("id"), req.FormValue("msg")}
	//nosy parrot
	chatroom.Parrot.Receive <- Chatter{req.FormValue("id"), req.FormValue("msg")}
	resp.Write([]byte("ok"))
}


/////////////// web
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
	chatroom_lobby := ChatroomLobby{Chan: make(chan string, 100)}
	chatline := make(chan Message)
	parrot := Chatbot{Receive:make(chan Message), Terms: make(map[string]string, 99), Speak: chatline}
	chatroom := Chatroom{Ch: chatline, Participants: make(Users), Lobby: chatroom_lobby, Parrot: parrot }
	go chatroom.Run()
	go chatroom.Parrot.Run()

	http.HandleFunc("/join", chatroom.join)
	http.HandleFunc("/say", chatroom.say)
	http.HandleFunc("/exit", chatroom.exit)
	http.HandleFunc("/chatroom", chatroom.MainPage)

	http.HandleFunc("/ajax.js", ServJS);
	http.HandleFunc("/", WebIndex);

	http.ListenAndServe(":4000", nil)



}
