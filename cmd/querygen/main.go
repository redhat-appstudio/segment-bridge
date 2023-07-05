package main

import (
	"fmt"

	"github.com/redhat-appstudio/segment-bridge.git/querygen"
)

func main() {
	fmt.Println(querygen.GenQuery())
}
