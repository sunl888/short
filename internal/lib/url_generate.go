package lib

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

const (
	VAL   = 0x3FFFFFFF
	INDEX = 0x0000003D
)

var (
	alphabet = []byte("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// 短链接生成：
// 将原长链接进行md5校验和计算，生成32位字符串
// 将32位字符串每8位划分一段，得到4段子串。将每个字串（16进制形式）转化为整型数值，与0x3FFFFFFF按位与运算，生成一个30位的数值
// 将上述生成的30位数值按5位为单位依次提取，得到的数值与0x0000003D按位与，获取一个0-61的整型数值，作为从字符数组中提取字符的索引。得到6个字符就生成了一个短链
// 4段字串共可以生成4个短链
func GenerateShortUrl(hash string) [4]string {
	var result [4]string
	for i := 0; i < 4; i++ {
		tmpUrl := hash[i*8 : (i+1)*8]
		calcTmpUrl, _ := strconv.ParseInt(tmpUrl, 16, 64)
		tmpVal := int64(VAL) & calcTmpUrl
		var index int64
		var uri []byte
		for j := 0; j < 6; j++ {
			index = INDEX & tmpVal
			uri = append(uri, alphabet[index])
			tmpVal >>= 5
		}
		result[i] = string(uri)
	}
	return result
}

// 生成指定字符串的 md5
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	c := m.Sum(nil)
	return hex.EncodeToString(c)
}
