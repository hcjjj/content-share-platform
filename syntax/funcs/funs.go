package main

// Func0 单一返回值
func Func0(name string) string {
	return "hello, world"
}

// Func1 多个返回值
func Func1(a, b, c int, str1 string) (string, error) {
	return "", nil
}

// Func2 带名字的返回值
func Func2(a int, b int) (str string, err error) {
	str = "hello"
	// 带名字的返回值，可以直接 return
	return
}

// Func3 带名字的返回值
func Func3(a int, b int) (str string, err error) {
	res := "hello"
	// 虽然带名字，但是我们并没有用
	return res, nil
}

func Invoke() {
	str, err := Func2(1, 2)
	println(str, err)
	// 忽略返回值
	_, _ = Func2(1, 3)

	// 部分忽略返回值
	// str 是已经声明好了
	str, _ = Func2(3, 4)
	// str1 是新变量，需要使用 :=
	str1, _ := Func2(3, 4)
	println(str1)

	// str2 是新变量，需要使用 :=
	str2, err := Func2(3, 4)
	println(str2)
}

// Recursive 递归
// 这个方法运行的时候会出现错误
func Recursive() {
	Recursive()
}

func Func4() {
	myFunc3 := Func3
	_, _ = myFunc3(1, 2)
}

func Func5() {
	fn := func(name string) string {
		return "hello, " + name
	}

	fn("hcjjj")
}

// Func6 的返回值是一个方法，
func Func6() func(name string) string {
	return func(name string) string {
		return "hello," + name
	}
}

func Func6Invoke() {
	// sayHello 的签名就是 func(name string) string
	sayHello := Func6()
	sayHello("hcjjj")
}

// Func7 演示匿名方法
func Func7() {
	hello := func() string {
		return "hello, world"
	}()
	println(hello)
}

func Closure(name string) func() string {

	// 返回的这个函数，就是一个闭包。
	// 它引用到了 Closure 这个方法的入参
	return func() string {
		return "hello, " + name
	}
}

// YourName 不定参数例子
// 一个人可能有很多个别名，也可能一个都没有
func YourName(name string, alias ...string) {
	if len(alias) > 0 {
		println(alias[0])
	}
}

func YourNameInvoke() {
	YourName("hcjjj")
	YourName("hcjjj", "hcj")
	YourName("hcjjj", "hcj", "h")
}

func Defer() {
	defer func() {
		println("第一个 defer")
	}()

	defer func() {
		println("第二个 defer")
	}()
}

func DeferNum() {
	for i := 0; i < 100; i++ {
		defer func() {
			println("hello, world")
		}()
	}
}

func DeferOpenCode() {
	defer func() {
		println("hello, world")
	}()

	// 假如说这个是你的业务代码
	YourBusiness()

	// 开放编码就类似于编译器帮你把 defer 的内容放过来了这里
	//func() {
	//	println("hello, world")
	//}()
}

func YourBusiness() {
	// 假如说这是你的业务代码
}

func DeferConditional(input bool) {
	if input {
		defer func() {
			println("hello, world")
		}()
	}
}

func DeferClosure() {
	i := 0
	defer func() {
		println(i)
	}()
	i = 1
}

func DeferClosureV1() {
	i := 0
	defer func(val int) {
		println(val)
	}(i)
	i = 1
}

func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}

func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}

func DeferClosureLoopV3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}

func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturnV1() (a int) {
	a = 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturnV2() *MyStruct {
	a := &MyStruct{
		name: "Jerry",
	}
	defer func() {
		a.name = "Tom"
	}()

	return a
}

type MyStruct struct {
	name string
}
