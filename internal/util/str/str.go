package str

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"regexp"
	"strings"
	"sync"
)

func RegeJsonStr(str string) (res string) {
	strs := extractJSONFromText(str)
	res = strings.Join(strs, "")
	return
}

// extractJSONFromText 从文本中提取所有JSON对象
func extractJSONFromText(text string) []string {
	// 定义一个正则表达式来匹配可能的JSON起点
	re := regexp.MustCompile(`\{`)

	// 用于保存提取的JSON部分
	var jsonParts []string

	for {
		// 查找匹配的起点
		match := re.FindStringIndex(text)
		if match == nil {
			break
		}

		start := match[0]
		end := findJSONEnd(text[start:])
		if end != -1 {
			jsonParts = append(jsonParts, text[start:start+end+1])
			// 更新文本，跳过已匹配的部分
			text = text[start+end+1:]
		} else {
			// 如果没有找到结束位置，跳出循环
			break
		}
	}

	return jsonParts
}

// findJSONEnd 找到JSON对象的结束位置
func findJSONEnd(text string) int {
	braceCount := 0
	inQuotes := false

	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '{':
			if !inQuotes {
				braceCount++
			}
		case '}':
			if !inQuotes {
				braceCount--
				if braceCount == 0 {
					return i
				}
			}
		case '"':
			if i == 0 || text[i-1] != '\\' {
				inQuotes = !inQuotes
			}
		}
	}

	return -1
}

// popNElements 从切片中弹出前 n 个元素
func PopNElements(slice *[]string, n int) []string {
	if n > len(*slice) {
		n = len(*slice)
	}
	// 获取前 n 个元素
	poppedElements := (*slice)[:n]
	// 更新原切片，去掉前 n 个元素
	*slice = append((*slice)[n:], poppedElements...)
	return poppedElements
}

var popSync sync.Mutex

func PopElement(slice *[]string) string {
	popSync.Lock()
	defer popSync.Unlock()
	if len(*slice) == 0 {
		return ""
	}
	if len(*slice) == 1 {
		return (*slice)[0]
	}
	var res string
	if len(*slice) > 1 {
		res = (*slice)[0]
		*slice = append((*slice)[1:], (*slice)[0])
	}
	return res
}

// 判断字符串是否在切片中
func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func FilterSpecialChars(input string) string {
	// 使用正则表达式匹配非字母、数字和空格的字符，并将其替换为空字符串
	reg := regexp.MustCompile(`[^a-zA-Z0-9 \p{L}\p{N} \p{Han}\p{P}]`)
	filtered := reg.ReplaceAllString(input, " ")
	return filtered
}

func AddQuery(link string, params g.Map) string {
	paramsStr := ""
	for k, v := range params {
		paramsStr += fmt.Sprintf("&%s=%v", k, v)
	}
	if len(paramsStr) > 0 {
		paramsStr = "?" + paramsStr
	}
	return link + paramsStr
}

func SplitSlice(data []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}
