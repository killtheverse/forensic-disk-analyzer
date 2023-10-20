Name: Rahul Kureel
ASU ID: 1225362476
Email: rkureel@asu.edu


This Go program is designed to analyze disk image files. It includes functions to analyze both Master Boot Record (MBR) and GUID Partition Table (GPT) disk structures. The program operates on a binary disk image file.

Here's a brief description of the code:

- The program reads the JSON data from the PartitionTypesJson variable in the init function. This JSON data likely contains partition type codes and their corresponding descriptions.

- The AnalyzeImage function is the entry point for analyzing disk images. It takes the path to a disk image file as an argument.

- The function analyzeMBRImage is used to analyze MBR disk structures within the image.

- The function analyzeGPTImage is used to analyze GPT disk structures within the image.

- The program reads the disk image file, reads specific sections of the image, and extracts information about the partitions, including partition type, starting LBA address, and ending LBA address.

- It converts the numeric values, such as LBA addresses, from binary representation to hexadecimal and decimal formats and prints the information in a human-readable format.

- The code uses various offsets and constants to locate specific information within the disk image, such as partition entries in MBR and GPT headers.
