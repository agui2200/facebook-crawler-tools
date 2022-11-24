package global

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 判断string类型变量是否为空
func IsEmpty(str ...string) bool {
	for _, str := range str {
		if len(str) <= 0 {
			return true
		}
	}

	return false
}

// 生成指定数量的随机字符串
func RandomStr(len int) string {

	buff := make([]byte, int(math.Ceil(float64(len)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	str = strings.ToLower(str[:len])
	str = strings.ReplaceAll(str, "_", "1")
	str = strings.ReplaceAll(str, "-", "2")
	str = strings.ReplaceAll(str, "=", "3")
	str = strings.ReplaceAll(str, "/", "4")
	return str // strip 1 extra character we get from odd length results
}

// 取出中间文本
func StrBetween(str, starting, ending string) string {
	s := strings.Index(str, starting)
	if s < 0 {
		return ""
	}
	s += len(starting)
	e := strings.Index(str[s:], ending)
	if e < 0 {
		return ""
	}
	return str[s : s+e]
}

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Base64Decode(input string) string {
	str, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	return string(str)
}

// string数组切片去重复+去空白
func RemoveDuplicateElement(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok { //如果字典中找不到元素，ok=false，!ok为true，就往切片中append元素。
			temp[item] = struct{}{}

			//判断数据是否为空
			item = strings.TrimSpace(item)
			if len(item) <= 0 || item == "" {
				continue
			}
			result = append(result, item)
		}
	}
	return result
}

// 获取时间
func GetTime(format string) string {
	time := time.Unix(time.Now().Unix(), 0).Format(format)
	return time
}

// 获取UUID
func GetUUID() (bool, string) {
	guid := uuid.New().String()
	if len(guid) <= 0 {
		log.Println("get_uuid error!")
		return false, ""
	}

	return true, guid
}

// 生成token
func GeneratetToken() (bool, string) {
	_, token := GetUUID()
	w := md5.New()
	io.WriteString(w, token)

	//将token写入到w中
	token = fmt.Sprintf("%X", w.Sum(nil))
	if len(token) <= 0 {
		return false, ""
	}

	return true, token
}
