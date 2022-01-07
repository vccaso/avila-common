package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	fmt.Println("module avila-common")
	fmt.Println(errors.New("testing error"))
}
