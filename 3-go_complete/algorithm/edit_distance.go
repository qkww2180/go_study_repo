package algorithm

import "cmp"

// 用递归实现编辑距离
func EditDistance[T cmp.Ordered](s1, s2 []T) int {
	i := len(s1)
	j := len(s2)
	if i == 0 {
		return j
	}
	if j == 0 {
		return i
	}
	if s1[i-1] == s2[j-1] {
		return EditDistance(s1[0:i-1], s2[0:j-1])
	} else {
		subED := min(EditDistance(s1[0:i-1], s2[0:j-1]), EditDistance(s1[0:i-1], s2[0:j]), EditDistance(s1[0:i], s2[0:j-1]))
		return 1 + subED
	}
}

func ed[T cmp.Ordered](s1, s2 []T, arr [][]int) int {
	i := len(s1)
	j := len(s2)
	if arr[i][j] != -1 {
		return arr[i][j]
	} else {
		var distance int
		if s1[i-1] == s2[j-1] {
			distance = ed(s1[0:i-1], s2[0:j-1], arr)
		} else {
			subED := min(ed(s1[0:i-1], s2[0:j-1], arr), ed(s1[0:i-1], s2[0:j], arr), ed(s1[0:i], s2[0:j-1], arr))
			distance = 1 + subED
		}
		arr[i][j] = distance
		return distance
	}
}

// 编辑距离。自顶向下，仍然用到递归，但是借助于arr避免了重复计算
func EditDistanceTopDown[T cmp.Ordered](s1, s2 []T) int {
	i := len(s1)
	j := len(s2)
	//初始化arr是一个i+1行，j+1列的二维切片
	arr := make([][]int, i+1)
	for row := range arr {
		arr[row] = make([]int, j+1)
	}
	//把arr的第一行赋为0,1,2,3...，第一列赋为0,1,2,3...，其他位置全部赋为-1
	for m := 0; m < i+1; m++ {
		for n := 0; n < j+1; n++ {
			if m == 0 {
				arr[m][n] = n
			} else if n == 0 {
				arr[m][n] = m
			} else {
				arr[m][n] = -1
			}
		}
	}
	return ed(s1, s2, arr)
}

// 编辑距离。自下而上，没有用递归，空间开销O(M*N)
func EditDistanceButtomUp[T cmp.Ordered](s1, s2 []T) int {
	i := len(s1)
	j := len(s2)
	//初始化arr是一个i+1行，j+1列的二维切片
	arr := make([][]int, i+1)
	for row := range arr {
		arr[row] = make([]int, j+1)
	}
	for m := 0; m < i+1; m++ {
		for n := 0; n < j+1; n++ {
			if m == 0 { //首行赋为0,1,2,3...
				arr[m][n] = n
			} else if n == 0 { //首列赋为0,1,2,3...
				arr[m][n] = m
			} else {
				if s1[m-1] == s2[n-1] {
					arr[m][n] = arr[m-1][n-1]
				} else {
					dist := min(arr[m-1][n-1], arr[m-1][n], arr[m][n-1])
					arr[m][n] = 1 + dist
				}
			}
		}
	}
	return arr[i][j]
}

// 编辑距离。自下而上，空间开销2N
func EditDistanceButtomUp_WithSpaceON[T cmp.Ordered](s1, s2 []T) int {
	//使s2成为更短的那个，节省空间开销
	if len(s1) < len(s2) {
		s1, s2 = s2, s1
	}

	i := len(s1)
	j := len(s2)
	//只维护最下面的那两行
	prevRow := make([]int, j+1)
	currRow := make([]int, j+1)
	for n := 0; n < j+1; n++ { //初始化prevRow
		prevRow[n] = n
	}

	for m := 1; m < i+1; m++ {
		//填充currRow
		for n := 0; n < j+1; n++ {
			if n == 0 {
				currRow[n] = m
			} else {
				if s1[m-1] == s2[n-1] {
					currRow[n] = prevRow[n-1]
				} else {
					dist := min(prevRow[n-1], prevRow[n], currRow[n-1])
					currRow[n] = 1 + dist
				}
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
