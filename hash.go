package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

func CalculateHash(file *os.File, hasher hash.Hash) (string, error) {
	buffer := make([]byte, BUFFER_SIZE)
	hasher.Reset()
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		hasher.Write(buffer[:n])
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func StoreHashes(file_path string) error {
	file, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer file.Close()

	md5Hasher := md5.New()
	md5Hash, err := CalculateHash(file, md5Hasher)
	if err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	shaHasher := sha256.New()
	shaHash, err := CalculateHash(file, shaHasher)
	if err != nil {
		return err
	}

	filename := filepath.Base(file_path)
	
	md5HashFile, err := os.Create("MD5-" + filename + ".txt")
	if err != nil {
		return err
	}
	defer md5HashFile.Close()

	_, err = md5HashFile.Write([]byte(md5Hash))
	if err != nil {
		return err
	}

	sha256HashFile, err := os.Create("SHA-256-" + filename + ".txt")
	if err != nil {
		return err
	}

	_, err = sha256HashFile.Write([]byte(shaHash))
	if err != nil {
		return err
	}

	return nil
}