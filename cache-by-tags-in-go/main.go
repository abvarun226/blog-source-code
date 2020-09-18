package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

func main() {
	client := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		},
	)

	app := cache{Client: client}

	k := "data:key1"
	v := "randomstringdata"

	tags := make([]string, 0)
	for i := 1; i <= 60; i++ {
		tags = append(tags, fmt.Sprintf("post%d", i))
	}
	app.SetByTags(k, v, 30*time.Minute, tags)

	k = "data:key2"
	tags = nil
	for i := 10; i <= 60; i++ {
		tags = append(tags, fmt.Sprintf("post%d", i))
	}
	app.SetByTags(k, v, 30*time.Minute, tags)

	k = "data:key3"
	tags = nil
	for i := 50; i <= 60; i++ {
		tags = append(tags, fmt.Sprintf("post%d", i))
	}
	app.SetByTags(k, v, 30*time.Minute, tags)

	tags = nil
	for i := 1; i <= 10; i++ {
		tags = append(tags, fmt.Sprintf("post%d", i))
	}
	app.Invalidate(tags)
}

type cache struct {
	client redis.Client
}

// SetByTags will set cache by given tags.
func (c *cache) SetByTags(key, value string, expiry time.Duration, tags []string) {
	t := time.Now()

	pipe := c.Client.TxPipeline()
	for _, tag := range tags {
		pipe.SAdd("comment_by_tags:"+tag, key)
	}
	pipe.Set(key, value, expiry)

	if _, err := pipe.Exec(); err != nil {
		log.Printf("error in pipeline: %v", err)
	}

	log.Printf("SetByTags: time take = %dms", time.Since(t).Milliseconds())
}

// Invalidate will invalidate cache with given tags.
func (c *cache) Invalidate(tags []string) {
	t := time.Now()
	keys := make([]string, 0)
	for _, tag := range tags {
		tagKey := "comment_by_tags:" + tag
		k, _ := c.Client.SMembers(tagKey).Result()
		keys = append(keys, tagKey)
		keys = append(keys, k...)
	}
	c.Client.Del(keys...)
	log.Printf("Invalidate: time take = %dms", time.Since(t).Milliseconds())
}
