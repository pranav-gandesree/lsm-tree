package main

import (
	"fmt"
	"lsm/internal/kv"
	"strconv"
	"sync"
)

func main() {
	store := kv.CreateStore()

	// go store.PutData("name", "pranav")
	// go store.PutData("age", "21")

	// time.Sleep(time.Second)

	// name, ok := store.GetData("name")
	// age, ok := store.GetData("age")
	// if ok {
	// 	fmt.Println("Name:", name)
	// 	fmt.Println("Age", age)
	// }

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			store.PutData(strconv.Itoa(i), "value")
		}(i)
	}

	for i := 0; i < 100; i++ {
		key := strconv.Itoa(i)

		value, ok := store.GetData(key)
		if ok {
			fmt.Println(key, value)
		}
	}

	wg.Wait()
}
