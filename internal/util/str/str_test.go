package str

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"regexp"
	"testing"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func TestRegeJsonStr(t *testing.T) {
	// 正常的 JSON 字符串
	str1 := `12312{"key":"value","testArr":[{"title":"a","value":"val"},{"title":"b","value":"val"}],"testObj":{"title":"a","value":"val"}}啊手动阀`
	expected1 := `{"key":"value","testArr":[{"title":"a","value":"val"},{"title":"b","value":"val"}],"testObj":{"title":"a","value":"val"}}`

	// 不包含 JSON 的字符串
	str2 := "This is not JSON"
	expected2 := ""

	// 边界情况，空字符串
	str3 := ""
	expected3 := ""

	// 包含 简单JSON 的字符串
	str4 := `This is {"title":"a","value":"val"} JSON`
	expected4 := `{"title":"a","value":"val"}`

	// 测试函数
	result1 := RegeJsonStr(str1)
	result2 := RegeJsonStr(str2)
	result3 := RegeJsonStr(str3)
	result4 := RegeJsonStr(str4)

	// 验证结果
	if result1 != expected1 {
		t.Errorf("Expected %s, got %s", expected1, result1)
	}

	if result2 != expected2 {
		t.Errorf("Expected %s, got %s", expected2, result2)
	}

	if result3 != expected3 {
		t.Errorf("Expected %s, got %s", expected3, result3)
	}

	if result4 != expected4 {
		t.Errorf("Expected %s, got %s", expected4, result4)
	}
}

var bg_musics = []string{
	"90.mp3",
	"91.mp3",
	"92.mp3",
	"93.mp3",
	"94.mp3",
	"95.mp3",
	"96.mp3",
	"97.mp3",
	"98.mp3",
	"99.mp3",
	"100.mp3",
	"101.mp3",
	"102.mp3",
	"103.mp3",
	"104.mp3",
	"105.mp3",
	"106.mp3",
}

func TestPop(t *testing.T) {
	fmt.Println(len(bg_musics))
	for i := 0; i < 5; i++ {
		filename := PopElement(&bg_musics)
		fmt.Println(filename)
	}
	fmt.Println(len(bg_musics))
	fmt.Println((bg_musics))
	// discount := "-0%"

	// // 定义正则表达式模式，匹配discount为0的情况
	// pattern := `^(-0)%$`

	// // 编译正则表达式
	// re := regexp.MustCompile(pattern)

	// // 进行匹配
	// if re.MatchString(discount) {
	// 	fmt.Println("Discount为0")
	// } else {
	// 	fmt.Println("Discount不为0")
	// }
}

func TestFilterSpecialChars(t *testing.T) {
	// str := "-10%"
	// re := regexp.MustCompile(`(\d*\.?\d+)`)
	// match := re.FindStringSubmatch(str)

	// if len(match) > 1 {
	// 	fmt.Println(match[1])
	// }
	str := "Hello, 你好，123!@#$%^&*()_+ ♠♣♧♡♥❤❥❣♂♀✲☀☼☾☽◐◑☺☻☎☏✿❀№↑↓←→√×÷★℃℉°◆◇⊙■□△▽¿½☯✡㍿卍卐♂♀✚〓㎡♪♫♩♬㊚㊛囍㊒㊖Φ♀♂‖$@*&#※卍卐Ψ♫♬♭♩♪♯♮⌒¶∮‖€￡¥$§"
	fmt.Println(FilterSpecialChars(str))
	discount := "-0%"
	// 创建一个正则表达式来匹配数字部分
	re := regexp.MustCompile(`(\d*\.?\d+)`)
	match := re.FindStringSubmatch(discount)

	if len(match) > 1 {
		discount = match[1]
	}
	fmt.Println(discount)
}

func TestSplitSlice(t *testing.T) {
	slice := SplitSlice(bg_musics, 10)
	g.Dump(slice)
}
