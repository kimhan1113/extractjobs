package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cheggaaa/pb"
	"github.com/sqweek/dialog"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	dir, _ := dialog.Directory().Browse()

	count := 0
	files, err := ioutil.ReadDir(dir)

	// fmt.Println(len(files))
	// os.Exit(3)

	fmt.Println("hi")

	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup

	file_length := len(files)
	bar := pb.StartNew(file_length)
	wg.Add(len(files))
	go func() {

		for _, f := range files {
			defer wg.Done()
			// 파일 경로 설정
			outFile := filepath.Join(dir, f.Name())
			str := ".xlsx"

			// 엑셀파일만 필터
			if strings.HasSuffix(f.Name(), str) {
				// fmt.Println(outFile)
				xlsx, err := excelize.OpenFile(outFile)

				if err != nil {
					fmt.Println(err)
					return
				}

				activeSheet := xlsx.GetActiveSheetIndex()
				activeSheetName := xlsx.GetSheetName(activeSheet)

				rows, _ := xlsx.GetRows(activeSheetName)
				count += len(rows) - 1
				bar.Increment()

				// time.Sleep(time.Second * 1)
			}
		}
	}()
	wg.Wait()
	fmt.Println("총 rows수는: ", count)

}
