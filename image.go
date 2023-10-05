package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var PartitionType map[byte]string

func init() {
	err := json.Unmarshal(PartitionTypesJson, &PartitionType)
	if err != nil {
		log.Fatal(err)
	}
}

func AnalyzeImage(filepath string) error {
	buffer := make([]byte, 2*BUFFER_SIZE)

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Read(buffer)
	if err == io.EOF {
		return err
	} else if n < 2*BUFFER_SIZE {
		log.Fatalln("File is not large enough.")
	} else if err != nil {
		return err
	}


	partitionEntryOne := buffer[PARTITION_ENTRY_1_OFFSET:PARTITION_ENTRY_1_OFFSET+PARTITION_ENTRY_SIZE]
    partitionEntryTwo := buffer[PARTITION_ENTRY_1_OFFSET+PARTITION_ENTRY_SIZE:PARTITION_ENTRY_1_OFFSET+2*PARTITION_ENTRY_SIZE]
    partitionEntryThree := buffer[PARTITION_ENTRY_1_OFFSET+2*PARTITION_ENTRY_SIZE:PARTITION_ENTRY_1_OFFSET+3*PARTITION_ENTRY_SIZE]
    partitionEntryFour := buffer[PARTITION_ENTRY_1_OFFSET+3*PARTITION_ENTRY_SIZE:PARTITION_ENTRY_1_OFFSET+4*PARTITION_ENTRY_SIZE]
    if partitionEntryOne[PARTITION_TYPE_OFFSET] == 0x07 || 
    partitionEntryTwo[PARTITION_TYPE_OFFSET] == 0x07 ||
    partitionEntryThree[PARTITION_TYPE_OFFSET] == 0x07 ||
    partitionEntryFour[PARTITION_TYPE_OFFSET] == 0x07 {
		err := analyzeMBRImage(file)
		if err != nil {
			return err
		}
	} else if partitionEntryOne[PARTITION_TYPE_OFFSET] == 0xee {
		err := analyzeGPTImage(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func analyzeMBRImage(file *os.File) error {
	buffer := make([]byte, BUFFER_SIZE)

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	n, err := file.Read(buffer)
	if err == io.EOF {
		return err
	} else if n < BUFFER_SIZE {
		log.Fatalln("File is not large enough.")
	} else if err != nil {
		return err
	}

	var partitionLBAAddresses []uint32

	for partitionNo:=1; partitionNo<=4; partitionNo++ {
		partitionEntry := buffer[PARTITION_ENTRY_1_OFFSET + (partitionNo-1)*PARTITION_ENTRY_SIZE : PARTITION_ENTRY_1_OFFSET + partitionNo*PARTITION_ENTRY_SIZE]
		partitionTypeByte := partitionEntry[PARTITION_TYPE_OFFSET]
		if partitionTypeByte == 0 {
			continue
		}
		partitionLBAAddress := binary.LittleEndian.Uint32(partitionEntry[PARTITION_LBA_OFFSET : PARTITION_LBA_OFFSET+4])*512
		partitionLBAAddresses = append(partitionLBAAddresses, partitionLBAAddress)
		partitionSize := binary.LittleEndian.Uint32(partitionEntry[PARTITION_SIZE_OFFSET : PARTITION_SIZE_OFFSET+4])*512
		fmt.Printf("(%02x) %s %d %d\n", partitionTypeByte, PartitionType[partitionTypeByte], partitionLBAAddress, partitionSize)
	}

	for i, partitionLBAAddress := range partitionLBAAddresses {
		fmt.Printf("Partition number: %d\n", i+1)
		if _, err := file.Seek(int64(partitionLBAAddress), 0); err != nil {
			return err
		}
		_, err := file.Read(buffer)
		if err != nil {
			return err
		}
		fmt.Printf("First 16 bytes of boot record: ")
		for _, b := range buffer[0:16] {
			fmt.Printf("%02x ", b)
		}
		fmt.Printf("\n")
		fmt.Print("ASCII:                         ")
		for _, b := range buffer[0:16] {
			if b>=33 && b<=126 {
				fmt.Printf(" %c ", b)
			} else {
				fmt.Print(" . ")
			}
		}
		fmt.Printf("\n")
	}
	return nil
}

func analyzeGPTImage(file *os.File) error {
	println("GPT Image")
	return nil
}
