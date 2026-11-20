// Command clienthash 计算客户端文件指纹，并直接输出可粘贴进允许清单的 JSON 条目。
//
// 用法（在放客户端文件的机器上运行）：
//
//	go run ./cmd/clienthash -exe DNF.exe -pvf Script.pvf -launcher 登录器.exe -label v211-20261120
//
// 输出示例：
//
//	{"label":"v211-20261120","exe":"…","pvf":"…","launcher":"…"}
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func main() {
	exe := flag.String("exe", "", "DNF.exe 路径")
	pvf := flag.String("pvf", "", "Script.pvf 路径")
	launcher := flag.String("launcher", "", "登录器可执行文件路径")
	label := flag.String("label", "", "版本标签，写入允许清单")
	flag.Parse()

	if *exe == "" || *pvf == "" || *launcher == "" || *label == "" {
		fmt.Fprintln(os.Stderr, "必须同时提供 -exe -pvf -launcher -label")
		flag.Usage()
		os.Exit(2)
	}

	entry := map[string]string{"label": *label}
	for name, path := range map[string]string{"exe": *exe, "pvf": *pvf, "launcher": *launcher} {
		sum, err := hashFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "计算 %s 失败: %v\n", name, err)
			os.Exit(1)
		}
		entry[name] = sum
		fmt.Fprintf(os.Stderr, "%s: %s\n", name, sum)
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成 JSON 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(encoded))
}
