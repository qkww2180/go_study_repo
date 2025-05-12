package algorithm

//直接用递归求斐波那契数列的第n个数
func Fibonacci(n int) int {
	if n <= 1 {
		return n //凡是递归，一定要有终止条件，否则会进入无限循环
	}
	return Fibonacci(n-1) + Fibonacci(n-2) //递归调用函数自身
}

func f(n int, array []int) int {
	if array[n] != -1 {
		return array[n]
	} else {
		result := f(n-1, array) + f(n-2, array)
		array[n] = result
		return result
	}
}

// 自上而下，把大问题化解为小问题。把小问题的解存起来，避免重复计算。时间和空间复杂度都是O(N)
func FibonacciTopDown(n int) int {
	if n <= 1 {
		return n
	}

	array := make([]int, n+1)
	for i := 0; i < n+1; i++ {
		array[i] = -1
	}
	array[0] = 0
	array[1] = 1

	return f(n, array)
}

// 自下而上，先求小问题，再逐步求更大的问题。时间和空间复杂度都是O(N)
func FibonacciButtomUp(n int) int {
	array := make([]int, n+1)
	array[0] = 0
	array[1] = 1
	for i := 2; i <= n; i++ {
		array[i] = array[i-1] + array[i-2]
	}
	return array[n]
}

// 自下而上。O(1)的空间复杂度
func FibonacciButtomUp_WithSpaceO1(n int) int {
	if n <= 1 {
		return n
	}
	prev2 := 0
	prev1 := 1
	for i := 2; i <= n; i++ {
		s := prev1 + prev2
		prev2 = prev1
		prev1 = s
	}
	return prev1
}

// 上台阶
func Steps(n int) int {
	if n <= 2 {
		return n
	}
	return Steps(n-1) + Steps(n-2)
}
