package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var startDir = "./../.."
var checksum = make([]byte, 32)

func processDir(dir string) {
	items, _ := os.ReadDir(dir)
	for _, item := range items {
		if strings.HasPrefix(item.Name(), ".") {
			continue
		}
		fullPath := filepath.Join(dir, item.Name())

		if item.IsDir() {
			processDir(fullPath)
			continue
		}

		hash := sha256.New()
		relativeToStart, _ := filepath.Rel(startDir, fullPath)
		hash.Write([]byte(relativeToStart))
		fileContent, _ := os.ReadFile(fullPath)
		hash.Write(fileContent)
		itemHash := hash.Sum([]byte{})

		fmt.Printf("%s %s\n", hex.EncodeToString(itemHash), relativeToStart)

		xorItemHash(itemHash)
	}
}

func xorItemHash(itemHash []byte) {
	for i := 0; i < 32; i++ {
		checksum[i] = checksum[i] ^ itemHash[i]
	}
}

func main() {
	processDir(startDir)
	fmt.Printf("Checksum : %s", hex.EncodeToString(checksum))
}
