package server

import (
	"net"
	"sync"
)

var ServConf = initiateServConf()
// var LogChannel = make(chan string)

type ServerConf struct {
	Mu         sync.Mutex
	Port       string
	Users      map[string]UserConnInfo
	Chat       []string
	LastString string
}

func initiateServConf() *ServerConf {
	return &ServerConf{
		Port:       "8989",
		Users:      make(map[string]UserConnInfo),
		Chat:       []string{},
		LastString: "",
	}
}

type UserConnInfo struct {
	UserName string
	// UserChannel    chan Message
	UserConnection net.Conn
}

func CreateUserConnInfo(username string, connection net.Conn) UserConnInfo {
	return UserConnInfo{
		UserName:       username,
		UserConnection: connection,
	}
}

type Message struct {
	Message string
	Name    string
	Time    string
}

// func Init() ServerConf {
// 	return ServerConf{Port: "8989"}
// }
