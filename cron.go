package litecron

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/robfig/cron"
	"gopkg.in/redsync.v1"
)

const (
	DefaultMutexPrefix = "litecron/defaults"
	DefaultMutexFator = 0.05
)

type handle func()

func (h handle) Name() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

type Handler struct {
	cron string
	handle handle
}

func newHandler(cron string,f handle) *Handler {
	return &Handler{
		cron: cron,
		handle:f,
	}
}

type RedisConfig struct {
	DNS string
}

type MutexConfig struct {
	RedisConfig *RedisConfig
	Prefix string
	Factor float64
}

type Cron struct {
	cronClient *cron.Cron
	sync *redsync.Redsync
	MutexConfig *MutexConfig
}

func NewCron(config *MutexConfig) *Cron {
	c := new(Cron)

	p := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",config.RedisConfig.DNS)
		},
	}
	var pools []redsync.Pool
	pools = append(pools,p)
	c.sync = redsync.New(pools)

	if config.Prefix == "" {
		config.Prefix = DefaultMutexPrefix
	}
	if config.Factor <= 0 {
		config.Factor = DefaultMutexFator
	}
	c.MutexConfig = config
	c.cronClient = cron.New()
	return c
}

func (c *Cron) Register(cronScheue string,h handle) {
	c.cronClient.AddFunc(cronScheue,wrapperHandle(c,newHandler(cronScheue,h)))
}

func (c *Cron) Run() {
	c.cronClient.Run()
}

func (c *Cron) lock(h *Handler) (bool,error) {
	schedule, err := cron.Parse(h.cron)
	if err != nil {
		return false,err
	}
	now := time.Now()
	d := schedule.Next(now).Sub(now)
	d = d - time.Duration(float64(d)*c.MutexConfig.Factor)
	mutex := c.sync.NewMutex(fmt.Sprintf("%s/%s",c.MutexConfig.Prefix,h.handle.Name()),redsync.SetExpiry(d),redsync.SetTries(1))
	if err := mutex.Lock(); err != nil {
		return false,err
	}
	log.Printf("[LiteCron][Info] job will locking still:%s %s %s\n",h.cron,h.handle.Name(),schedule.Next(now))
	return true,nil
}

func wrapperHandle(c *Cron,h *Handler)  handle {
	return func() {
		log.Printf("[LiteCron] start run job:%s %s\n",h.cron,h.handle.Name())
		s,err := c.lock(h)
		if err != nil {
			log.Printf("[LiteCron][Error] can't run job:%s %s %s\n",h.cron,h.handle.Name(),err.Error())
			return
		}
		if !s {
			log.Printf("[LiteCron][Info] job done with other processor:%s %s\n",h.cron,h.handle.Name())
			return
		}
		h.handle()
		log.Printf("[LiteCron] job done:%s %s",h.cron,h.handle.Name())
	}
}