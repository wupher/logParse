package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	app := &cli.App{
		Name:  "box_log_count",
		Usage: "计数包厢房态",
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 1 {
				fmt.Printf("请输入房态日志文件名来进行解析 \n")
				_ = cli.Exit("必须提供日志文件路径", 1)
			}
			fileName := c.Args().Get(0)
			fmt.Printf("解析日志文件： %v \n", fileName)
			//TOOD 传参
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func convertLogs(logFiles []string, csv string) {
	initData()
	t := time.Now()

	for _, v := range logFiles {
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
	outPut(csv)
}
