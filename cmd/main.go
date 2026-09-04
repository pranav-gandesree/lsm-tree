package main

import (
	"log"

	kv "github.com/pranav-gandesree/lsm-tree/internal"
)

func main() {
	memtable, err := kv.CreateMemTable()

	if err != nil {
		log.Fatal(err)
	}

	memtable.PrintMap()

	if err := memtable.PutData("p", "hi"); err != nil {
		log.Fatal(err)
	}

	if err := memtable.PutData("r", "helloe"); err != nil {
		log.Fatal(err)
	}

	if err := memtable.DeleteData("r"); err != nil {
		log.Fatal(err)
	}

	memtable.PrintMap()

	// // time.Sleep(time.Second)
	// user1, ok1 := store.GetData("0")
	// user2, ok2 := store.GetData("1")

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
