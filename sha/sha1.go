package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

var filterPath *string // 目录过滤规则
var filterFile *string // 文件滤规则
var fileHash []string  // 最终结果数组

var wg sync.WaitGroup // 定义一个同步等待的组

func main() {
	// 全部核心运行程序
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 接收参数
	root := flags()
	// log.Println(root, filter)
	wg.Add(1)
	go run(root)
	wg.Wait()              // 阻塞等待所有组内成员都执行完毕退栈
	sort.Strings(fileHash) // 排序
	// log.Println(fileHash)
	// 写的文件名使用时间形式
	fileName := "sha1-each-file-in-html-dir"
	// 写入到文件
	err := writeToFile(fileName, fileHash)
	if err != nil {
		log.Println(fmt.Sprintf("写文件错误:%s", err.Error()))
	}
}

// 接收参数
func flags() string {
	root := flag.String("root", "./", "root directory")               // 要生成哈希的根目录
	filterPath = flag.String("filterpath", "", "directory to filter") // 要排除的目录
	filterFile = flag.String("filterfile", "", "file to filter")      // 要排除的文件
	help := flag.Bool("help", false, "Use the help")                  // 帮助
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}
	if *filterPath == "*" || *filterFile == "*" {
		log.Println("Filter rules cannot be '*'")
		os.Exit(1)
	}
	return *root
}

// 运行
func run(root string) {
	defer wg.Done()
	list := listFiles(root)
	for _, name := range list {
		// 拼接全路径
		fpath := filepath.Join(root, name)
		// 构造文件结构
		fio, _ := os.Lstat(fpath)
		// 如果遍历的当前文件是个目录，则进入该目录进行递归
		if fio.IsDir() {
			// 验证目录过滤规则
			if verifyFilter(fpath, true) == true {
				continue
			}
			wg.Add(1) // 为同步等待组增加一个成员
			go run(fpath)
			// run(fpath)
		} else {
			// 验证文件过滤规则
			if verifyFilter(fpath, false) == true {
				continue
			}
			// 获取文件信息
			info, err := os.Stat(fpath)
			if err != nil {
				log.Println(err)
			} else {
				oneFile := fmt.Sprintf("%s,%s,%d", fpath, sha1ToString(fpath), info.Size())
				fileHash = append(fileHash, oneFile)
			}
		}
	}
}

// 列出当前目录下的所有目录、文件
func listFiles(dirname string) []string {
	f, _ := os.Open(dirname)
	names, _ := f.Readdirnames(-1)
	f.Close()
	sort.Strings(names)
	return names
}

// 对字符串进行SHA1哈希
func sha1ToString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 验证过滤规则-目录或文件过滤规则
func verifyFilter(name string, isDir bool) bool {
	var filter *string
	// 分析逻辑一样，这里判断后把过滤规则赋值给新的变量
	if isDir == true {
		filter = filterPath
	} else {
		filter = filterFile
	}
	if *filter == "" {
		return false
	}
	// log.Println(*filter)
	if len(name) < len(strings.Trim(*filter, "*")) {
		return false
	}
	// 根据*号的位置，使用不同的匹配规则
	i := strings.Index(*filter, "*") // 查找*位置
	if i == 0 {
		log.Println((*filter)[1:], name[(len(name)-(len(*filter)-1)):])
		if (*filter)[1:] == name[(len(name)-(len(*filter)-1)):] {
			return true
		}
	} else if i == (len(*filter) - 1) {
		if (*filter)[:len(*filter)-1] == name[:len(*filter)-1] {
			return true
		}
	} else {
		if name == *filter {
			return true
		}
	}
	return false
}

// 写内容到文件
func writeToFile(fileName string, data []string) error {
	var str string
	for _, v := range data {
		str += fmt.Sprintf("%s\n", v)
	}
	err := ioutil.WriteFile(fmt.Sprintf("%s.csv", fileName), []byte(str), 0666)
	return err
}
