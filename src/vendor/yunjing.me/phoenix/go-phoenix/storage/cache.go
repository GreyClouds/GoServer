package storage

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/go-redis/redis"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

// 缓存管理
type CacheManager struct {
	client *redis.Ring
}

func NewCacheManager(addrs []string, idx int, passwd string) *CacheManager {
	results := make(map[string]string)
	for i, v := range addrs {
		results[fmt.Sprintf("cluster#%d", i+1)] = v
	}

	client := redis.NewRing(&redis.RingOptions{
		Addrs:    results,
		DB:       idx,
		Password: passwd,
	})

	err := client.Ping().Err()
	if err != nil {
		log.Printf("[Phoenix]连接缓存Redis时出错: %v", err)
		return nil
	}

	log.Printf("[Phoenix]数据缓存启用: %v | %d", addrs, idx)

	return &CacheManager{
		client: client,
	}
}

// ---------------------------------------------------------------------------

// 转换键名
func (self *CacheManager) to(name string) string {
	return fmt.Sprintf("cache:%s", name)
}

// 清空缓存数据
func (self *CacheManager) Clean(k string) error {
	err := self.client.Del(self.to(k)).Err()
	if err != nil {
		log.Printf("[Phoenix]缓存数据清除时出错: %v", err)
		return err
	}

	return nil
}

// 加载缓存数据
func (self *CacheManager) Load(k string, v interface{}) error {
	raw, err := self.client.Get(self.to(k)).Bytes()
	if err != nil {
		log.Printf("[Phoenix]缓存数据加载时读取数据出错: %v", err)
		return err
	}

	err = msgpack.Unmarshal(raw, v)
	if err != nil {
		log.Printf("[Phoenix]缓存数据加载时解析数据出错: %v", err)
		return err
	}

	return nil
}

// 保存缓存数据
func (self *CacheManager) Save(k string, v interface{}) error {
	raw, err := msgpack.Marshal(v)
	if err != nil {
		log.Printf("[Phoenix]缓存数据保存时生成数据出错: %v", err)
		return err
	}

	err = self.client.Set(self.to(k), raw, 0).Err()
	if err != nil {
		log.Printf("[Phoenix]缓存数据保存时设置数据出错: %v", err)
		return err
	}

	log.Printf("[Phoenix]缓存%s数据: %d字节", k, len(raw))

	return nil
}

// ---------------------------------------------------------------------------

// 加载二进制数据
func (self *CacheManager) LoadRaw(k string) ([]byte, error) {
	return self.client.Get(self.to(k)).Bytes()
}

// 保存二进制数据
func (self *CacheManager) SaveRaw(k string, raw []byte) error {
	return self.client.Set(self.to(k), raw, 0).Err()
}

// ---------------------------------------------------------------------------

func (self *CacheManager) Client() *redis.Ring {
	return self.client
}

// 生成唯一编号
func (self *CacheManager) GenID(key string, step int32) uint32 {
	keyv := fmt.Sprintf("id:%s", key)
	stepv := int32(1)
	if step > 0 {
		stepv += rand.Int31n(step)
	}

	v := self.client.IncrBy(keyv, int64(stepv)).Val()
	if v == 0 {
		return 0
	}

	return uint32(v)
}

// ---------------------------------------------------------------------------
