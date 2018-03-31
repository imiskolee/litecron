package litecron

var defaultCron *Cron

func InitDefaultCron(config *MutexConfig) *Cron {
	if defaultCron != nil {
		panic("[LiteCron][Error] defaultCron init twice.")
	}
	defaultCron = NewCron(config)
	return defaultCron
}

func Register(c string,f handle) {
	if defaultCron == nil {
		panic("[LiteCron][Error] can't register cron before InitDefaultCron")
	}
	defaultCron.Register(c,f)
}

func Run() {
	if defaultCron == nil {
		panic("[LiteCron][Error] can't run cron before InitDefaultCron")
	}
	defaultCron.Run()
}



