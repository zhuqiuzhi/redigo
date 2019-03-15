package main

import (
    "log"
    "time"

    "github.com/gomodule/redigo/redis"
    "github.com/youtube/vitess/go/pools"
    "golang.org/x/net/context"
)

// ResourceConn adapts a Redigo connection to a Vitess Resource.
type ResourceConn struct {
    redis.Conn
}

func (r ResourceConn) Close() {
    r.Conn.Close()
}

func main() {
    p := pools.NewResourcePool(func() (pools.Resource, error) {
        c, err := redis.Dial("tcp", ":6379")
        return ResourceConn{c}, err
    }, 1, 2, time.Minute)
    defer p.Close()

    ctx := context.TODO()
    r, err := p.Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer p.Put(r)

    c := r.(ResourceConn)
    n, err := c.Do("INFO")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("info=%s", n)
}
