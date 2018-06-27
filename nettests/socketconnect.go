package nettests

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"
)

var toHex = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)
	return
}

func receiveSock(conn net.Conn, c chan string) {
	DEBUG := 0

	if DEBUG > 0 {
		fmt.Fprintln(os.Stderr, "in receiveSock")
	}

	msg := make([]byte, 4096)

	l, err := bufio.NewReader(conn).Read(msg)

	if err != nil {
		return
	}
	if l > 0 {
		if DEBUG > 0 {
			fmt.Fprintf(os.Stderr, ">>>message is %s length(%d)", msg, l)
		}
	}

	if DEBUG > 0 {
		fmt.Fprintf(os.Stderr, "message is %s", msg)
	}

	c <- fmt.Sprintf("%d:%s", l, hex.EncodeToString(msg[:l]))

}

func SocketConnect(srcIP net.Addr, host string, port int) (string, error) {
	DEBUG := 0

	constr := fmt.Sprintf("%s:%d", host, port)
	//	d := &net.Dialer{LocalAddr: "1.1.1.1"}

	if DEBUG > 0 {
		fmt.Fprintf(os.Stderr, "Socket connect>>>%s<<<\n", constr)
	}

	timeStart := time.Now()

	conn, err := net.DialTimeout("tcp", constr, time.Millisecond*400)

	rtt := time.Now().Sub(timeStart)

	if DEBUG > 0 {
		fmt.Printf("RTT>>>%v<<<\n", rtt)
	}

	if err != nil {
		if DEBUG > 0 {
			fmt.Printf(">>>>>err:%v<<<<\n", err)
		}
		return fmt.Sprintf("Error:%v", err), err
	} else {
		if DEBUG > 0 {
			fmt.Fprintln(os.Stderr, "no error going to receive\n")
		}
		c := make(chan string)
		msg := ""

		go receiveSock(conn, c)

		select {
		case msg = <-c:
		case <-time.After(time.Millisecond * 400):
			msg = "0:"
		}

		if DEBUG > 0 {
			fmt.Printf(">>>%v<<<\n", msg)
		}
		return msg, nil
	}

}
