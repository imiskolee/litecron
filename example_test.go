package litecron_test

import (
	"log"
	"time"

	"github.com/imiskolee/litecron"
)

func exampleHandle() {
	log.Print("Hello,LiteCron.")
}

func ExampleDefault() {
	//Step 1: init default cron client.
	litecron.InitDefaultCron(&litecron.MutexConfig{
		RedisConfig:&litecron.RedisConfig{
			DNS:"127.0.0.1:6379",
		},
		Prefix:"litecron/examples/defaults/",
		Factor:0.01,
	})
	//Register a cron job for every 2 seconds.
	litecron.Register("@every 2s",exampleHandle)
	go litecron.Run()
	time.Sleep(10 * time.Second)
	// Output:
}

func ExampleNew() {
	cron := litecron.NewCron(&litecron.MutexConfig{
		RedisConfig:&litecron.RedisConfig{
			DNS:"127.0.0.1:6379",
		},
		Prefix:"litecron/examples/advances/",
		Factor:0.01,
	})
	cron.Register("@every 5s",exampleHandle)
	go cron.Run()
	time.Sleep(10 * time.Second)
	// Output:
}

func ExampleMulti() {
	f := func() {
		cron := litecron.NewCron(&litecron.MutexConfig{
			RedisConfig: &litecron.RedisConfig{
				DNS: "127.0.0.1:6379",
			},
			Prefix: "litecron/examples/multi/",
			Factor: 0.01,
		})
		cron.Register("@every 3s", exampleHandle)
		go cron.Run()
	}
	go f() // processor 1
	go f() // processor 2
	time.Sleep(10 * time.Second)
	// Output:
}

