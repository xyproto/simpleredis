package simpleredis

import (
	"github.com/xyproto/pinterface"
)

// For implementing pinterface.ICreator

type RedisCreator struct {
	pool    *ConnectionPool
	dbindex int
}

func NewCreator(pool *ConnectionPool, dbindex int) *RedisCreator {
	return &RedisCreator{pool, dbindex}
}

func (c *RedisCreator) NewList(id string) pinterface.IList {
	return &List{c.pool, id, c.dbindex}
}

func (c *RedisCreator) NewSet(id string) pinterface.ISet {
	return &Set{c.pool, id, c.dbindex}
}

func (c *RedisCreator) NewHashMap(id string) pinterface.IHashMap {
	return &HashMap{c.pool, id, c.dbindex}
}

func (c *RedisCreator) NewKeyValue(id string) pinterface.IKeyValue {
	return &KeyValue{c.pool, id, c.dbindex}
}
