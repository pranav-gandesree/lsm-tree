package main

import (
	"fmt"
	"log"
	"lsm/internal/kv"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	store := kv.CreateStore[int, User]()

	user1 := User{
		Name: "stefan salvatore",
		Age:  160,
	}

	user2 := User{
		Name: "jon snow",
		Age:  30,
	}

	if err := store.PutData(0, user1); err != nil {
		log.Fatal(err)
	}

	if err := store.PutData(1, user2); err != nil {
		log.Fatal(err)
	}

	if err := store.DeleteData(1); err != nil {
		log.Fatal(err)
	}

	// time.Sleep(time.Second)
	user1, ok1 := store.GetData(0)
	user2, ok2 := store.GetData(1)

	if ok1 {
		fmt.Println("Name:", user1.Name)
		fmt.Println("Age:", user1.Age)
	}

	if ok2 {
		fmt.Println("Name:", user2.Name)
		fmt.Println("Age:", user2.Age)
	}

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
