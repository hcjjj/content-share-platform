package main

func main() {
	res := DeferReturn()
	println(res)
	println(DeferReturnV1())
	println(DeferReturnV2().name)
}
