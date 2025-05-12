package algorithm

import "cmp"

//递归的方式求最长公共子序列的长度
func LCS[T cmp.Ordered](s1, s2 []T) int {
	i, j := len(s1), len(s2)
	if i == 0 || j == 0 {
		return 0
	}
	if s1[i-1] == s2[j-1] {
		return 1 + LCS(s1[:i-1], s2[:j-1])
	} else {
		max := max(LCS(s1[:i-1], s2[:j]), LCS(s1[:i], s2[:j-1]))
		return max
	}
}

func lcs[T cmp.Ordered](s1, s2 []T, arr [][]int) int {
	i := len(s1)
	j := len(s2)
	if arr[i][j] != -1 {
		return arr[i][j]
	} else {
		var length int
		if s1[i-1] == s2[j-1] {
			length = 1 + lcs(s1[0:i-1], s2[0:j-1], arr)
		} else {
			length = max(lcs(s1[:i-1], s2[:j], arr), lcs(s1[:i], s2[:j-1], arr))
		}
		arr[i][j] = length
		return length
	}
}

// 最长公共子序列。自顶向下，仍然用到递归，但是借助于arr避免了重复计算
func LCSTopDown[T cmp.Ordered](s1, s2 []T) int {
	i := len(s1)
	j := len(s2)
	//初始化arr是一个i+1行，j+1列的二维切片
	arr := make([][]int, i+1)
	for row := range arr {
		arr[row] = make([]int, j+1)
	}
	//把arr的第一行和第一列全赋为0，其他位置全部赋为-1
	for m := 0; m < i+1; m++ {
		for n := 0; n < j+1; n++ {
			if m == 0 {
				arr[m][n] = 0
			} else if n == 0 {
				arr[m][n] = 0
			} else {
				arr[m][n] = -1
			}
		}
	}
	return lcs(s1, s2, arr)
}

// 最长公共子序列。自下而上，没有用递归，空间开销O(M*N)
func LCSButtomUp[T cmp.Ordered](s1, s2 []T) int {
	i := len(s1)
	j := len(s2)
	//初始化arr是一个i+1行，j+1列的二维切片
	arr := make([][]int, i+1)
	for row := range arr {
		arr[row] = make([]int, j+1) //初始全为0
	}
	for m := 0; m < i+1; m++ {
		for n := 0; n < j+1; n++ {
			if m == 0 { //首行赋为0
				arr[m][n] = 0
			} else if n == 0 { //首列赋为0
				arr[m][n] = 0
			} else {
				if s1[m-1] == s2[n-1] {
					arr[m][n] = 1 + arr[m-1][n-1]
				} else {
					max := max(arr[m-1][n], arr[m][n-1])
					arr[m][n] = max
				}
			}
		}
	}
	return arr[i][j]
}

// 最长公共子序列。自下而上，没有用递归，空间开销O(N)
func LCSButtomUp_WithSpaceON[T cmp.Ordered](s1, s2 []T) int {
	//使s2成为更短的那个，节省空间开销
	if len(s1) < len(s2) {
		s1, s2 = s2, s1
	}

	i := len(s1)
	j := len(s2)
	//只维护最下面的那两行
	prevRow := make([]int, j+1) //初始全为0
	currRow := make([]int, j+1)
	for m := 1; m < i+1; m++ {
		//填充currRow
		for n := 1; n < j+1; n++ {
			if s1[m-1] == s2[n-1] {
				currRow[n] = 1 + prevRow[n-1]
			} else {
				max := max(prevRow[n], currRow[n-1])
				currRow[n] = max
			}
		}
		//把currRow拷贝给prevRow
		// for n, ele := range currRow {
		// 	prevRow[n] = ele
		// }
		copy(prevRow, currRow) //go标准库的copy作了性能优化
	}
	return currRow[j]
}
