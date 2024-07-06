package main

import (
	"github.com/go-fuego/fuego"
)

func main(){
       s := fuego.NewServer(
		fuego.WithAddr(":8080"),
	   )

 fuego.Get(s,"/",func(c fuego.ContextNoBody) (string, error) {
		return "Hello world",nil
	})
	s.Run()
}