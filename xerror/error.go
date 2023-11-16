package xerror

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/xerrors"
)

func ExampleWrapError() {
	err1 := errors.New("this is error 1")
	err2 := fmt.Errorf("this is error 2: [%w]", err1)
	fmt.Println(err2)
	if errors.Unwrap(err2) == err1 {
		fmt.Println("err1 is err2's wrap")
	}
}

func exampleCallStack1() error {
	_, err := strconv.Atoi("sdf")
	return xerrors.Errorf("this is error 1: [%w]", err)
}

func exampleCallStack2() error {
	err := exampleCallStack1()
	return xerrors.Errorf("this is error 2: [%w]", err)
}

func exampleCallStack3() error {
	err := exampleCallStack2()
	return xerrors.Errorf("this is error 3: [%w]", err)
}

func ExampleCallStack() {
	err := exampleCallStack3()
	fmt.Printf("%+v\n", err)
}
