package main

import (
	"fmt"
	"log"
	"os/exec"

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

	//check the exec runtime result
	cmdstr := fmt.Sprintf("echo %x | xxd -r -p | protoc --decode_raw", data)
	output := runshell(cmdstr)
	fmt.Printf("the shell cmd output is:\n%s\n", output)
	message := fmt.Sprintf("%T", u)
	fmt.Println("fmt type:", message)
	cmdstr = fmt.Sprintf("echo %x | xxd -r -p | protoc --decode %s proto/myobject.proto", data, message)
	output = runshell(cmdstr)
	fmt.Printf("the shell cmd output is:\n%s\n", output)
}
func runshell(shell string) []byte {
	cmd := exec.Command("sh", "-c", shell)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return output
}
