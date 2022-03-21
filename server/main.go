package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path"

	"go.uber.org/zap"
)

var logger *zap.Logger

const TRANSFER_DIR = "transferred"

func handleConnection(conn net.Conn) {
	var fileNameSz uint64
	var fileSz uint64
	var fileName string

	defer conn.Close()

	err := binary.Read(conn, binary.LittleEndian, &fileNameSz)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	buf := make([]byte, fileNameSz)
	n, err := io.ReadFull(conn, buf)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if n != int(fileNameSz) {
		logger.Error("Incorrect file name read. Aborting.")
		return
	}
	fileName = string(buf)

	err = binary.Read(conn, binary.LittleEndian, &fileSz)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Printf("FileNameLen: %d\nFileName: %s\nFileSize: %d\n", fileNameSz, fileName, fileSz)

	file, err := os.Create(path.Join(TRANSFER_DIR, fileName))
	defer file.Close()
	if err != nil {
		logger.Error("Could not create file: " + fileName)
		return
	}

	_, err = io.CopyN(file, conn, int64(fileSz))
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Printf("Successfully transferred %s.\n", fileName)

}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	listener, err := net.ListenTCP("tcp", addr)

	if err != nil {
		logger.Error(err.Error())
		return
	}
	_, err = os.Stat(TRANSFER_DIR)
	if os.IsNotExist(err) {
		err = os.Mkdir("transferred", os.ModeDir)
	}
	if err != nil {
		logger.Error(err.Error())
		return
	}

	fmt.Printf("Server started on port %d.\n", 8080)

	defer listener.Close()

	// Begin looping forever, waiting for the next connection
	for {
		conn, err := listener.Accept() // Wait for the next connection.
		if err != nil {
			logger.Error(err.Error())
			return
		}
		fmt.Printf("Client IP %s connected.\n", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}
