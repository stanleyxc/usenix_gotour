
package main

import (
	//"io"
	//"strings"
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
)


//session 3
//faking ajax
type Ajax struct {
	Chan chan string
}

func (ajax Ajax) ShowMessage(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("hello"))
	for msg, more := <-ajax.Chan; more; msg, more = <-ajax.Chan {
		resp.Write([]byte(msg))
		time.Sleep(time.Millisecond *10)
	}
}

func (ajax Ajax) LogPage(resp http.ResponseWriter, req *http.Request) {
	 log, more := <-ajax.Chan
	 if more {
		 resp.Write([]byte(log))
	 } else {
		 resp.Write([]byte("no more log activity\n"))
	 }
}

///// end faking ajax

func (ajax Ajax) GenerateLogs() {
		go func() {
			for i:= 0; i < 10000; i++ {
				text := fmt.Sprintf("Log message %d\n", i)
				ajax.Chan <- text
				time.Sleep(time.Second *1)
			}
			close(ajax.Chan)
		}()
}
func (ajax Ajax) Log(msg string) {
	ajax.Chan <- msg

}

func WebIndex(resp http.ResponseWriter, req *http.Request) {

	 HTML, _ := ioutil.ReadFile("./WebLog.html")
	 resp.Write(HTML)
}
func ServJS(resp http.ResponseWriter, req *http.Request) {
	JS, _ := ioutil.ReadFile("./ajax.js")
	resp.Write(JS)
}

func main() {

	ajax := Ajax{Chan: make(chan string, 100)}
	ajax.GenerateLogs()

	http.HandleFunc("/log", ajax.LogPage)
	http.HandleFunc("/ajax.js", ServJS)
	http.HandleFunc("/", WebIndex);
	http.ListenAndServe(":4000", nil)



}
