package main

const (
	BUFFER_SIZE = 512
	PARTITION_ENTRY_SIZE = 0x10
	PARTITION_ENTRY_1_OFFSET = 0x01be
	BOOT_SIGNATURE = 0x01fe
	PARTITION_ENTRY_FLAG_OFFSET = 0x00
	PARTITION_TYPE_OFFSET = 0x04
	PARTITION_LBA_OFFSET = 0x08
	PARTITION_SIZE_OFFSET = 0x0c
	GPT_HEADER_OFFSET = 0x200
	GPT_HEADER_PARTITION_ENTRIES_STARTING_LBA_OFFSET = 0x48
	PARTITION_ENTRY_START_LBA_OFFSET = 0x20
	PARTITION_ENTRY_END_LBA_OFFSET = 0x28
)
