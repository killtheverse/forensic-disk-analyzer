package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var PartitionType map[byte]string

func init() {
	err := json.Unmarshal(PartitionTypesJson, &PartitionType)
	if err != nil {
		fmt.Println(err)
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
		return fmt.Errorf("file is not large enough")
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
		return fmt.Errorf("file is not large enough")
	} else if err != nil {
		return err
	}

	var partitionLBAAddresses []uint32

	for partitionNo:=1; partitionNo<=4; partitionNo++ {
		partitionEntry := buffer[PARTITION_ENTRY_1_OFFSET + (partitionNo-1)*PARTITION_ENTRY_SIZE : PARTITION_ENTRY_1_OFFSET + partitionNo*PARTITION_ENTRY_SIZE]
		partitionTypeByte := partitionEntry[PARTITION_TYPE_OFFSET]
		if partitionTypeByte == 0 {
			fmt.Printf("(00) Empty, 0, 0\n")
			continue
		}
		partitionLBA := binary.LittleEndian.Uint32(partitionEntry[PARTITION_LBA_OFFSET : PARTITION_LBA_OFFSET+4])
		partitionLBAAddresses = append(partitionLBAAddresses, partitionLBA*512)
		partitionSize := binary.LittleEndian.Uint32(partitionEntry[PARTITION_SIZE_OFFSET : PARTITION_SIZE_OFFSET+4])*512
		fmt.Printf("(%02x) %s, %d, %d\n", partitionTypeByte, PartitionType[partitionTypeByte], partitionLBA, partitionSize)
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
	buffer := make([]byte, BUFFER_SIZE)

	if _, err := file.Seek(GPT_HEADER_OFFSET, 0); err != nil {
		return err
	}
	n, err := file.Read(buffer)
	if err == io.EOF {
		return err
	} else if n < BUFFER_SIZE {
		return fmt.Errorf("file is not large enough")
	} else if err != nil {
		return err
	}

	partitionStartingLBA := binary.LittleEndian.Uint64(buffer[GPT_HEADER_PARTITION_ENTRIES_STARTING_LBA_OFFSET:GPT_HEADER_PARTITION_ENTRIES_STARTING_LBA_OFFSET+8])
	partitionEntry := make([]byte, 128)
	partitionEntryAddress := partitionStartingLBA*512
	for partitionNo:=1; partitionNo<=128; partitionNo+=1 {
		if _, err := file.Seek(int64(partitionEntryAddress), 0); err != nil {
			return err
		}
		_, err := file.Read(partitionEntry)
		if err != nil {
			return err
		}

		skipEntry := true
		for _, b := range(partitionEntry) {
			if b != 0 {
				skipEntry = false
			}
		}
		if skipEntry {
			continue
		}
		
		startingLBA := binary.LittleEndian.Uint64(partitionEntry[PARTITION_ENTRY_START_LBA_OFFSET:PARTITION_ENTRY_START_LBA_OFFSET+8])
		endingLBA := binary.LittleEndian.Uint64(partitionEntry[PARTITION_ENTRY_END_LBA_OFFSET:PARTITION_ENTRY_END_LBA_OFFSET+8])

		fmt.Printf("Partition number: %d\n", partitionNo)
		fmt.Printf("Partition Type GUID : ")
		for _, b := range(partitionEntry[0:16]) {
			fmt.Printf("%02x", b)
		}
		fmt.Printf("\n")
		fmt.Printf("Starting LBA address in hex: 0x%x\n", startingLBA)
		fmt.Printf("ending LBA address in hex: 0x%x\n", endingLBA)
		fmt.Printf("starting LBA address in Decimal: %d\n", startingLBA)
		fmt.Printf("ending LBA address in Decimal: %d\n", endingLBA)
		fmt.Printf("\n")

		partitionEntryAddress += 128
	}

	return nil
}
