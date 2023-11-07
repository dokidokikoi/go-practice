package service

import (
	"fmt"
	"time"

	"github.com/liushuochen/gotable"
	"github.com/liushuochen/gotable/table"
)

func PrintData() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			printProcessData()
		}
	}
}

func printProcessData() {
	fmt.Print("\033[2J") // 清空终端输出
	fmt.Print("\033[H")  // 光标移动到屏幕左上角
	tableMeta, err := gotable.CreateByStruct(&record{})
	if err != nil {
		fmt.Println(err)
	}
	dataCenter.lock.RLock()
	defer dataCenter.lock.RUnlock()

	output := map[string]string{
		"Mid":       dataCenter.records[ProcessKey].Mid,
		"Reply":     dataCenter.records[ProcessKey].Reply,
		"DeviceId":  dataCenter.records[ProcessKey].DeviceId,
		"Timestamp": time.Unix(dataCenter.records[ProcessKey].Timestamp/1000, 0).Format("2006-01-02 15:04:05"),
	}
	tableMeta.AddRow(output)
	fmt.Println(tableMeta)

	table := CreateTable(dataCenter.records[ProcessKey].Add)
	fmt.Println("Add")
	fmt.Println(table)

	table = CreateTable(dataCenter.records[ProcessKey].Del)
	fmt.Println("Del")
	fmt.Println(table)

	table = CreateTable(dataCenter.records[ProcessKey].Update)
	fmt.Println("Update")
	fmt.Println(table)

	process := make([]Data, 0)
	for _, d := range dataCenter.data[ProcessKey] {
		process = append(process, d)
	}
	table = CreateTable(process)
	fmt.Println("Full")
	fmt.Println(table)
	fmt.Printf("总进程数：%d", len(process))
}

func CreateTable(process []Data) *table.Table {
	table, err := gotable.CreateByStruct(&Process{})
	if err != nil {
		fmt.Println(err)
	}
	cols := make([]map[string]string, 0)

	for _, d := range process {
		cols = append(cols, d.ToMap())
	}
	table.AddRows(cols)

	return table
}
