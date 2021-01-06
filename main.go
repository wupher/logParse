package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	fileList := []string{
		"/Users/fanwu/Documents/网盘/Nutstore/Documents/K米/Workspace/logs/19.log",
	}
	initData()
	t := time.Now()

	for _, v := range fileList {
		s := time.Now()
		file, isGzip, err := getFile(v)
		if err != nil {
			fmt.Printf("read file %v error, error %v\n", v, err)
			continue
		}
		err = process(file, isGzip)
		if err != nil {
			fmt.Printf("Handle file %v err %v \n", v, err)
		}
		_ = file.Close()
		fmt.Printf("File %v time take - %v \n", v, time.Since(s))
	}
	fmt.Printf("All Time taken - %v \n", time.Since(t))
	outPut()
}

func parseArgs() []string {
	args := os.Args[1:]
	if len(args) != 3 {
		log.Fatalln("Format ./coreader 200727 0 190(date start end)")
	}
	format := args[0]
	//Unicode 转码，貌似我不一定需要
	start, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalln("Format ./coreader 200727 0 190(date start end)")
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalln("Format ./coreader 200727 0 190(date start end)")
	}
	fileList := make([]string, 0)
	for i := start; i < end; i++ {
		s := fmt.Sprintf(format, i)
		fileList = append(fileList, s)
	}
	return fileList
}
