package main

import (
	"fmt"
	"net"
	"sync"
)

// Server 结构体，表示服务器
type Server struct {
	Ip   string
	Port int
	// 在线用户的列表
	OnlineMap map[string]*User
	// mapLock 是一个读写互斥锁，用于保护 map 的并发访问
	mapLock sync.RWMutex
	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个新的服务器实例
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// ListenMessager 监听Message广播信息channel的GoRoutine,一旦有消息就发送给全部在线的User
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// BroadCast 广播消息给所有在线用户
func (s *Server) BroadCast(user *User, msg string) {
	var sendMeg = fmt.Sprintf("[%s]: %s", user.Addr, msg)
	// 遍历 map，给每个 user 发送 msg
	s.Message <- sendMeg
}

// Handler 处理客户端连接
func (s *Server) Handler(conn net.Conn) {
	// 当前链接的业务
	// fmt.Println("链接建立成功")
	var user = NewUser(conn)
	// 用户上线,加入OnlineMap
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()
	// 广播当前用户上线消息
	s.BroadCast(user, "上线了")
	// 当前handler阻塞
	select {}
}

// Start 启动服务器的接口
func (s *Server) Start() {
	//Socket监听
	var listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//关闭监听
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("listener.Close err:", err)
		}
	}(listener)
	//启动监听Message的GoRoutine
	go s.ListenMessager()
	// 接受客户端链接并处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		// 处理客户端链接
		go s.Handler(conn)
	}

}
