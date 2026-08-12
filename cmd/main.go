package main

import (
	"log"
	"lsm/internal/kv"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	store, err := kv.CreateStore[int, User]()

	if err != nil {
		log.Fatal(err)
	}

	store.PrintMap()

	user1 := User{
		Name: "katherine",
		Age:  500,
	}

	user2 := User{
		Name: "jon snow",
		Age:  30,
	}

	if err := store.PutData(3, user1); err != nil {
		log.Fatal(err)
	}

	if err := store.PutData(4, user2); err != nil {
		log.Fatal(err)
	}

	if err := store.DeleteData(1); err != nil {
		log.Fatal(err)
	}

	store.PrintMap()

	// // time.Sleep(time.Second)
	// user1, ok1 := store.GetData(0)
	// user2, ok2 := store.GetData(1)

	// if ok1 {
	// 	fmt.Println("Name:", user1.Name)
	// 	fmt.Println("Age:", user1.Age)
	// }

	// if ok2 {
	// 	fmt.Println("Name:", user2.Name)
	// 	fmt.Println("Age:", user2.Age)
	// }

	// var wg sync.WaitGroup

	// for i := 0; i < 100; i++ {
	// 	wg.Add(1)

	// 	go func(i int) {
	// 		defer wg.Done()
	// 		store.PutData(strconv.Itoa(i), "value")
	// 	}(i)
	// }

	// for i := 0; i < 100; i++ {
	// 	key := strconv.Itoa(i)

	// 	value, ok := store.GetData(key)
	// 	if ok {
	// 		fmt.Println(key, value)
	// 	}
	// }

	// wg.Wait()
}
