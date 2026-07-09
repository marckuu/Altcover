package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting program...")
	fmt.Println(os.Getenv("CONN"))

}
