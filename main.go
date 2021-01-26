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
		Name:  "Daily Report Generator CSV Convert",
		Usage: "转换日报生成 CSV 记录",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "logs",
				Aliases: []string{"L"},
				Usage:   "load logs from `FILES` ",
			},
			&cli.StringFlag{
				Name:    "csv",
				Aliases: []string{"C"},
				Usage:   "output csv filesName",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 1 {
				fmt.Printf("请输入房态日志文件名来进行解析 \n")
				_ = cli.Exit("必须提供日志文件路径", 1)
			}
			logFileNames := c.StringSlice("logs")
			fmt.Printf("解析日志文件： %v \n", logFileNames)
			csvFileName := c.String("csv")
			fmt.Printf("转换为 CSV: %v \n", csvFileName)
			convertLogs(logFileNames, csvFileName)
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
		fmt.Printf("文件 %v 花费时间 - %v 秒\n", v, time.Since(s))
	}
	fmt.Printf("全部文件共花费 - %v 秒\n", time.Since(t))
	outPut(csv)
}
