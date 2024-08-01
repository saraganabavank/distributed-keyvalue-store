package db

import (
	"fmt"
	"hash/fnv"
	"log"

	bbolt "go.etcd.io/bbolt"
)

type Shard struct {
	db         *bbolt.DB
	bucket     string
	shardIndex int
}

func GetShard(path string, bucket string, shardIndex int) (*Shard, error) {
	// dbs := make([]DB,0)
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bbolt.Tx) error {
		// Check if the bucket already exists.
		bucket_ := tx.Bucket([]byte(bucket))
		if bucket_ == nil {
			// Create the bucket if it doesn't exist.
			_, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			fmt.Println("Bucket created:", bucket)
		} else {
			fmt.Println("Bucket already exists:", bucket)
		}
		return nil
	})

	return &Shard{db, bucket, shardIndex}, nil
}

func (db *Shard) Get(key string) (string, error) {
	var value []byte
	err := db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(db.bucket))
		value = bucket.Get([]byte(key))
		return nil
	})

	if err != nil {
		return "", err
	}
	if value == nil {

		return "", fmt.Errorf("Key not exist")
	}

	return string(value), nil
}

func (db *Shard) Set(key string, value string) (string, error) {

	err := db.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(db.bucket))
		err := bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return "failed", err
	}
	return "ok", err
}

func (db *Shard) Delete(key string) (string, error) {
	err := db.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(db.bucket))
		err := bucket.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return "failed", err
	}

	return "ok", nil
}
func getShardIndex(key string, noOfShards int) int64 {
	hash := fnv.New64()
	hash.Write([]byte(key))
	shardIndex := hash.Sum64() % uint64(noOfShards)
	return int64(shardIndex)
}
