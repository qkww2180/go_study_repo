package algorithm

// 最长公共子串。原生解法
func LongestCommonSubstring(s1, s2 string) (string, int) {
	var maxLen int = 0
	var endIndex int //公共子串的末尾位置(不包含)
	for i := 0; i < len(s1); i++ {
		for j := 0; j < len(s2); j++ {
			//获取s1[i:]和s2[j:]的最长公共前缀
			prefixLen := 0 //目前得到的公共前缀的长度
			for k := 0; k < len(s2)-j && k < len(s1)-i; k++ {
				if s1[i+k] == s2[j+k] {
					prefixLen++
					if prefixLen > maxLen {
						maxLen = prefixLen
						endIndex = i + k + 1
					}
				} else {
					break
				}
			}
		}
	}
	return s1[endIndex-maxLen : endIndex], maxLen
}

// 最长公共子串。自下而上动态规划，空间开销M*N
func LongestCommonSubstringDP(s1, s2 string) (string, int) {
	dp := make([][]int, len(s1)+1)
	for i := 0; i < len(s1)+1; i++ {
		dp[i] = make([]int, len(s2)+1)
	}
	var maxLen int = 0
	var endIndex int //公共子串的末尾位置(不包含)
	for i := 0; i < len(s1); i++ {
		for j := 0; j < len(s2); j++ {
			if s1[i] == s2[j] {
				dp[i+1][j+1] = dp[i][j] + 1
				if dp[i+1][j+1] > maxLen {
					maxLen = dp[i+1][j+1]
					endIndex = i + 1
				}
			}
		}
	}
	return s1[endIndex-maxLen : endIndex], maxLen
}

// 最长公共子串。自下而上动态规划，空间开销2N
func LongestCommonSubstringDP_WithSpaceON(s1, s2 string) (string, int) {
	//使s2成为更短的那个，节省空间开销
	if len(s1) < len(s2) {
		s1, s2 = s2, s1
	}

	//只维护最下面的那两行
	prevRow := make([]int, len(s2)+1)
	currRow := make([]int, len(s2)+1)

	var maxLen int = 0
	var endIndex int //公共子串的末尾位置(不包含)
	for i := 0; i < len(s1); i++ {
		//填充currRow
		for j := 0; j < len(s2); j++ {
			if s1[i] == s2[j] {
				currRow[j+1] = prevRow[j] + 1
				if currRow[j+1] > maxLen {
					maxLen = currRow[j+1]
					endIndex = i + 1

				}
			} else {
				currRow[j+1] = 0 //对于每一个j都要给currRow[j+1]赋值，否则它会复用上一行的值
			}
		}
		//把currRow拷贝给prevRow
		for n, ele := range currRow {
			prevRow[n] = ele
		}
	}
	return s1[endIndex-maxLen : endIndex], maxLen
}

// 最长公共子串。自下而上动态规划，空间开销2N，使用go标准库的copy
func LongestCommonSubstringDP_WithSpaceON_StdCopy(s1, s2 string) (string, int) {
	//使s2成为更短的那个，节省空间开销
	if len(s1) < len(s2) {
		s1, s2 = s2, s1
	}

	//只维护最下面的那两行
	prevRow := make([]int, len(s2)+1)
	currRow := make([]int, len(s2)+1)

	var maxLen int = 0
	var endIndex int //公共子串的末尾位置(不包含)
	for i := 0; i < len(s1); i++ {
		//填充currRow
		for j := 0; j < len(s2); j++ {
			if s1[i] == s2[j] {
				currRow[j+1] = prevRow[j] + 1
				if currRow[j+1] > maxLen {
					maxLen = currRow[j+1]
					endIndex = i + 1

				}
			} else {
				currRow[j+1] = 0 //对于每一个j都要给currRow[j+1]赋值，否则它会复用上一行的值
			}
		}
		//把currRow拷贝给prevRow
		copy(prevRow, currRow) //go标准库的copy作了性能优化
	}
	return s1[endIndex-maxLen : endIndex], maxLen
}
