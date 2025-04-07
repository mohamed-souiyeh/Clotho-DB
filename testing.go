package main

import (
	filemanager "Clotho/file_manager"
	"Clotho/file_manager/page"
	"fmt"
)

func main() {

	// data := make([]byte, 42)
	// buff := make([]byte, 200)

	// count, err := binary.Encode(buff, binary.NativeEndian, []byte("helloworld"))

	// if err != nil {
	// 	print("\nerr: ", err.Error())
	// 	print("\ncount: ", count)
	// }

	// copy(data[13:], buff)

	// var val int32 = 42
	// binary.Encode(buff, binary.NativeEndian, val)

	// print("\nsize: ", binary.Size([]byte("helloworld")))

	// print("copy count: ", copy(data[13 - binary.Size(val):13], buff))

	// clear(buff)

	// bytesCount, _ := binary.Decode(data[3:], binary.NativeEndian, buff)

	// print("\nbytes consumed: ", bytesCount, ", the decoded string: ", string(buff), ", expected: helloworld\n")

	// var val64 int32 = 1074563468

	// bytesCount, _ = binary.Decode(data[13 - binary.Size(val):], binary.NativeEndian, &val64)

	// print("bytes consumed: ", bytesCount, ", the decoded int: ",val64, ", expected: 42\n")

	// // fmt pkg

	// fmt.Printf("byte slice: %v\n", []byte{1, 2, 3, 4,5})

	// fmt.Printf("os.Getenv(\"HOME\"): %v\n", os.Getenv("HOME"))

	// dbDirectory, err := os.OpenRoot("./Makefile")

	// val, _ := err.(*os.PathError)

	// fmt.Printf("dbDirectory: %+v\nerr: %+v\n", dbDirectory, val)

	// err = os.MkdirAll("./makefile", 0755)

	// fmt.Printf("err: %+v\n", err)


	data := "hello world"

	p := page.NewPage(32)

	eSize, _ := page.EncodingSize(data)

	fmt.Println("encoding size of Data: ", eSize)

	p.SetString(data, 0)

	fmt.Println("page buffer len if data set at offset 5: ", len(p.Bytes()))






	fm, err := filemanager.NewFileManager("./makefile", 512)

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	if fm != nil {
		fmt.Printf("fm: %v\n", fm)
	}

	length, err := fm.Length("test")

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Println("test length (block number): ", length)

	fm.Extend("test")

	length, err = fm.Length("test")

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Println("test length (block number): ", length)
}

/*
  	JS
      r = function(pwd, fr, to, by) {
          return _.map(
            _.range(fr, to, by), function(i) { return pwd[i] }
          ).join('');
      }
      print = console.log

    PYTHON
      >>> def r(pwd, fr, to, by):
      ...     return ''.join([pwd[i] for i in range(fr, to, by)])

Output:

>>> print(r(pwd, 0, 4, 1)) -
2c45
/*
>>> print(r(pwd, 4, 20, 2)) -
39ee035f

/*
>>> print(r(pwd, 5, 30, 3)) -
7ea325b3e

/*
>>> print(r(pwd, 5, 40, 4)) -
7fa20ee0d

/*
>>> print(r(pwd, 6, 40, 6)) -
90f84d

/*
>>> print(r(pwd, 7, 8, 1)) -
1

/*
>>> print(r(pwd, 10, 40, 3)) -
ea527efd3d

/*
>>> print(r(pwd, 10, 40, 5)) -
e05e41

/*
>>> print(r(pwd, 15, 40, 6)) -
00105

/*
>>> print(r(pwd, 20, 40, 3))
5b3ee10
*/

// var memory [40]bool
// var count int = 0

// func decode(final []byte, str string, fr, to, by int) {
// 	for j, i := 0, fr; i < to; {
// 		if memory[i] == false {
// 			final[i] = str[j]
// 			memory[i] = true
// 		}
// 		i += by
// 		j++
// 		count++
// 	}
// }

// func main() {
// 	fragments := []struct {
// 		str string
// 		fr  int
// 		to  int
// 		by  int
// 	}{
// 		{
// 			/*
// 				>>> print(r(pwd, 0, 4, 1)) -
// 				2c45
// 			*/
// 			str: "2c45",
// 			fr:  0,
// 			to:  4,
// 			by:  1,
// 		},
// 		{
// 			/*
// 				>>> print(r(pwd, 4, 20, 2)) -
// 				39ee035f
// 			*/
// 			str: "39ee035f",
// 			fr:  4,
// 			to:  20,
// 			by:  2,
// 		},
// 		{
// 			/*
// 				>>> print(r(pwd, 5, 30, 3)) -
// 				7ea325b3e
// 			*/
// 			str: "7ea325b3e",
// 			fr:  5,
// 			to:  30,
// 			by:  3,
// 		},
// 		{
// 			/*
// 				>>> print(r(pwd, 5, 40, 4))
// 				7fa20ee0d
// 			*/
// 			str: "7fa20ee0d",
// 			fr:  5,
// 			to:  40,
// 			by:  4,
// 		},
// 		{
// 			/*
// 				>>> print(r(pwd, 6, 40, 6))
// 				90f84d
// 			*/
// 			str: "90f84d",
// 			fr:  6,
// 			to:  40,
// 			by:  6,
// 		},
// 		//------------
// 		{
// 			/*
// 				>>> print(r(pwd, 7, 8, 1))
// 				1
// 			*/
// 			str: "1",
// 			fr:  7,
// 			to:  8,
// 			by:  1,
// 		},
// 		{
// 			/*
// 				>>> print(r(pwd, 10, 40, 3))
// 				ea527efd3d
// 			*/
// 			str: "ea527efd3d",
// 			fr:  10,
// 			to:  40,
// 			by:  3,
// 		},
// 		{
// 			/*
// 				>>> print(r(pwd, 10, 40, 5))
// 				e05e41
// 			*/
// 			str: "e05e41",
// 			fr:  10,
// 			to:  40,
// 			by:  5,
// 		},
// 		{
// 			/*
// 				>>> print(r(pwd, 15, 40, 6))
// 				00105
// 			*/
// 			str: "00105",
// 			fr:  15,
// 			to:  40,
// 			by:  6,
// 		},
// 		{
// 			/*
// 				>>> print(r(pwd, 20, 40, 3))
// 				5b3ee10
// 			*/
// 			str: "5b3ee10",
// 			fr:  20,
// 			to:  40,
// 			by:  3,
// 		},
// 	}

// 	var final []byte = make([]byte, 40, 40)

// 	for _, v := range fragments {
// 		decode(final, v.str, v.fr, v.to, v.by)
// 	}

// 	fmt.Printf("count: %v\n", count)
// 	fmt.Printf("string(final): %v\n", string(final))
// }
