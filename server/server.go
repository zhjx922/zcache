package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"

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

func newAccpet(c net.Conn) {
	defer c.Close()

	buffer := make([]byte, 0)

	for {
		tmpBuffer := make([]byte, 1024)
		n, err := c.Read(tmpBuffer)
		if err != nil {
			//log.Println("conn read error:", err)
			return
		}
		//log.Println("tmpBuffer:", string(tmpBuffer))
		buffer = checkData(c, append(buffer, tmpBuffer[:n]...))
	}
}

func checkData(c net.Conn, buffer []byte) []byte {
	index := bytes.Index(buffer, []byte("\r\n"))

	//数据不完整，继续获取
	if index <= 0 {
		log.Println("数据不完整，继续获取")
		return buffer
	}

	command := bytes.Split(buffer[:index], []byte(" "))
	log.Println("当前的命令是:", string(command[0]))
	if string(command[0]) == "set" {
		log.Println("当前的命令是set")
		//当前block总长度
		l, err := strconv.Atoi(string(command[4]))
		if err != nil {
			log.Println("conn read error:", err)
			return []byte("")
		}
		total := index + 2 + l + 2
		//数据不完整，继续获取
		log.Println("当前数据长度:", len(buffer))
		log.Println("Block总长度:", total)
		if len(buffer) < total {
			log.Println("数据不完整，继续获取")
			return buffer
		}
		log.Println("数据完整啦")
		flags, _ := strconv.ParseInt(string(command[2]), 10, 64)
		expire, _ := strconv.ParseInt(string(command[3]), 10, 64)
		if expire == 0 {
			expire = flags
			flags = 0
		}
		cache.Add(string(command[1]), string(buffer[index+2:total-2]), flags, expire)
		c.Write([]byte("STORED\r\n"))
		if len(buffer) > total {
			fmt.Printf("%v\n", buffer)
			log.Println("粘包，继续处理")
			return buffer[total+1:]
		}
	} else if string(command[0]) == "get" {
		log.Println("当前的命令是get")
		value, flags, ok := cache.Get(string(command[1]))
		if ok {
			c.Write([]byte("VALUE " + string(command[1]) + " " + strconv.FormatInt(flags, 10) + " " + strconv.Itoa(len(value)) + "\r\n"))
			c.Write([]byte(value + "\r\n"))
		}
		c.Write([]byte("END\r\n"))
		return buffer[index+2:]
	} else if string(command[0]) == "delete" {
		log.Println("当前的命令是delete")
		if d := cache.Delete(string(command[1])); d {
			c.Write([]byte("DELETED\r\n"))
		} else {
			c.Write([]byte("NOT_FOUND\r\n"))
		}
		return buffer[index+2:]
	}
	return []byte("")
}

//Start 启动
func Start(port string) {
	l, err := net.Listen("tcp", ":"+port)
	defer l.Close()

	if err != nil {
		log.Fatalf("Got an error:  %s", err)
	}

	fmt.Println("Server启动成功")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting from %s", l)
		} else {
			go newAccpet(conn)
		}
	}
}
