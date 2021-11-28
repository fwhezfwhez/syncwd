package syncwd

import "fmt"

func Printf(f string, a ... interface{}) {
	fmt.Printf(fmt.Sprintf("syncwd info: %s \n", f), a...)
}

func Errorf(f string, a ... interface{}) {
	fmt.Printf(fmt.Sprintf("syncwd err: %s \n", f), a...)
}
