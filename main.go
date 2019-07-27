package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

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
	fmt.Printf("Decode raw output is:\n%s\n", output)

	message := fmt.Sprintf("%T", u)
	fmt.Println("program runtime type:", message)

	cmdstr = fmt.Sprintf("echo %x | xxd -r -p | protoc --decode %s proto/myobject.proto", data, message)
	output = runshell(cmdstr)
	fmt.Printf("decode with runtime type output is:\n%s\n", output)

	cmdstr = "awk '$1 == \"package\" {print $2}' proto/myobject.proto"
	output = runshell(cmdstr)
	fmt.Printf("package name from shell is:\n%s\n", output)
	pkg := fmt.Sprintf("%s", output)
	println("check the pluspure package", pureCmdStringPlus(pkg))

	cmdstr = "awk '$1 == \"message\" {print $2}' proto/myobject.proto"
	output = runshell(cmdstr)
	fmt.Printf("message name from shell is:\n%s\n", output)
	messages := strings.Split(fmt.Sprintf("%s", output), "\n")
	message = pureCmdString(pkg) + "." + messages[0]
	fmt.Println("messageType=", message)
	cmdstr = fmt.Sprintf("echo %x | xxd -r -p | protoc --decode %s proto/myobject.proto", data, message)
	println("cmd =", cmdstr)
	output = runshell(cmdstr)
	fmt.Printf("the shell cmd output is:\n%s\n", output)

	hardCoreResult := hardcoreDecode("proto/myobject.proto", data)
	fmt.Printf("Hardcode Decode result=%s\n", hardCoreResult)

	cmdstr = "awk '$1 == \"message\" {print $2}' proto/cpdcp_control.proto"
	output = runshell(cmdstr)
	fmt.Printf("the shell cmd output is:\n%s\n", output)
	message = fmt.Sprintf("%s", output)
	fmt.Println("count=", strings.Count(message, "\n"))
	messages = strings.Split(message, "\n")
	fmt.Println("messages=", messages)

}
func pureCmdString(str string) string {
	return strings.Trim(strings.Trim(strings.Trim(strings.Trim(str, "\n"), "\r"), " "), ";")
}

func pureCmdStringPlus(str string) string {
	return strings.Replace(strings.Replace(strings.Replace(strings.Replace(str, "\n", "", -1), "\r", "", -1), " ", "", -1), ";", "", -1)
}
func filterPkg(proto string) string {
	cmdstr := fmt.Sprintf("awk '$1 == \"package\" {print $2}' %s", proto)
	output := runshell(cmdstr)
	return pureCmdStringPlus(fmt.Sprintf("%s", output))

}
func filterMessageTypes(proto string) []string {
	cmdstr := fmt.Sprintf("awk '$1 == \"message\" {print $2}' %s", proto)
	output := runshell(cmdstr)
	messages := strings.Split(fmt.Sprintf("%s", output), "\n")
	fmt.Printf("before filter return %s\n", messages)
	for i := 0; i < len(messages); {
		if messages[i] == "\n" {
			messages = append(messages[:i], messages[i+1:]...)
		} else {
			i++
		}
	}
	fmt.Printf("After filter return %s\n", messages)
	return messages
}

func hardcoreDecode(proto string, data []byte) []byte {
	var pkgMesg, cmdstr string
	pkg := filterPkg(proto)
	types := filterMessageTypes(proto)
	for k, message := range types {
		pkgMesg = pkg + "." + message
		fmt.Printf("decode the %v type %s", k, pkgMesg)

		cmdstr = fmt.Sprintf("echo %x | xxd -r -p | protoc --decode %s %s", data, pkgMesg, proto)
		fmt.Println("cmd =", cmdstr)
		cmd := exec.Command("sh", "-c", cmdstr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("DecodeFail on messageType", pkgMesg, "continue...")
			continue
		} else {
			return output
		}
	}

	//finally give a raw decode
	cmdstr = fmt.Sprintf("echo %x | xxd -r -p | protoc --decode_raw", data)
	return runshell(cmdstr)
}
func runshell(shell string) []byte {
	cmd := exec.Command("sh", "-c", shell)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return output
}
