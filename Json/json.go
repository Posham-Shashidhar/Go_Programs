package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
}

func main() {
	person := Person{
		Name:    "Shashidhar",
		Age:     23,
		Country: "India",
	}

	jsonData, err := json.MarshalIndent(person, "", "    ")
	if err != nil {
		fmt.Println("Error in marshalling", err)
		return
	}
	file, err := os.Create("person.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {

		fmt.Println("error in writing the file", err)

	}
	fmt.Println("JSON data successfully ")

	// file1, err1 := os.Open("person.json")
	// if err1 != nil {
	// 	fmt.Println("Error opening file",err)
	// 	return 
	// }
    // defer file1.Close()
	// var people []Person
	// err1 = json.NewDecoder(file1).Decode(&people)
	// if err1 != nil {
	// 	fmt.Println("Encode JSON file error :", err)
	// 	return
	// }

	// fmt.Println("Decoded data of JSON:")
	// for _, p := range people {
	// 	fmt.Printf("Name : %s , Age : %d , Country : %s\n", p.Name, p.Age, p.Country)
	// }

	printJSONData("person.json")

}
func printJSONData(filename string) error {
    jsonData, err := ioutil.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("error reading file: %v", err)
    }


    var person Person
    err = json.Unmarshal(jsonData, &person)
    if err != nil {
        return fmt.Errorf("error parsing JSON: %v", err)
    }

    fmt.Println("Name:", person.Name)
    fmt.Println("Age:", person.Age)
    fmt.Println("Country:", person.Country)

    return nil
}