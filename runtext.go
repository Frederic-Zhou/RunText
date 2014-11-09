package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/djimenez/iconv-go"
	"os"
	"os/exec"
	"strings"
)

func main() {

	path := flag.String("p", "", "*(必要) 文件路径")
	count := flag.Int("c", 0, "读取行数,默认为0，代表读取所有行")
	substr := flag.String("s", "", "搜索内容，默认搜索全部")
	linenum := flag.Bool("l", false, "是否显示行数，默认不显示")
	ignore := flag.Int("i", 0, "忽略的行数，默认不忽略")
	fromEncoding := flag.String("ef", "utf-8", "文档的编码，默认'utf-8'")
	toEncoding := flag.String("et", "utf-8", "欲转换为的编码，默认'utf-8'")
	cmd := flag.String("cmd", "", "外部命令，每一行读取完毕后执行的命令，会将行内容传递给命令作为参数，默认无命令")
	isAll := false
	/*扩展*/

	flag.Parse()

	if *count == 0 {
		isAll = true
	}

	selectTxt(*path, isAll, *ignore, *count, *substr, *linenum, *fromEncoding, *toEncoding, *cmd)
	fmt.Println("Read Over!!")

}

func selectTxt(path string, isAll bool, ignore int, count int, substr string, linenum bool, fromEncoding string, toEncoding string, cmd string) {

	// 打开文件流
	fp, err := os.Open(path)
	if err != nil {
		fmt.Println("[open]:", err.Error())
		fmt.Println("type '--help' to show more infomation:")
		return
	}
	defer fp.Close()

	// 创建文件阅读器
	r := bufio.NewReader(fp)

	// 忽略掉指定行数
	for i := 0; i < ignore; i++ {
		if _, _, err := r.ReadLine(); err != nil {
			if err.Error() != "EOF" {
				fmt.Println("[ignore line]:", err.Error())
			}
			break
		}
	}

	//循环指定行号，或者全部行
	for i := 0; i < count || isAll; i++ {
		line, prefix, err := r.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("[read line]:", err.Error())
			}
			break
		}

		_ = prefix //暂时丢弃prefix

		var linestring string
		if fromEncoding != toEncoding {
			linestring, _ = iconv.ConvertString(string(line), fromEncoding, toEncoding) // 转换编码并转换程字符串
		} else {
			linestring = string(line)
		}

		//判断行中是否有指定包含的内容
		if strings.Contains(linestring, substr) {
			// 根据指定参数判断是否显示行号
			if linenum {
				fmt.Println(i+1+ignore, linestring)
			} else {
				fmt.Println(string(line))
			}

			if cmd != "" {
				command := exec.Command(cmd, linestring)
				buf, err := command.Output()
				if err != nil {
					fmt.Println("[command err]:", err.Error())
				} else {
					fmt.Println(cmd, ":", string(buf))
				}

			}

		}
	}

}
