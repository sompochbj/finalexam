package main

import "github.com/sompochbj/finalexam/customer"

func main() {
	r := customer.SetupRouter()
	r.Run(":2019")
}
