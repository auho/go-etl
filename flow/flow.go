package flow

import (
	"fmt"
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
	sourceConfig.MaxConcurrent = 4
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

	fmt.Println("start...")
	go func() {
		fmt.Println(source.State.GetTitle())
		fmt.Println(" ")

		for _, a := range actions {
			fmt.Println(a.GetTitle())
			fmt.Println(" ")
		}

		lines := 2 + len(actions)*2
		t := time.NewTicker(time.Millisecond * 500)

		for range t.C {
			fmt.Printf("%c[%dA\r%c[K%c[1;40;32m %s %c[0m", 0x1B, lines-1, 0x1B, 0x1B, source.State.GetStatus(), 0x1B)
			for _, a := range actions {
				fmt.Printf("%c[2B\r%c[K%c[1;40;32m %s %c[0m", 0x1B, 0x1B, 0x1B, a.GetStatus(), 0x1B)
			}

			fmt.Printf("%c[1B", 0x1B)
		}
	}()

	wg.Add(1)
	go func() {
		for {
			if items, ok := source.Consume(); ok {
				time.Sleep(time.Second * 3)
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

	time.Sleep(time.Second)
	fmt.Println("done")
}
