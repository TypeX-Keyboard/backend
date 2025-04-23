package common

import (
	"time"

	"golang.org/x/exp/rand"
)

func RandStrSlice(slice []string) {
	// 初始化随机数种子
	seed := time.Now().UnixNano()
	rand.Seed(uint64(seed))

	// 使用 rand.Shuffle 函数打乱字符串数组的顺序
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

}

func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
