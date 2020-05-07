package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		panic("无效的域名")
	}

	domain := strings.ToUpper(os.Args[1])

	url := "https://raw.githubusercontent.com/daliyo/whois-spider/master/WHOIS.txt"

	res, err := http.Get(url)
	checkErr(err)

	snr := bufio.NewScanner(res.Body)
	srvlist := map[string]string{}
	for snr.Scan() {
		parts := strings.Split(snr.Text(), ": ")
		srvlist[strings.ToUpper(parts[0])] = parts[1]
	}

	parts := strings.SplitN(domain, ".", 2)

	if srv, ok := srvlist["."+parts[1]]; ok {

		addr, err := net.ResolveTCPAddr("tcp", srv+":43")
		checkErr(err)

		conn, err := net.DialTCP("tcp", nil, addr)
		checkErr(err)

		defer conn.Close()

		_, err = conn.Write([]byte(domain + "\n"))
		checkErr(err)

		buff, err := ioutil.ReadAll(conn)
		checkErr(err)

		fmt.Println(string(buff))
	} else {
		fmt.Println("未找到有效的WHOIS服务器")
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
