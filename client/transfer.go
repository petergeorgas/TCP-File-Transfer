package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: transfer <file>\n")
		os.Exit(2)
	}

	logger, _ := zap.NewProduction()

	fileName := args[1]
	file, err := os.Open(fileName)
	if err != nil {
		logger.Error("Could not open file: " + fileName)
		return
	}
	baseFileName := filepath.Base(fileName)
	fmt.Println("FILEPATH: ", baseFileName)

	defer file.Close()

	fmt.Println("Client Started!")
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer conn.Close()
	if fileName == "" {
		logger.Error("No file name found. Cannot continue.")
	}

	fileStats, _ := file.Stat()
	fileSz := fileStats.Size()

	err = binary.Write(conn, binary.LittleEndian, uint64(len(baseFileName)))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	n, err := io.WriteString(conn, baseFileName)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if n != len(baseFileName) {
		logger.Error("Incomplete file name sent. Aborting.")
		return
	}
	err = binary.Write(conn, binary.LittleEndian, uint64(fileSz))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	fmt.Println("File Info Sent...")
	bytesCopied, err := io.CopyN(conn, file, fileSz)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if bytesCopied != fileSz {
		logger.Error("Did not complete file transfer. Aborted.")
		return
	}
	fmt.Printf("%d B copied.\n", bytesCopied)
}
