package service

import (
	"fmt"
	"time"

	"github.com/GerogeGol/yadro-test-problem/domain/store"
)

type TableInfo struct {
	Number      int
	Profit      float64
	WorkingTime time.Duration
}

func (i TableInfo) String() string {
	hours := int(i.WorkingTime.Hours())
	minutes := int(i.WorkingTime.Minutes()) - 60*hours
	time := store.NewDayTime(hours, minutes)
	return fmt.Sprintf("%d %.f %s", i.Number, i.Profit, time)
}
