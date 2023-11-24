package main

import (
	"fmt"
	"net"
	"netcat/server"
)

func main() {
	servConf := server.ServConf
	servConf.Configure()
	go server.LoggerCreate()
	logChannel := server.LogChannel

	logChannel <- fmt.Sprintf("Attempt listening on localhost:%s", server.ServConf.Port)
	listener, err := net.Listen("tcp", "localhost:"+servConf.Port)
	if err != nil {
		fmt.Println(err)
		logChannel <- fmt.Sprintf("Unseccessfull listening attempt on port:%s", servConf.Port)
		logChannel <- fmt.Sprint(err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening of localhost:%s\n", servConf.Port)

	BrCaChannel := make(chan server.Message)
	go server.BroadCast(BrCaChannel)

	for {
		conn, err := listener.Accept()
		logChannel <- fmt.Sprintf("Connection attempt from %s", conn.LocalAddr().String())
		if len(server.ServConf.Users) == 10 {
			logChannel <- fmt.Sprintf("Connection attempt from %s, connection not possible, group is full", conn.LocalAddr().String())
			conn.Write([]byte(string("The chat is full connect later!")))
			conn.Close()
			continue
		}
		if err != nil {
			logChannel <- fmt.Sprintf("Connection error, connection unsuccessfull: %v", conn.LocalAddr().String())
			logChannel <- fmt.Sprint(err)
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
		server.LogChannel <- fmt.Sprint(err)
		return
	}

	// channel := make(chan server.Message)

	userConnInfo := server.CreateUserConnInfo(name, conn) //smth like constructor

	// UserConnInfo := server.UserConnInfo{
	// 	UserName: name,
	// 	// UserChannel:    channel,
	// 	UserConnection: conn,
	// }

	server.AppendUser(name, userConnInfo)

	go server.TakeInput(conn, name, BrCaChannel)
}
