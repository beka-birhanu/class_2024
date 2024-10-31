package main

import "fmt"

func main4() {
	arr := [3]int{1, 2, 3}
	fmt.Println("Array:", arr)

	slice := []int{4, 5, 6}
	slice = append(slice, 7)
	fmt.Println("Slice:", slice)

	myMap := make(map[string]int)
	myMap["Alice"] = 25
	myMap["Bob"] = 30
	fmt.Println("Map:", myMap)
	fmt.Println("Alice's age:", myMap["Alice"])

	for i, v := range slice {
		fmt.Printf("Index: %d, Value: %d\n", i, v)
	}

	for key, value := range myMap {
		fmt.Printf("%s is %d years old\n", key, value)
	}
}
