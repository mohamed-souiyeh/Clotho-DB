package main

import "encoding/binary"



func main () {

	data := make([]byte, 42)
	buff := make([]byte, 200)
	
	count, err := binary.Encode(buff, binary.NativeEndian, []byte("helloworld"))

	if err != nil {
		print("\nerr: ", err.Error())
		print("\ncount: ", count)
	}

	copy(data[13:], buff)

	var val int32 = 42
	binary.Encode(buff, binary.NativeEndian, val)

	print("\nsize: ", binary.Size([]byte("helloworld")))

	print("copy count: ", copy(data[13 - binary.Size(val):13], buff))
	
	clear(buff)

	bytesCount, _ := binary.Decode(data[3:], binary.NativeEndian, buff)

	print("\nbytes consumed: ", bytesCount, ", the decoded string: ", string(buff), ", expected: helloworld\n")
	
	var val64 int32 = 1074563468

	bytesCount, _ = binary.Decode(data[13 - binary.Size(val):], binary.NativeEndian, &val64)

	print("bytes consumed: ", bytesCount, ", the decoded int: ",val64, ", expected: 42")


}