package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

func (s *ServerConf) Configure() {
	Args := os.Args
	switch len(Args) {
	case 1:
		break
	case 2:
		s.Port = Args[1]
	default:
		fmt.Println("[Usage]: ./TCPChat $port")
		return
	}
}

func ShowImage(conn net.Conn) {
	file, err := os.Open("image.txt")
	if err != nil {
		conn.Write([]byte("Image Not found"))
		LogChannel <- fmt.Sprintf("Error image impossbile to draw: %s", conn.LocalAddr().String())
		LogChannel <- fmt.Sprint(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		conn.Write([]byte(line + "\n"))
	}

	if err := scanner.Err(); err != nil {
		LogChannel <- "Scanner error"
		LogChannel <- fmt.Sprint(err)
		return
	}

}

func AppendUser(name string, userConnInfo UserConnInfo) {
	ServConf.Mu.Lock()
	ServConf.Users[name] = userConnInfo //append into struct data about connected user
	ServConf.Mu.Unlock()
}

func TakeName(conn net.Conn) (string, error) {
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	buffer := make([]byte, 1024)
	length, err := conn.Read(buffer)
	if err != nil {
		LogChannel <- fmt.Sprintf("Error reading name: %s", conn.LocalAddr().String())
		LogChannel <- fmt.Sprint(err)
		errString := fmt.Sprintf("Error unresolved name: %v \n", err)
		conn.Write([]byte(errString))
		conn.Close()
		return "", err
	}
	name := string(buffer[:length-1])
	if _, ispresent := ServConf.Users[name]; ispresent {
		LogChannel <- fmt.Sprintf("Connection attempt from %s", conn.LocalAddr().String())
		conn.Write([]byte("Username taken, choose different one \n"))
		return "", errors.New("USERNAME TAKEN")
	}
	LogChannel <- fmt.Sprintf("Connection from %s, with USERNAME: %s", conn.LocalAddr().String(), name)
	return name, nil
}

func TakeInput(conn net.Conn, name string, BrCaChannel chan Message) {
	writeHistory(conn)
	sayHello(name) //broadcasting "%s name has joined our chat..."
	for {
		buffer := make([]byte, 1024)
		length, err := conn.Read(buffer)
		if err != nil {
			LogChannel <- fmt.Sprintf("Connection closed: %s, with USERNAME %s", conn.LocalAddr().String(), name)
			conn.Close()
			ServConf.Mu.Lock()
			delete(ServConf.Users, name)
			ServConf.Mu.Unlock()
			sayGoodBye(name)
			return
		}
		input := Message{
			Message: string(buffer[:length]),
			Name:    name,
			Time:    time.Now().Format(time.DateTime),
		}
		BrCaChannel <- input
	}
}

func BroadCast(BrCaChannel chan Message) {
	for input := range BrCaChannel {
		if len(input.Message) == 1 {
			conn := ServConf.Users[input.Name].UserConnection                                            //to write input string -> [time][name]:#
			conn.Write([]byte(fmt.Sprintf("\r[%s][%s]:", time.Now().Format(time.DateTime), input.Name))) //case message is empty prints nothing
			continue
		}
		ServConf.Mu.Lock()
		for _, v := range ServConf.Users {
			message := fmt.Sprintf("\r[%s][%s]:%s", input.Time, input.Name, input.Message) //string to broadcast

			if v.UserName == input.Name { //if connection is of the user whos' message is broadcasting
				ServConf.Chat = append(ServConf.Chat, message)
				ServConf.LastString = message            //I only use this for empty space (:/)
				emptySpace(v.UserConnection, input.Name) //do we really need this??
				v.UserConnection.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(time.DateTime), input.Name)))
				continue
			}
			emptySpace(v.UserConnection, v.UserName)                                                               //to clean the input string of the users broadcasting to
			v.UserConnection.Write([]byte(message))                                                                //broadcasting message
			v.UserConnection.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(time.DateTime), v.UserName))) //input stiring for each user -> [time][name]:#
		}
		ServConf.Mu.Unlock()
		// conn := ServConf.Users[input.Name].UserConnection
		// emptySpace(conn, input.Name)
		// conn.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(time.DateTime), input.Name)))
	}
}

func writeHistory(conn net.Conn) {
	for _, v := range ServConf.Chat {
		conn.Write([]byte(v))
	}
}

func sayHello(name string) {
	ServConf.Mu.Lock()
	for _, v := range ServConf.Users {
		if v.UserName == name {
			v.UserConnection.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(time.DateTime), v.UserName))) //input string -> [time][name]:#
			continue
		}
		emptySpace(v.UserConnection, v.UserName)                                           //writes on last line (last line for each user is their input string, so we need to delete that string, and start to write on it)
		v.UserConnection.Write([]byte(fmt.Sprintf("\r%s has joined our chat...\n", name))) //has joined the chat
		ServConf.LastString = fmt.Sprintf("%s has joined our chat...\n", name)
		v.UserConnection.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(time.DateTime), v.UserName))) //input string different for each user -> [time][name]: #
	}
	ServConf.Mu.Unlock()
}

func sayGoodBye(name string) {
	ServConf.Mu.Lock()
	for _, v := range ServConf.Users {
		if v.UserName == name {
			continue
		}
		emptySpace(v.UserConnection, v.UserName)
		v.UserConnection.Write([]byte(fmt.Sprintf("\r%s has left our chat...\n", name)))
		ServConf.LastString = fmt.Sprintf("\r%s has left our chat...\n", name)
		v.UserConnection.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(time.DateTime), v.UserName)))
	}
	ServConf.Mu.Unlock()
}

func emptySpace(conn net.Conn, name string) {
	var empty string

	for i := 0; i < max(len(name)+24, len(ServConf.LastString)); i++ {
		empty += " "
	}
	conn.Write([]byte(fmt.Sprint("\r" + empty + "\r")))
}
