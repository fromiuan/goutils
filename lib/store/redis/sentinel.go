package redis

import (
	"net"
	"strings"
	"sync"
	"time"

	"share/logging"

	"github.com/gomodule/redigo/redis"
)

type SentinelPool struct {
	sentinelAddr string
	masterName   string
	masterAddr   string
	db           string
	password     string
	pools        map[string]*redis.Pool
	l            sync.RWMutex
}

func NewSentinelPool(sentinelAddr string, masterName string, masterAddr string, db string, password string) *SentinelPool {
	sc := &SentinelPool{
		masterName:   masterName,
		masterAddr:   masterAddr,
		sentinelAddr: sentinelAddr,
		db:           db,
		password:     password,
		pools:        make(map[string]*redis.Pool),
	}

	if len(sentinelAddr) > 0 {
		go sc.listen()
	}
	return sc
}

func (this *SentinelPool) Get() redis.Conn {
	return this.switchPool().Get()
}

func (this *SentinelPool) ActiveCount() int {
	return this.switchPool().ActiveCount()
}

// ------------------------------------------------------------------------

// Maseter and slave switch
func (this *SentinelPool) switchPool() *redis.Pool {

	this.l.Lock()
	defer this.l.Unlock()

	if _, ok := this.pools[this.masterAddr]; !ok {
		this.pools[this.masterAddr] = &redis.Pool{
			MaxIdle:     80,
			MaxActive:   10000,
			IdleTimeout: 600 * time.Second,
			Dial: func() (redis.Conn, error) {
				con, err := redis.Dial("tcp", this.masterAddr)
				if err != nil {
					return nil, err
				}
				_, err = con.Do("AUTH", this.password)
				if err == nil {
					con.Do("SELECT", this.db)
				}

				return con, err
			},
		}
	}

	return this.pools[this.masterAddr]

}

// Monitor sentinel message
func (this *SentinelPool) listen() {
	conn, _ := redis.Dial("tcp", this.sentinelAddr)
	pubsub := redis.PubSubConn{Conn: conn}

	if err := pubsub.Subscribe("+switch-master"); err != nil {
		return
	}

	for {

		switch msg := pubsub.Receive().(type) {
		case redis.Message:
			switch msg.Channel {
			case "+switch-master":
				parts := strings.Split(string(msg.Data), " ")
				if parts[0] != this.masterName {
					logging.Info("Redis sentinel: Ignore new %s addr", parts[0])
					continue
				}

				masterAddr := net.JoinHostPort(parts[3], parts[4])
				logging.Info("Redis sentinel: New %q addr is %s", this.masterName, masterAddr)

				this.l.Lock()
				this.masterAddr = masterAddr
				this.l.Unlock()

			default:
				logging.Info("Redis sentinel: Unsupported message: %s", msg)
			}
		case redis.Subscription:
			// Ignore.
		default:
			logging.Warning("%v", msg)
		}
	}
}
