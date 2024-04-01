// this is go file
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {



	// array()
	// slice()
	// add(1, 12)
	// fmt.Println(add(1, 2))
	// a, _ := add(1, 2)
	// fmt.Println(a)

	// fmt.Println(testcount(1))
	// fmt.Println(fact(4))

	// structs()

	// mapss()

	//    r:=rect{wight: 10,height : 20}
	//    fmt.Println("area",r.area())
	//    fmt.Println("Perimeter",r.permi())

	//    rp:= &r
	//    fmt.Println("Area",rp.area())
	//    fmt.Println("Permieter",rp.permi())
	// panics()
	files()

}
func array() {
	var arr1 = [3]int{1, 2}
	arr2 := [5]int{5, 6, 7, 8, 9}

	var arr3 = [...]int{6, 5, 8}
	arr4 := [...]int{5, 5, 8, 8, 8}
	twoD := make([]int, 3)

	fmt.Println(twoD)
	fmt.Println(arr1, arr2)
	fmt.Println(len(arr1))
	fmt.Println(cap(arr1))
	fmt.Println(arr3, arr4)
}
func slice() {
	var s []string
	// fmt.Println("Slices",s,s==nil,len(s)==0)

	s = make([]string, 3)
	// fmt.Println(s,len(s),cap(s))

	myslice1 := []int{1, 2, 3}
	myslice2 := []int{4, 5, 6}
	myslice3 := append(myslice1, myslice2...)
	fmt.Println(myslice3)

	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	s = append(s, "d")
	fmt.Println("S", s)

	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("Copy", c)

	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)

}
func add(a int, b int) (result int, txt string) {
	result = a + b
	txt = "Hello!kjrgndkjgn"
	return
}
func testcount(x int) int {
	if x == 11 {
		return 0
	}
	fmt.Println(x)
	return testcount(x + 1)
}
func fact(x float64) (y float64) {
	if x > 0 {
		y = x * fact(x-1)
	} else {
		y = 1
	}
	return
}

func structs() {

	type Employee struct {
		id     int
		name   string
		age    int
		PhNo   int
		job    string
		salary int
	}

	var e1 Employee
	var e2 Employee

	e1.id = 10
	e1.name = "Shashidhar"
	e1.age = 21
	e1.PhNo = 1234567890
	e1.job = "Monitoring"
	e1.salary = 10000
	e2.id = 11
	e2.name = "Asheaf"
	e2.age = 90
	e2.PhNo = 544545452
	e2.job = "Monitoring"
	e2.salary = 10000

	fmt.Println(e1)
	fmt.Println(e2)

	fmt.Println(e1.name)
	fmt.Println(e2.name)

}
func mapss() {

	m := make(map[string]int)

	m["k1"] = 2
	m["k2"] = 3

	fmt.Println(m)
	delete(m, "k1")
	fmt.Println(m)
	clear(m)
	fmt.Println(m)

	var a = make(map[string]int)
	a["kw"] = 12
	a["gd"] = 13
	fmt.Println(a)

}

type rect struct {
	wight, height int
}

func (r *rect) area() int {
	return r.wight * r.height
}
func (r rect) permi() int {
	return 2*r.height + 2*r.wight
}
func panics() {
	fmt.Printf("Before panic")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	panic("Something went wroung")

	// fmt.Println("After the panic")
}
func files() {

	file, err := os.Open("example.txt")
	if err != nil {
		log.Fatalf("Error opening file :%v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	fmt.Println("File content (ReadAll)")
	fmt.Println(string(data))

	file.Seek(0, 0)

	fmt.Println("File Content (ReadAll):")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Printf(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning file: %v", err)
	}

	// Reset file cursor to the beginning
	file.Seek(0, 0)

	// Buffered reading
	fmt.Println("\nFile Content (Buffered Reading):")
	bufferedReader := bufio.NewReader(file)
	for {
		line, err := bufferedReader.ReadString('\n')
		if err != nil {
			break // End of file
		}
		fmt.Print(line)
	}

	// Close the file
	if err := file.Close(); err != nil {
		log.Fatalf("Error closing file: %v", err)
	}

	

}

