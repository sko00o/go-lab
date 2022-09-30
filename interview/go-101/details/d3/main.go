package main

func main() {
}

/*
The type deduction rule for the left untyped operand of a
bit-shift operation depends on whether or not the right
operand is a constant.
*/

const M = 2

// Compiles okay. 1.0 is deduced as an int value.
var _ = 1.0 << M // 使用 const 关键字初始化的参数作位运算， 参数被认定为 int 类型

var N = 2

// Fails to compile. 1.0 is deduced as a float64 value.
// var _ = 1.0 << N  // 使用 var 关键字初始化的参数作位运算， 参数被认定为 float64 类型
