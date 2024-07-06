package main

import (
	"github.com/go-fuego/fuego"
)

type MyInput struct {
	Name string `json:"name" validate:"required"`
}

type MyOutput struct {
	Message string `json:"message"`
}

func main(){
   s := fuego.NewServer(
	fuego.WithAddr(":8080"),
   )
   //here we are sending tha data in body so that we are giving context with body
    fuego.Post(s,"/",func(c *fuego.ContextWithBody[MyInput]) (MyOutput,error) {
     body, err := c.Body()
	 if err!=nil {
		return MyOutput{},err
	 }
       
	  return MyOutput{
		Message: "Hello" + body.Name,
	  },nil

	})
	s.Run()
}