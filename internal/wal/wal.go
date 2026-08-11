package wal

import (
	"fmt"
	"os"
)

func AppendData(data string) error {

	// f, err := os.Create("data/wal.log") //destroys the existing logs too

	f, err := os.OpenFile(
		"data/wal.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err // donot use log.fatal, it ends and doesnt run f.close
	}

	defer f.Close()

	_, err = f.WriteString(data + "\n")

	fmt.Println("done")

	return err
}
