package main

func Byte() {
	var a byte = 'a'
	// 输出的是的 a 的ASC II 表达 97
	println(a)

	var str string = "this is string"
	var bs []byte = []byte(str)
	println(bs)
}
