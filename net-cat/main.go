package main

import (
	"fmt"
	"log"
	"net"
	"netcat/server"
	// "github.com/jroimartin/gocui"
)

func main() {
	// g, err := gocui.NewGui(gocui.OutputNormal)
	// if err != nil {
	// 	// handle error
	// 	return
	// }
	// defer g.Close()

	// // Set GUI managers and key bindings
	// // ...

	// if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
	// 	// handle error
	// 	return
	// }
	//this shit is dangerous
	servConf := server.ServConf
	servConf.Configure()
	file := server.LoggerCreate(log.Default())
	defer file.Close()
	// logChannel := server.LogChannel

	log.Printf("Attempt listening on localhost:%s", server.ServConf.Port)
	listener, err := net.Listen("tcp", "localhost:"+servConf.Port)
	if err != nil {
		// fmt.Println(err)
		log.Printf("Unseccessfull listening attempt on port:%s", servConf.Port)
		log.Println(err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening of localhost:%s\n", servConf.Port)

	BrCaChannel := make(chan server.Message)
	go server.BroadCast(BrCaChannel)

	for {
		conn, err := listener.Accept()
		log.Printf("Connection attempt from %s", conn.LocalAddr().String())
		if len(server.ServConf.Users) == 10 {
			log.Printf("Connection attempt from %s, connection not possible, group is full", conn.LocalAddr().String())
			conn.Write([]byte(string("The chat is full connect later!")))
			conn.Close()
			continue
		}
		if err != nil {
			log.Printf("Connection error, connection unsuccessfull: %v", conn.LocalAddr().String())
			log.Println(err)
			fmt.Printf("Connection error: %v \n", err)
			continue
		}
		go HandleConnection(conn, BrCaChannel)
	}
}

func HandleConnection(conn net.Conn, BrCaChannel chan server.Message) {
	server.ShowImage(conn)
	var name string
	var err error
	for { //until user inputs username that is not taken or error returns other error than "USERNAME TAKEN"
		name, err = server.TakeName(conn)
		if err == nil || err.Error() != "USERNAME TAKEN" {
			break
		} else if err.Error() == "USERNAME TAKEN" {
			continue
		}
	}

	if err != nil { //connection closed in function if err occured
		log.Println(err)
		return
	}
	userConnInfo := server.CreateUserConnInfo(name, conn) //smth like constructor

	server.AppendUser(name, userConnInfo)

	go server.TakeInput(conn, name, BrCaChannel)
}
