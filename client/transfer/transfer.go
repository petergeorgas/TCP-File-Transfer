package transfer

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

var logger, err = zap.NewProduction()

func TransferFile(hostname, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(fileName)
	if err != nil {
		logger.Error("Could not open file: " + fileName)
		return
	}
	defer file.Close()

	baseFileName := filepath.Base(fileName)
	fmt.Println("FILEPATH: ", baseFileName)

	fmt.Println("Client Started!")
	conn, err := net.Dial("tcp", hostname)

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
