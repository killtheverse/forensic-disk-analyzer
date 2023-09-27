package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const BUFFER_SIZE = 1024

func calculateHash(file *os.File, hasher hash.Hash) (string, error) {
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

func storeHashes(file_path string) error {
	file, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer file.Close()

	md5Hasher := md5.New()
	md5Hash, err := calculateHash(file, md5Hasher)
	if err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	shaHasher := sha256.New()
	shaHash, err := calculateHash(file, shaHasher)
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

func main() {
	app := cli.NewApp()
	app.Name = "boot_info"
	app.Usage = "Analyzes the MBR/GPT information of forensic images"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "filepath",
			Aliases: []string{"f"},
			Usage: "`path` of the forensic image file",
			Required: true,
		},
	}
	app.Action = func(ctx *cli.Context) error {
		println(ctx.String("filepath"))
		storeHashes(ctx.String("filepath"))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}