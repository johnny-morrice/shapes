package shapes

var __STD_LIB *Library

func StdLib() *Library {
	if __STD_LIB == nil {
		__STD_LIB = &Library{}
	}

	return __STD_LIB
}
