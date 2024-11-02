package main

import (
	"rgr/controller"
)

func main() {
	c, err := controller.New()
	if err != nil {
		panic(err.Error())
	}
	defer c.Destroy()
	
	for c.Index() {
	}

}
