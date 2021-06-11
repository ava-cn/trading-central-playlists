package consoles

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

func InitCorn() {
	var (
		c   *cron.Cron
		err error
	)

	c = cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
	if _, err = c.AddFunc("@hourly", RunSync); err != nil {
		panic(err)
	}

	go c.Start()
	defer c.Stop()

	//select {
	//case <-time.After(time.Second * 10):
	//	return
	//}

}
