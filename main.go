package main

func main() {
	var s = NewServer("127.0.0.1", 8888)
	s.Start()
}
