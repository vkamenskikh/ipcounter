Application read a text file and calculate amount of rows and unique IPs
The file format:
  145.67.23.4
  8.34.5.23
  89.54.3.124
  89.54.3.124
  3.45.71.5
  ...

Usage:
go run . fileName [1]
- fileName - full file path with the name
- 1 - optional parameter to instruct applicaiton to use single thread for file reading, which is faster for HDD (vs to SSD)

Application uses multiple threads to read a file by blocks. Amount of thread is equal to amount of CPUs

