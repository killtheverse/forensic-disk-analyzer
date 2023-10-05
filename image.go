package main

import (
	"io"
	"log"
	"os"
)

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
		err := analyzeMBRImage(buffer)
		if err != nil {
			return err
		}
	} else if partitionEntryOne[PARTITION_TYPE_OFFSET] == 0xee {
		err := analyzeGPTImage(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func analyzeMBRImage(buffer []byte) error {
	println("MBR Image")
	return nil
}

func analyzeGPTImage(buffer []byte) error {
	println("GPT Image")
	return nil
}
