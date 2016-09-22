package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/zhjx922/zcache/lru"
)

var cache *lru.Cache

func init() {
	cache = lru.NewCache(0)
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Error:", err.Error())
		return
	}
}

func accept(c net.Conn) {
	defer c.Close()
	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)

		if err != nil {
			log.Println("conn read error:", err)
			return
		}
		params := strings.Fields(string(buf[:n]))

		if params[0] == "set" {
			expire, _ := strconv.ParseInt(params[2], 10, 64)
			cache.Add(params[1], params[5], expire)
			c.Write([]byte("STORED\r\n"))
		} else if params[0] == "get" {
			value, ok := cache.Get(params[1])
			if ok {
				c.Write([]byte("VALUE " + params[1] + " 0 " + strconv.Itoa(len(value)) + "\r\n"))
				c.Write([]byte(value + "\r\n"))
			}
			c.Write([]byte("END\r\n"))
		} else if params[0] == "delete" {
			//DELETED\r\n
			cache.Delete(params[1])
			c.Write([]byte("DELETED\r\n"))
		}

		fmt.Printf("%v\n", params)

	}
}

//Start 启动
func Start(port string) {
	l, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Got an error:  %s", err)
	}

	fmt.Println("Server启动成功")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting from %s", l)
		} else {
			go accept(conn)
		}
	}
}
