package main

import (
	"log"

	"github.com/xmdhs/hidethread/get"
)

func main() {
	log.Println("开始")
	get.Range(0, 1092244, 5)
	log.Println("结束")
}
