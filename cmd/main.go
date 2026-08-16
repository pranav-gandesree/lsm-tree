package main

import (
	"log"

	kv "github.com/pranav-gandesree/lsm-tree/internal"
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
		Name: "ned stark",
		Age:  40,
	}

	user2 := User{
		Name: "jamie lannister",
		Age:  35,
	}

	if err := store.PutData(7, user1); err != nil {
		log.Fatal(err)
	}

	if err := store.PutData(8, user2); err != nil {
		log.Fatal(err)
	}

	if err := store.DeleteData(5); err != nil {
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
