package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// NewUser NewUser函数用于创建一个新的User对象
func NewUser(conn net.Conn) *User {
	//获取连接的远程地址
	var userAddr = conn.RemoteAddr().String()
	//创建一个新的User对象
	user := &User{
		//设置User对象的Name属性为远程地址
		Name: userAddr,
		//设置User对象的Addr属性为远程地址
		Addr: userAddr,
		//创建一个channel用于接收消息
		C: make(chan string),
		//设置User对象的conn属性为传入的连接
		conn: conn,
	}
	//启动监听当前User channel消息的GoRoutine
	go user.ListenMessage()
	//返回User对象
	return user
}

// ListenMessage 监听当前User channel的方法,一旦有消息,就直接发送给对端客户端
func (u *User) ListenMessage() {
	// 无限循环
	for {
		// 从通道C中读取消息
		msg := <-u.C
		// 将消息写入连接
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			return
		}
	}
}
