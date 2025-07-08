package levenshtein

func levenshteinDistance(s1, s2 string) int {
	m := len(s1)
	n := len(s2)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i == 0 {
				dp[i][j] = j
			} else if j == 0 {
				dp[i][j] = i
			} else {
				if s1[i-1] == s2[j-1] {
					dp[i][j] = dp[i-1][j-1]
				} else {
					dp[i][j] = 1 + min(dp[i][j-1], dp[i-1][j], dp[i-1][j-1])
				}
			}
		}
	}

	return dp[m][n]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func similarityPercentage(s1, s2 string) float64 {
	maxLen := max(len(s1), len(s2))
	if maxLen == 0 {
		return 100.0
	}
	distance := levenshteinDistance(s1, s2)
	return 100.0 - (float64(distance)/float64(maxLen))*100.0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Fuzzy(substr string, str string) float64 {
	return similarityPercentage(substr, str)
}
