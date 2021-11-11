package main

import (
	"github.com/boltdb/bolt"
	"log"
)

func main() {
	// 1. 打开数据库
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Panic("open db failed")
	}
	defer db.Close()

	// 2. 找到抽屉bucket(如果没有就创建)
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("b1"))
			if err != nil {
				log.Panic("create bucket failed")
			}
		}

		bucket.Put([]byte("111"), []byte("hello"))
		bucket.Put([]byte("222"), []byte("world"))

		return nil
	})

	// 3. 写数据
	// 4. 读数据
}
