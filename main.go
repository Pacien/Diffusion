/*

	This file is part of Diffusion (https://github.com/Pacien/Diffusion)

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in
	all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.

*/

// Diffusion
package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
	"strings"
)

var broadcastKey = flag.String("broadcastKey", "", "Secret key for braodcasting.")
var clientKey = flag.String("revieveKey", "", "Secret key for recieving.")

func splitRequest(ws *websocket.Conn, base string) (channelName string, arguments []string) {
	request := strings.Split(ws.Request().RequestURI[len(base):], "?")
	channelName = request[0]
	if len(request) > 1 {
		arguments = strings.Split(request[1], "&")
	}
	return
}

func auth(submittedKey string, key *string) bool {
	if *key == "" || submittedKey == *key {
		return true
	}
	return false
}

var clients = make(map[string]map[int]chan []byte)

func broadcastHandler(ws *websocket.Conn) {

	channelName, arguments := splitRequest(ws, "/b/")

	log.Println("New broadcaster on #" + channelName)

	if !auth(arguments[0], broadcastKey) {
		return
	}

	var message = make([]byte, 512)

	for {
		ws.Read(message)

		log.Println("#" + channelName + ": " + string(message))

		for _, client := range clients[channelName] {
			client <- message
		}
	}

}

func clientHandler(ws *websocket.Conn) {

	channelName, arguments := splitRequest(ws, "/")

	log.Println("New client on #" + channelName)

	if !auth(arguments[0], clientKey) {
		return
	}

	var channel = make(chan []byte)

	clientID := len(clients[channelName])
	if clientID < 1 {
		clients[channelName] = make(map[int]chan []byte)
	}
	clients[channelName][clientID] = channel
	defer delete(clients[channelName], clientID)

	for {
		ws.Write(<-channel)
	}

}

func main() {

	flag.Parse()

	http.Handle("/b/", websocket.Handler(broadcastHandler))
	http.Handle("/", websocket.Handler(clientHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
