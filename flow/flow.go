package flow

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/storage/database"
)

func RunFlow(config goEtl.DbConfig, dataName string, idName string, actions []action.Actionor) {
	var wg sync.WaitGroup

	fields := []string{idName}
	for _, a := range actions {
		fields = append(fields, a.GetFields()...)
	}

	fields = goEtl.RemoveReplicaSliceString(fields)

	sourceConfig := database.NewDbSourceConfig()
	sourceConfig.MaxConcurrent = runtime.NumCPU() * 2
	sourceConfig.Size = 2000
	sourceConfig.Table = dataName
	sourceConfig.Driver = config.Driver
	sourceConfig.Dsn = config.Dsn
	sourceConfig.PKeyName = idName
	sourceConfig.Fields = fields

	source := database.NewDbSource(sourceConfig)
	source.Start()

	for _, a := range actions {
		a.Start()
	}

	startTime := time.Now()
	fmt.Println("start...", startTime.Format("2006-01-02 15:04:05"))
	anchorTicker := time.NewTicker(time.Millisecond * 500)

	go func() {
		fmt.Println(source.State.GetTitle())
		fmt.Println(" ")

		for _, a := range actions {
			fmt.Println(a.GetTitle())
			fmt.Println(" ")
		}

		lines := 2 + len(actions)*2

		for range anchorTicker.C {
			fmt.Printf("%c[%dA\r%c[K%c[1;40;32m %s %c[0m", 0x1B, lines-1, 0x1B, 0x1B, source.State.GetStatus(), 0x1B)
			for _, a := range actions {
				fmt.Printf("%c[2B\r%c[K%c[1;40;32m %s %c[0m", 0x1B, 0x1B, 0x1B, a.GetStatus(), 0x1B)
			}

			fmt.Printf("%c[1B\r", 0x1B)
		}
	}()

	wg.Add(1)
	go func() {
		for {
			if items, ok := source.Consume(); ok {
				for _, a := range actions {
					a.Receive(items)
				}
			} else {
				break
			}
		}

		wg.Done()
	}()

	wg.Wait()

	for _, a := range actions {
		a.Done()
		a.Close()
	}

	endTime := time.Now()
	time.Sleep(time.Second)
	anchorTicker.Stop()

	duration := endTime.Sub(startTime)
	fmt.Println("done", endTime.Format("2006-01-02 15:04:05"), "耗时:",
		fmt.Sprintf("%d%s %d%s", int(math.Floor(duration.Seconds()/60)), "/分", int(math.Ceil(duration.Seconds()))%60, "/秒"), "\n ")
}
