package helper

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// IsLetters 字符串是否全(英文)字母组成.
func (ts *TsStr) IsLetters(str string) bool {
	for _, r := range str {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return str != ""
}

// IsEmpty 字符串是否为空(包括空格).
func (ts *TsStr) IsEmpty(str string) bool {
	if len(str) == 0 || len(ts.Trim(str)) == 0 {
		return true
	}

	return false
}

// IsUpper 字符串是否全部大写.
func (ts *TsStr) IsUpper(str string) bool {
	for _, r := range str {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return str != ""
}

// IsLower 字符串是否全部小写.
func (ts *TsStr) IsLower(str string) bool {
	for _, r := range str {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return str != ""
}

// HasLetter 字符串是否含有(英文)字母.
func (ts *TsStr) HasLetter(str string) bool {
	for _, r := range str {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}

// IsUtf8 字符串是否UTF-8编码.
func (ts *TsStr) IsUtf8(str string) bool {
	return str != "" && utf8.ValidString(str)
}

// IsASCII 是否IsASCII字符串.
func (ts *TsStr) IsASCII(str string) bool {
	//return str != "" && RegAscii.MatchString(str)
	n := len(str)
	for i := 0; i < n; i++ {
		if str[i] > 127 {
			return false
		}
	}

	return str != ""
}

// IsMultibyte 字符串是否含有多字节字符.
func (ts *TsStr) IsMultibyte(str string) bool {
	return str != "" && RegMultiByte.MatchString(str)
}

// HasFullWidth 是否含有全角字符.
func (ts *TsStr) HasFullWidth(str string) bool {
	return str != "" && RegFullWidth.MatchString(str)
}

// HasHalfWidth 是否含有半角字符.
func (ts *TsStr) HasHalfWidth(str string) bool {
	return str != "" && RegHalfWidth.MatchString(str)
}

// IsEnglish 字符串是否纯英文.letterCase是否检查大小写,枚举值(CaseNone,CASE_LOWER,CASE_UPPER).
func (ts *TsStr) IsEnglish(str string, letterCase LetterCase) bool {
	switch letterCase {
	case CaseNone:
		return ts.IsLetters(str)
	case CaseLower:
		return str != "" && RegAlphaLower.MatchString(str)
	case CaseUpper:
		return str != "" && RegAlphaUpper.MatchString(str)
	default:
		return ts.IsLetters(str)
	}
}

// HasEnglish 是否含有英文字符,HasLetter的别名.
func (ts *TsStr) HasEnglish(str string) bool {
	return ts.HasLetter(str)
}

// HasChinese 字符串是否含有中文.
func (ts *TsStr) HasChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}

	return false
}

// IsChinese 字符串是否全部中文.
func (ts *TsStr) IsChinese(str string) bool {
	return str != "" && RegChineseAll.MatchString(str)
}

// IsChineseName 字符串是否中文名称.
func (ts *TsStr) IsChineseName(str string) bool {
	return str != "" && RegChineseName.MatchString(str)
}

// HasSpecialChar 字符串是否含有特殊字符.
func (ts *TsStr) HasSpecialChar(str string) (res bool) {
	if str == "" {
		return
	}

	for _, r := range str {
		// IsPunct 判断 r 是否为一个标点字符 (类别 P)
		// IsSymbol 判断 r 是否为一个符号字符
		// IsMark 判断 r 是否为一个 mark 字符 (类别 M)
		if unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsMark(r) {
			res = true
			return
		}
	}

	return
}

// IsJSON 字符串是否合法的json格式.
func (ts *TsStr) IsJSON(str string) bool {
	length := len(str)
	if length == 0 {
		return false
	} else if (str[0] != '{' || str[length-1] != '}') && (str[0] != '[' || str[length-1] != ']') {
		return false
	}

	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// IsIP 检查字符串是否IP地址.
func (ts *TsStr) IsIP(str string) bool {
	return str != "" && net.ParseIP(str) != nil
}

// IsIPv4 检查字符串是否IPv4地址.
func (ts *TsStr) IsIPv4(str string) bool {
	ipAddr := net.ParseIP(str)
	// 不是合法的IP地址
	if ipAddr == nil {
		return false
	}

	return ipAddr.To4() != nil && strings.ContainsRune(str, '.')
}

// IsIPv6 检查字符串是否IPv6地址.
func (ts *TsStr) IsIPv6(str string) bool {
	ipAddr := net.ParseIP(str)
	return ipAddr != nil && strings.ContainsRune(str, ':')
}

// IsPort 字符串或数字是否端口号.
func (ts *TsStr) IsPort(val interface{}) bool {
	if TConv.IsInt(val) {
		port := TConv.ToInt(val)
		if port > 0 && port < 65536 {
			return true
		}
	}

	return false
}

// IsDNSName 是否DNS名称.
func (ts *TsStr) IsDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		// constraints already violated
		return false
	}
	return !ts.IsIP(str) && RegDNSName.MatchString(str)
}

// IsDialString 是否网络拨号字符串(形如127.0.0.1:80),用于net.Dial()检查.
func (ts *TsStr) IsDialString(str string) bool {
	h, p, err := net.SplitHostPort(str)
	if err == nil && h != "" && p != "" && (ts.IsDNSName(h) || ts.IsIP(h)) && ts.IsPort(p) {
		return true
	}

	return false
}

// IsMACAddr 是否MAC物理网卡地址.
func (ts *TsStr) IsMACAddr(str string) bool {
	_, err := net.ParseMAC(str)
	return err == nil
}

// IsHost 字符串是否主机名(IP或DNS名称).
func (ts *TsStr) IsHost(str string) bool {
	return ts.IsIP(str) || ts.IsDNSName(str)
}

// IsEmail 检查字符串是否邮箱.参数validateTrue,是否验证邮箱主机的真实性.
func (ts *TsStr) IsEmail(email string, validateHost bool) (bool, error) {
	//长度检查
	length := len(email)
	at := strings.LastIndexByte(email, '@')
	if (length < 6 || length > 254) || (at <= 0 || at > length-3) {
		return false, fmt.Errorf("invalid email length")
	}

	// 验证邮箱格式
	chkFormat := RegEmail.MatchString(email)
	if !chkFormat {
		return false, fmt.Errorf("invalid email format")
	}

	// 验证主机
	if validateHost {
		host := email[at+1:]
		if _, err := net.LookupMX(host); err != nil {
			// 因无法确定mx主机的smtp端口,所以去掉Hello/Mail/Rcpt检查邮箱是否存在
			// 仅检查主机是否有效
			// 仅对国内几家大的邮件厂家进行检查
			if _, err = net.LookupIP(host); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

// IsMobileCN 检查字符串是否中国大陆手机号.
func (ts *TsStr) IsMobileCN(str string) bool {
	return str != "" && RegMobileCN.MatchString(str)
}

// IsTel 是否固定电话或400/800电话.
func (ts *TsStr) IsTel(str string) bool {
	return str != "" && RegTelephone.MatchString(str)
}

// IsPhone 是否电话号码(手机或固话).
func (ts *TsStr) IsPhone(str string) bool {
	return str != "" && RegPhone.MatchString(str)
}

// IsCreditNo 检查是否(15或18位)身份证号码,并返回经校验的号码.
func (ts *TsStr) IsCreditNo(str string) (bool, string) {
	chk := str != "" && RegCreditNo.MatchString(str)
	if !chk {
		return false, ""
	}

	// 检查省份代码
	if _, chk = CreditArea[str[0:2]]; !chk {
		return false, ""
	}

	// 将15位身份证升级到18位
	length := len(str)
	if length == 15 {
		// 先转为17位,如果身份证顺序码是996 997 998 999,这些是为百岁以上老人的特殊编码
		if chk, _ = ts.DStrPos(str[12:], []string{"996", "997", "998", "999"}, false); chk {
			str = str[0:6] + "18" + str[6:]
		} else {
			str = str[0:6] + "19" + str[6:]
		}

		// 再加上校验码
		code := append([]byte{}, creditChecksum(str))
		str += string(code)
	}

	// 检查生日
	birthday := str[6:10] + "-" + str[10:12] + "-" + str[12:14]
	chk, tim := TTime.IsDate2time(birthday)
	now := TTime.Time()
	if !chk {
		return false, ""
	} else if tim >= now {
		return false, ""
	}

	// 18位身份证需要验证最后一位校验位
	if length == 18 {
		str = strings.ToUpper(str)
		if str[17] != creditChecksum(str) {
			return false, ""
		}
	}

	return true, str
}

// IsAlphaNumeric 是否字母或数字.
func (ts *TsStr) IsAlphaNumeric(str string) bool {
	return str != "" && RegAlphaNumeric.MatchString(str)
}

// IsHEXColor 检查是否十六进制颜色,并返回带"#"的修正值.
func (ts *TsStr) IsHEXColor(str string) (bool, string) {
	chk := str != "" && RegRGBColor.MatchString(str)
	if chk && !strings.ContainsRune(str, '#') {
		str = "#" + strings.ToUpper(str)
	}
	return chk, str
}

// IsRGBColor 检查字符串是否RGB颜色格式.
func (ts *TsStr) IsRGBColor(str string) bool {
	return str != "" && RegRGBColor.MatchString(str)
}

// IsBlank 是否空(空白)字符.
func (ts *TsStr) IsBlank(str string) bool {
	// Check length
	if len(str) > 0 {
		// Iterate string
		for i := range str {
			// Check about char different from whitespace
			// 227为全角空格
			if str[i] > 32 && str[i] != 227 {
				return false
			}
		}
	}
	return true
}

// IsWhitespaces 是否全部空白字符,不包括空字符串.
func (ts *TsStr) IsWhitespaces(str string) bool {
	return str != "" && RegWhitespaceAll.MatchString(str)
}

// HasWhitespace 是否带有空白字符.
func (ts *TsStr) HasWhitespace(str string) bool {
	return str != "" && RegWhitespaceHas.MatchString(str)
}

// IsBase64 是否base64字符串.
func (ts *TsStr) IsBase64(str string) bool {
	return str != "" && RegBase64.MatchString(str)
}

// IsBase64Image 是否base64编码的图片.
func (ts *TsStr) IsBase64Image(str string) bool {
	if str == "" || !strings.ContainsRune(str, ',') {
		return false
	}

	dataURI := strings.Split(str, ",")
	return RegBase64Image.MatchString(dataURI[0]) && RegBase64.MatchString(dataURI[1])
}

// IsRsaPublicKey 检查字符串是否RSA的公钥,keyLength为密钥长度.
func (ts *TsStr) IsRsaPublicKey(str string, keyLength int) bool {
	bb := bytes.NewBufferString(str)
	pemBytes, _ := ioutil.ReadAll(bb)

	// 获取公钥
	block, _ := pem.Decode(pemBytes)
	if block != nil && block.Type != "PUBLIC KEY" {
		return false
	}
	var der []byte
	var err error

	if block != nil {
		der = block.Bytes
	} else {
		der, err = base64.StdEncoding.DecodeString(str)
		if err != nil {
			return false
		}
	}

	key, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return false
	}
	pubkey, ok := key.(*rsa.PublicKey)
	if !ok {
		return false
	}
	bitlen := len(pubkey.N.Bytes()) * 8
	return bitlen == keyLength
}

// IsUrl 检查字符串是否URL.
func (ts *TsStr) IsUrl(str string) bool {
	if str == "" || len(str) <= 3 || utf8.RuneCountInString(str) >= 2083 || strings.HasPrefix(str, ".") {
		return false
	}

	res, err := url.ParseRequestURI(str)
	if err != nil {
		return false //Couldn't even parse the url
	}
	if len(res.Scheme) == 0 {
		return false //No Scheme found
	}

	return true
}

// IsUrlExists 检查URL是否存在.
func (ts *TsStr) IsUrlExists(str string) bool {
	if !ts.IsUrl(str) {
		return false
	}

	client := &http.Client{}
	resp, err := client.Head(str)
	if err != nil {
		return false
	} else if resp.StatusCode == 404 {
		return false
	}

	return true
}

// IsMd5 是否md5值.
func (ts *TsStr) IsMd5(str string) bool {
	return str != "" && RegMd5.MatchString(str)
}

// IsSha1 是否Sha1值.
func (ts *TsStr) IsSha1(str string) bool {
	return str != "" && RegSha1.MatchString(str)
}

// IsSha256 是否Sha256值.
func (ts *TsStr) IsSha256(str string) bool {
	return str != "" && RegSha256.MatchString(str)
}

// IsSha512 是否Sha512值.
func (ts *TsStr) IsSha512(str string) bool {
	return str != "" && RegSha512.MatchString(str)
}

// StartsWith 字符串str是否以substr开头.
func (ts *TsStr) StartsWith(str, substr string) bool {
	if str != "" && substr != "" && ts.MbSubstr(str, 0, len([]rune(substr))) == substr {
		return true
	}
	return false
}

// EndsWith 字符串str是否以substr结尾.
func (ts *TsStr) EndsWith(str, substr string) bool {
	if str != "" && substr != "" && ts.MbSubstr(str, -len([]rune(substr))) == substr {
		return true
	}
	return false
}

// IsArrayOrSlice 检查变量是否数组或切片.
// chkType检查类型,枚举值有(1仅数组,2仅切片,3数组或切片);结果为-1表示非,>=0表示是
func (ta *TsArr) IsArrayOrSlice(val interface{}, chkType uint8) int {
	return isArrayOrSliceHelper(val, chkType)
}

// IsMap 检查变量是否字典.
func (ta *TsArr) IsMap(val interface{}) bool {
	return isMap(val)
}

// IsMapBySprintf 是否是map,通过fmt.Sprintf判断
func (ta *TsArr) IsMapBySprintf(i interface{}) bool {
	m := fmt.Sprintf("%T", i)
	return strings.HasPrefix(m, "map[")
}

// IsDate2time 检查字符串是否日期格式,并转换为时间戳.注意,时间戳可能为负数(小于1970年时).
// 匹配如:
//	0000
//	0000-00
//	0000/00
//	0000-00-00
//	0000/00/00
//	0000-00-00 00
//	0000/00/00 00
//	0000-00-00 00:00
//	0000/00/00 00:00
//	0000-00-00 00:00:00
//	0000/00/00 00:00:00
// 等日期格式.
func (tk *TsTime) IsDate2time(str string) (bool, int64) {
	if str == "" {
		return false, 0
	} else if strings.ContainsRune(str, '/') {
		str = strings.Replace(str, "/", "-", -1)
	}

	chk := RegDatetime.MatchString(str)
	if !chk {
		return false, 0
	}

	leng := len(str)
	if leng < 19 {
		reference := "1970-01-01 00:00:00"
		str = str + reference[leng:19]
	}

	tim, err := TTime.Str2Timestamp(str)
	if err != nil {
		return false, 0
	}

	return true, tim
}

// IsNan 是否为“非数值”.
func (ti *TsInt) IsNan(val float64) bool {
	return math.IsNaN(val)
}

// IsString 变量是否字符串.
func (tc *TsConvert) IsString(val interface{}) bool {
	return tc.Gettype(val) == "string"
}

// IsBinary 字符串是否二进制.
func (tc *TsConvert) IsBinary(s string) bool {
	for _, b := range s {
		if 0 == b {
			return true
		}
	}
	return false
}

// IsNumeric 变量是否数值(不包含复数和科学计数法).
func (tc *TsConvert) IsNumeric(val interface{}) bool {
	return isNumeric(val)
}

// IsInt 变量是否整型数值.
func (tc *TsConvert) IsInt(val interface{}) bool {
	return isInt(val)
}

// IsFloat 变量是否浮点数值.
func (tc *TsConvert) IsFloat(val interface{}) bool {
	return isFloat(val)
}

// IsEmpty 检查变量是否为空.
func (tc *TsConvert) IsEmpty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// IsNil 检查变量是否空值.
func (tc *TsConvert) IsNil(val interface{}) bool {
	if val == nil {
		return true
	}

	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Interface:
		if rv.IsNil() {
			return true
		}
	}
	return false
}

// IsBool 是否布尔值.
func (tc *TsConvert) IsBool(val interface{}) bool {
	return val == true || val == false
}

// IsHex 是否十六进制字符串.
func (tc *TsConvert) IsHex(str string) bool {
	_, err := tc.Hex2Dec(str)
	return err == nil
}

// IsByte 变量是否字节切片.
func (tc *TsConvert) IsByte(val interface{}) bool {
	return tc.Gettype(val) == "[]uint8"
}

// IsStruct 变量是否结构体.
func (tc *TsConvert) IsStruct(val interface{}) bool {
	r := reflectPtr(reflect.ValueOf(val))
	return r.Kind() == reflect.Struct
}

// IsInterface 变量是否接口.
func (tc *TsConvert) IsInterface(val interface{}) bool {
	r := reflectPtr(reflect.ValueOf(val))
	return r.Kind() == reflect.Invalid
}

// IsOdd 变量是否奇数.
func (ti *TsInt) IsOdd(val int) bool {
	return val%2 != 0
}

// IsEven 变量是否偶数.
func (ti *TsInt) IsEven(val int) bool {
	return val%2 == 0
}

// IsRangeInt 数值是否在2个整数范围内.
func (ti *TsInt) IsRangeInt(value, left, right int) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}


// IsNegative 数值是否为负数.
func (tf *TsFloat) IsNegative(value float64) bool {
	return value < 0
}

// IsPositive 数值是否为正数.
func (tf *TsFloat) IsPositive(value float64) bool {
	return value > 0
}

// IsNonNegative 数值是否为非负数.
func (tf *TsFloat) IsNonNegative(value float64) bool {
	return value >= 0
}

// IsNonPositive 数值是否为非正数.
func (tf *TsFloat) IsNonPositive(value float64) bool {
	return value <= 0
}

// IsWhole 数值是否为整数.
func (tf *TsFloat) IsWhole(value float64) bool {
	return math.Remainder(value, 1) == 0
}

// IsRangeFloat32 数值是否在2个32位浮点数范围内.
func (tf *TsFloat) IsRangeFloat32(value, left, right float32) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// IsRangeFloat64 数值是否在2个64位浮点数范围内.
func (tf *TsFloat) IsRangeFloat64(value, left, right float64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// AverageFloat64 对浮点数序列求平均值.
func (tf *TsFloat) AverageFloat64(nums ...float64) (res float64) {
	length := len(nums)
	if length == 0 {
		return
	} else if length == 1 {
		res = nums[0]
	} else {
		total := tf.SumFloat64(nums...)
		res = total / float64(length)
	}

	return
}

// IsDir 是否目录(且存在)
func (tf *TsFile) IsDir(filePath string) bool {
	f, err := os.Lstat(filePath)
	if os.IsNotExist(err) || nil != err {
		return false
	}
	return f.IsDir()
}

// IsFileExist 文件是否存在
func (tf *TsFile) IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}
