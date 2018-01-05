package cdkey

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	pkgBean "webapi/bean"
)

const letterBytes = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func ParseGiftString(plain string) map[int32]int32 {
	results := make(map[int32]int32)

	arr1 := strings.Split(plain, ",")
	if n1 := len(arr1); n1 > 0 {
		for i := 0; i < n1; i++ {
			v := arr1[i]
			arr2 := strings.Split(v, ":")
			if n2 := len(arr2); n2 == 2 {
				v1, _ := strconv.ParseInt(arr2[0], 10, 32)
				v2, _ := strconv.ParseInt(arr2[1], 10, 32)
				if v1 > 0 && v2 > 0 {
					results[int32(v1)] = int32(v2)
				}
			}
		}
	}

	return results
}

func generateCDKeyString(num, length int) []string {
	cdkeys := make(map[string]bool)
	for len(cdkeys) != num {
		v := RandStringBytesMaskImprSrc(length)
		cdkeys[v] = true
	}

	i := 0
	results := make([]string, num)
	for v, _ := range cdkeys {
		results[i] = v
		i++
	}

	return results
}

func GenerateWithGiftID(category uint32, channel string, num int, deadline int64, gift int) (error, int, []string) {
	results := []string{}

	for num > len(results) {
		cdkeys := generateCDKeyString(num-len(results), 9)
		err, batch := pkgBean.BatchAddCDKeyWithGift(category, channel, gift, deadline, cdkeys)
		if err != nil {
			return err, 0, nil
		}

		results = append(results, batch...)
	}

	return nil, gift, results
}

func Generate(category uint32, channel string, num int, deadline int64, resources map[int32]int32) (error, int, []string) {
	err, bean := pkgBean.AddCDKeyGift(resources)
	if err != nil {
		return err, 0, nil
	}

	return GenerateWithGiftID(category, channel, num, deadline, bean.Gift)
}
