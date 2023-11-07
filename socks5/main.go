package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

func main() {
	server, err := net.Listen("tcp", ":1080")
	if err != nil {
		panic(err)
	}

	fmt.Println("proxy run on port 1080")
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}
		go process(conn)
	}
}

func process(cli net.Conn) {
	if err := Socks5Auth(cli); err != nil {
		fmt.Printf("Auth error: %v\n", err)
		cli.Close()
		return
	}

	target, err := Socks5Connect(cli)
	if err != nil {
		fmt.Printf("connect error: %v\n", err)
		cli.Close()
		return
	}

	Socks5Forward(cli, target)
}

func Socks5Auth(cli net.Conn) error {
	buf := make([]byte, 256)

	// 读取 VER 和 NMENTHODS
	n, err := io.ReadFull(cli, buf[:2])
	if n != 2 {
		return fmt.Errorf("Reading error: %v\n", err)
	}
	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return fmt.Errorf("Version not support: %v\n", ver)
	}
	// 读取 MENTHODS
	n, err = io.ReadFull(cli, buf[:nMethods])
	if n != nMethods {
		return fmt.Errorf("Reading error: %v\n", err)
	}

	// 无需认证
	n, err = cli.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return fmt.Errorf("write rsp: " + err.Error())
	}

	return nil
}

func Socks5Connect(cli net.Conn) (net.Conn, error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(cli, buf[:4])
	if n != 4 {
		return nil, fmt.Errorf("Reading error: %v\n", err)
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, fmt.Errorf("invalid ver/cmd")
	}

	addr := ""
	switch atyp {
	case 1:
		n, err := io.ReadFull(cli, buf[:4])
		if n != 4 {
			return nil, fmt.Errorf("Reading error: %v\n", err)
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case 3:
		n, err := io.ReadFull(cli, buf[:1])
		if n != 1 {
			return nil, fmt.Errorf("Reading error: %v\n", err)
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(cli, buf[:addrLen])
		if n != addrLen {
			return nil, fmt.Errorf("Reading error: %v\n", err)
		}
		addr = string(buf[:addrLen])
	case 4:
		return nil, errors.New("IPv6 not support")
	default:
		return nil, errors.New("invalid atyp")
	}

	n, err = io.ReadFull(cli, buf[:2])
	if n != 2 {
		return nil, fmt.Errorf("Reading error: %v\n", err)
	}
	port := binary.BigEndian.Uint16(buf[:2])

	destAddrPort := fmt.Sprintf("%s:%d", addr, port)
	dest, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		return nil, fmt.Errorf("Dial error: %v\n", err)
	}

	n, err = cli.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if n != 10 || err != nil {
		return nil, fmt.Errorf("write rsp: " + err.Error())
	}
	return dest, nil
}

func Socks5Forward(cli, target net.Conn) {
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}

	go forward(cli, target)
	go forward(target, cli)
}
