The application reads a text file and calculates amount of rows and unique IPs
The file format:<pre>
  145.67.23.4
  8.34.5.23
  89.54.3.124
  89.54.3.124
  3.45.71.5
  ...
</pre>
Usage:
go run . fileName [1]
- fileName - full file path with the name
- 1 - an optional parameter to instruct the applicaiton to use a single thread for a file reading, which is faster for HDD (vs to SSD)

The application uses multiple threads to read a file by blocks. The number of threads is equal to the number of CPUs.

