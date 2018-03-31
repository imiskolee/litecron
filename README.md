### LiteCron
------

#### What's LiteCron

LiteCron is a In-Processing distributed cron job processor. its easy to handle you cron job into you web app.
you don't need any special machine for run cron job.

#### UseCase

* running cron job into my web app.
* build a distributed cron job service(for replace system cron service).


#### UseAge

```go
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
```

#### Examples

1. [default cron client](example_test.go#L14)
2. [without default cron client](example_test.go#L30)
3. [mock multi processor](example_test.go#L44)

#### Deps & Thanks

* [RedSync - Distributed mutual exclusion lock ](https://github.com/go-redsync/redsync)
* [Cron - cron lib](https://github.com/robfig/cron)


