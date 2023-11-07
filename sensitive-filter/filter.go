package filter

import (
	"strings"
)

type filter struct {
	children map[string]*filter
	isEnd    bool
}

const InvalidWords = " ,~,!,@,#,$,%,^,&,*,(,),_,-,+,=,?,<,>,.,—,，,。,/,\\,|,《,》,？,;,:,：,',‘,；,“,"

var InvalidSet = make(map[string]struct{})

var sensitiveFilter = &filter{
	children: make(map[string]*filter),
	isEnd:    false,
}
var sensitiveWords = []string{
	"草泥马",
	"操你妈",
	"去死",
	"该死",
	"死",
	"狗东西",
	"你妈逼的",
	"你妈",
	"日",
	"逼",
}

func AddSensitiveWord(words ...string) {
	for _, word := range words {
		wordsFilter := sensitiveFilter
		arr := []rune(word)
		for i := range arr {
			if _, ok := wordsFilter.children[string(arr[i])]; !ok {
				tmpFilter := &filter{
					children: make(map[string]*filter),
					isEnd:    i == len(arr)-1,
				}
				wordsFilter.children[string(arr[i])] = tmpFilter
			} else {
				tmpFilter := wordsFilter.children[string(arr[i])]
				tmpFilter.isEnd = i == len(arr)-1
			}
			wordsFilter = wordsFilter.children[string(arr[i])]
		}
	}
}

func init() {
	words := strings.Split(InvalidWords, ",")
	for _, word := range words {
		InvalidSet[word] = struct{}{}
	}
	AddSensitiveWord(sensitiveWords...)
}

func ChangeSensitiveWords(text string) string {
	arr := []rune(text)
	nowMap := sensitiveFilter
	start := -1
	tag := -1
	for i := 0; i < len(arr); i++ {
		if _, ok := InvalidSet[(string(arr[i]))]; ok || string(arr[i]) == "," {
			continue
		}
		if thisMap, ok := nowMap.children[string(arr[i])]; ok {
			tag++
			if tag == 0 {
				start = i

			}
			isEnd := thisMap.isEnd
			if isEnd {
				for y := start; y < i+1; y++ {
					arr[y] = 42
				}
				nowMap = sensitiveFilter
				start = -1
				tag = -1
			} else {
				nowMap = nowMap.children[string(arr[i])]
			}

		} else {
			if start != -1 {
				i = start
			}
			nowMap = sensitiveFilter
			start = -1
			tag = -1
		}
	}

	return string(arr)
}
