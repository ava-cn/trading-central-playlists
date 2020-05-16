package consoles

import "github.com/robfig/cron"

func InitCorn() {
	var (
		c   *cron.Cron
		err error
	)

	c = cron.New()
	if err = c.AddFunc("* * * * * *", RunSync); err != nil {
		panic(err)
	}

	go c.Start()
	defer c.Stop()

	//select {
	//case <-time.After(time.Second * 10):
	//	return
	//}

}
