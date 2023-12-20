package main

import (
	"fmt"
	"net"
)

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		fmt.Printf("Received: %v", string(data))
		_, err = c.Write(data)
		if err != nil {
			panic("Write: " + err.Error())
		}
	}
}

func main() {
	/*
		lockPath := "/tmp/test_flock.lock"
		fileLock := flock.New(lockPath)
		for {
			locked, err := fileLock.TryLock()
			if err != nil {
				fmt.Printf("try lock error %v\n", err)
			} else if locked {
				break
			}
			fmt.Printf("try lock failed\n")
			time.Sleep(time.Second)
		}

		fmt.Printf("locked\n")
		time.Sleep(time.Hour)
	*/

	udsPath := "/tmp/echo.sock"
	//syscall.Unlink(udsPath)
	//if true {
	//	return
	//}

	l, err := net.Listen("unix", udsPath)
	if err != nil {
		println("listen error", err.Error())
		return
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			println("accept error", err.Error())
			return
		}

		go echoServer(fd)
	}
}
