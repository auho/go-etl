package database

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func getFieldAmount(field string, value interface{}, t *testing.T) int64 {
	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s` WHERE `%s` = ?", tableName, field)
	res, err := db.QueryFieldInterface("_count", query, value)
	if err != nil {
		t.Error(err)
	}

	return res.(int64)
}

func getAmount() int {
	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s`", tableName)
	res, err := db.QueryFieldInterface("_count", query)
	if err != nil {
		return 0
	}

	n, err := strconv.Atoi(string(res.([]uint8)))
	if err != nil {
		return 0
	}

	return n
}

func Test_DbTargetInsertSliceSlice(t *testing.T) {
	d := NewDbTargetInsertSliceSlice(targetConfig, []string{"name", "value"}, WithDbTargetSliceTruncate())
	d.Start()

	maxI := rand.Intn(100) + 5
	maxJ := rand.Intn(100) + 5

	nameValue := fmt.Sprintf("name slice %d", time.Now().Unix()/(rand.Int63n(999901)+1))
	for i := 0; i < maxI; i++ {
		items := make([][]interface{}, 0, maxJ)
		for j := 0; j < maxJ; j++ {
			items = append(items, []interface{}{nameValue, time.Now().Unix()})
		}

		d.Send(items)
	}

	d.Done()
	d.Close()
	//d.State()

	count := getFieldAmount("name", nameValue, t)
	if int64(maxI*maxJ) != count {
		t.Error(fmt.Sprintf("count is error. %d != %d", maxI*maxJ, count))
	}

	if d.State.itemAmount != count {
		t.Error(fmt.Sprintf("count is error. %d != %d", d.State.itemAmount, count))
	}
}

func Test_DbTargetInsertSliceMap(t *testing.T) {
	d := NewDbTargetInsertSliceMap(targetConfig, WithDbTargetMapTruncate())
	d.Start()

	maxI := rand.Intn(100) + 5
	maxJ := rand.Intn(100) + 5

	nameValue := fmt.Sprintf("name map %d", time.Now().Unix()/(rand.Int63n(999902)+1))
	for i := 0; i < maxI; i++ {
		items := make([]map[string]interface{}, 0, maxJ)
		for j := 0; j < maxJ; j++ {
			items = append(items, map[string]interface{}{"name": nameValue, "value": time.Now().Unix()})
		}

		d.Send(items)
	}

	d.Done()
	d.Close()

	count := getFieldAmount("name", nameValue, t)
	if int64(maxI*maxJ) != count {
		t.Error(fmt.Sprintf("count is error. %d != %d", maxI*maxJ, count))
	}

	if d.State.itemAmount != count {
		t.Error(fmt.Sprintf("count is error. %d != %d", d.State.itemAmount, count))
	}
}

func Test_DbTargetUpdateSliceMap(t *testing.T) {
	d := NewDbTargetUpdateSliceMap(targetConfig, pkName)
	d.Start()

	maxI := rand.Intn(10) + 1
	maxJ := rand.Intn(10) + 1
	value := rand.Intn(1000) + 3

	amount := 0
	for i := 0; i < maxI; i++ {
		items := make([]map[string]interface{}, 0, maxJ)
		for j := 0; j < maxJ; j++ {
			amount += 1
			items = append(items, map[string]interface{}{"id": amount, "value": value})
		}

		d.Send(items)
	}

	d.Done()
	d.Close()

	count := getFieldAmount("value", value, t)
	if int64(maxI*maxJ) != count {
		t.Error(fmt.Sprintf("count is error. %d != %d", maxI*maxJ, count))
	}

	if d.State.itemAmount != count {
		t.Error(fmt.Sprintf("count is error. %d != %d", d.State.itemAmount, count))
	}
}

func Test_DbSource(t *testing.T) {
	d := NewDbSource(sourceConfig)
	d.Start()

	amount := getAmount()

	count := 0
	page := 0
	for {
		items, ok := <-d.itemsChan
		if ok == false {
			break
		}

		if len(items) <= 0 {
			t.Error(fmt.Sprintf("items size is error %d != %d", len(items), sourceConfig.Size))
		}

		page += 1
		count += len(items)
	}

	if page > sourceConfig.Page {
		t.Error(fmt.Sprintf("page is error %d != %d", page, sourceConfig.Page))
	}

	if amount != count {
		t.Error(fmt.Sprintf("amount is error %d != %d", amount, count))
	}
}
