package util

import (
	"fmt"
	"runtime"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func LogAndForget(err error) {
	if err != nil {
		fmt.Printf("Error in %s: %s\n", GetCurrentFuncName(2), err.Error())
	}
}

func LogError(err error) error {
	if err != nil {
		fmt.Printf("Error in %s: %s\n", GetCurrentFuncName(2), err.Error())
	}
	return err
}

// GetCurrentFuncName : get name of function being called
func GetCurrentFuncName(numCallStack int) string {
	pc, _, _, _ := runtime.Caller(numCallStack)
	return fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
}
