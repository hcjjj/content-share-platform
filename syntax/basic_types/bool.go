package main

func Bool() {
	var a bool = true
	var b bool = false
	var c bool = a || b // true
	println(c)
	var d bool = a && b // false
	println(d)
	var e bool = !a
	println(e)
}
