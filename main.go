package main

import (
	"fmt"
	"log"

	"gpbdecoder/goproto/proto"

	"github.com/golang/protobuf/proto"
)

func main() {
	fmt.Println("hello GPB")
	u := myobject.User{
		Id:    proto.Int64(1),
		Name:  proto.String("Mike"),
		Email: proto.String("12345@html.com"),
	}
	data, err := proto.Marshal(&u)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	fmt.Printf("Before Marshal data:%v\n", u)
	fmt.Printf("After Marshal data:%x\n", data)
	ur := &myobject.User{}
	err = proto.Unmarshal(data, ur)
	fmt.Printf("After Unmarshal data:%v\n", ur)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
}
