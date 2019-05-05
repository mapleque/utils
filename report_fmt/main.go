package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var config = map[string]map[string]string{
	"_default": {
		"_line": "1. %s",
	},
	"status": {
		"_line": "<span style='color:$vmap'>$v</span>%s",
		"进行":    "blue",
		"完成":    "gray",
		"取消":    "red",
		"计划":    "green",
	},
	"category": {
		"_line": "<span style='color:$vmap'>[$v]%s</span>",
		"保险":    "#DB9019",
		"理财":    "#5ED5D1",
		"基础":    "#1A2D27",
		"数据":    "#FF6E97",
		"运维":    "#F1AAA6",
	},
	"content": {
		"_line": "$v",
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "need file args")
		return
	}

	file := os.Args[1]
	fileFmt(file)
}

func fileFmt(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	cont, err := ioutil.ReadAll(f)
	for _, line := range strings.Split(string(cont), "\n") {
		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}
		ret := lineDessect(line, "[%{status}][%{category}]%{content}")

		tar := config["_default"]["_line"]
		tar = lineFmt(tar, "status", ret["status"])
		tar = lineFmt(tar, "category", ret["category"])
		tar = lineFmt(tar, "content", ret["content"])
		fmt.Println(tar)
	}
}

func lineDessect(str, expr string) map[string]string {
	ret := map[string]string{}
	start, end, lastSep := 0, len(str), ""
	exprArr := strings.Split(expr, "%{")
	if len(exprArr) < 1 {
		panic(fmt.Errorf("invalid dissect expr %s", expr))
	}
	for _, exprItem := range exprArr {
		if !strings.Contains(exprItem, "}") {
			// 先把前面可能存在的不提取的字符去掉
			// 例如：abcd%{key}...
			//  切割出来的exprItem是abcd
			//  这样就在这里去掉了
			end = strings.Index(str, exprItem)
			start = end + len(exprItem)
			lastSep = exprItem
			continue
		}
		// 这里的exprItem是key}sep的形式
		// 只需要按照这个sep去截取字符串
		// 所得到的字符串就是当前key的值
		arr := strings.Split(exprItem, "}")
		if len(arr) != 2 {
			panic(fmt.Errorf("invalid dissect expr arr `%s` `%s`", exprItem, arr))
		}
		key, sep := arr[0], arr[1]
		if len(key) < 1 {
			panic(fmt.Errorf("invalid dissect expr key `%s` `%s`", exprItem, arr))
		}
		if len(sep) > 0 {
			end = strings.Index(str[start:], sep)
			if end < 0 {
				panic(fmt.Errorf("%s not match on %s with seperator %s", str[start:], key, sep))
			}
			end = start + end
		}
		// 这里如果key是+开头的，就说明是把它追加到已有的key上面
		if key[0] == '+' {
			key = key[1:]
			ret[key] += str[start-len(lastSep) : end]
		} else {
			ret[key] = str[start:end]
		}
		lastSep = sep
		start = end + len(sep)
		end = len(str)
	}
	return ret
}

func lineFmt(layout, k, v string) string {
	value := config[k]["_line"]
	value = strings.Replace(value, "$vmap", config[k][v], -1)
	value = strings.Replace(value, "$v", v, -1)
	return fmt.Sprintf(layout, value)
}
