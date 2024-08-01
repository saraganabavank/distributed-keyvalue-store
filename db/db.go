package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DB struct {
	db         []Shard
	noOfShards int
	IsReplica  bool
	port       int
	host       string
	Replicas   []string
}

func GetDB(noOfShards int, IsReplica bool, port int, host string) *DB {
	var db []Shard
	for index := 0; index < noOfShards; index++ {
		shard, err := GetShard("schema/shard-"+strconv.Itoa(index)+"-"+strconv.Itoa(port), "shard-"+strconv.Itoa(index), index)
		if err != nil {
			panic(err)
		}
		db = append(db, *shard)
	}

	return &DB{
		db:         db,
		noOfShards: noOfShards,
		IsReplica:  IsReplica,
		port:       port,
		host:       host,
		Replicas:   make([]string, 0),
	}

}
func (db *DB) CloseDB() {
	for _, x := range db.db {
		x.db.Close()
	}
}

func (dbs *DB) Get(key string) (string, error, int64) {
	shardIndex := getShardIndex(key, dbs.noOfShards)
	data, err := dbs.db[shardIndex].Get(key)
	return data, err, shardIndex
}

func (dbs *DB) Set(key string, value string) (string, error, int64) {
	shardIndex := getShardIndex(key, dbs.noOfShards)
	if !dbs.IsReplica || true {
		data, err := dbs.db[shardIndex].Set(key, value)
		if !dbs.IsReplica {
			SetOnReplica(dbs, key, value)
		}
		return data, err, shardIndex
	} else {
		return "", fmt.Errorf("Cannot modify connent on read replica."), shardIndex
	}

}

type KeyValue struct {
	Key   string
	Value string
}

func SetOnReplica(dbs *DB, key string, value string) (string, error, int64) {
	for _, replica := range dbs.Replicas {
		jsonData, _ := json.Marshal(&KeyValue{
			Key:   key,
			Value: value,
		})
		resp, _ := http.Post(replica+"/", "application/json", bytes.NewBuffer(jsonData))
		resp.Body.Close()
	}
	return "", nil, 0
}

func DeleteOnReplica(dbs *DB, key string) (string, error, int64) {
	for _, replica := range dbs.Replicas {
		req, _ := http.NewRequest("DELETE", replica+"/"+key, bytes.NewBuffer(nil))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
	}
	return "", nil, 0
}

func (dbs *DB) Delete(key string) (string, error, int64) {
	shardIndex := getShardIndex(key, dbs.noOfShards)

	fmt.Println("+++++++++++++++++++++++", key)
	data, err := dbs.db[shardIndex].Delete(key)
	fmt.Println("+++++++++++++++++++++++", data)
	if !dbs.IsReplica {
		DeleteOnReplica(dbs, key)
	}
	return data, err, shardIndex
}
