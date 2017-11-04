package shapes

import (
	"fmt"
)

type VmFunction func(runtime *Runtime, stackAddr Address)

type Library struct {
	Functions []VmFunction
	index     map[string]int
}

func (lib *Library) AddLibrary(other *Library) {
	for name, index := range other.index {
		lib.AddFunction(name, other.Functions[index])
	}
}

func (lib *Library) AddFunction(name string, vmFunc VmFunction) {
	if lib.index == nil {
		lib.index = map[string]int{}
	}

	index := len(lib.Functions)
	lib.Functions = append(lib.Functions, vmFunc)
	lib.index[name] = index
}

func (lib *Library) GetFunctionIndex(name string) (int, error) {
	index, ok := lib.index[name]

	if ok {
		return index, nil
	}

	return 0, fmt.Errorf("Unknown function '%s'", name)
}

func (lib *Library) GetFunction(name string) (VmFunction, error) {
	index, err := lib.GetFunctionIndex(name)

	if err != nil {
		return nil, err
	}

	return lib.Functions[index], nil
}

var __STD_LIB *Library

func StdLib() *Library {
	if __STD_LIB == nil {
		__STD_LIB = &Library{}
	}

	return __STD_LIB
}
