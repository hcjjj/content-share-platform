package main

import "unicode/utf8"

func String() {
	// He said:"hello, go!"
	println("He said:\"hello, go!\"")
	println(`我可以换行
这是新的行
但是这里不能有反引号
`)

	println("hello, " + "Go!")
	// 字符串只能和字符串拼接
	//println("hello, " + 123)

	// strings 包里面放着各种字符串相关操作的方法，需要的时候再查阅
	//strings.

	println(len("你好"))                      // 输出 6
	println(utf8.RuneCountInString("你好"))   // 输出2
	println(utf8.RuneCountInString("你好ab")) // 输出4
}
