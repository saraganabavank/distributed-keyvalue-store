package main

import (
	"flag"
	"fmt"
	"kvstore/db"
	"kvstore/raft"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Election struct {
	Master   string   `json:"Master"`
	Replicas []string `json:"Replicas"`
}

var (
	port    = flag.Int("port", 8080, "Port number for server")
	shard   = flag.Int("shard", 1, "No of shard for DB")
	replica = flag.Bool("replica", false, "No of shard for DB")
	master  = flag.Int("master", 0, "master node address")
	host    = flag.String("host", "http://localhost:", "Domine in which process are running on")
)

func main() {

	flag.Parse()
	master_ := *host + strconv.Itoa(*master)
	iam := *host + strconv.Itoa(*port)
	if *master == 0 {
		master_ = iam
	}

	raft := raft.GetRaft(master_, iam)
	go raft.HeartBeat()
	router := gin.Default()
	database := db.GetDB(*shard, *replica, *port, *host)

	if !*replica {
		iam = master_
	}

	defer database.CloseDB()
	router.POST("/replica", func(c *gin.Context) {
		replica___ := c.Query("replica")
		isPresent := false
		if *replica && raft.Master != raft.WhoAm {
			var data []string
			err := c.BindJSON(&data)
			if err == nil {
				raft.Replicas = data
				if len(replica___) > 0 {
					for _, replica_ := range raft.Replicas {
						if replica_ == replica___ {
							isPresent = true
						}
					}
					if !isPresent {
						raft.Replicas = append(raft.Replicas, replica___)
					}

					database.IsReplica = true
				}
			}
		} else {
			if len(replica___) > 0 {
				for _, replica_ := range raft.Replicas {
					if replica_ == replica___ {
						isPresent = true
					}
				}
				if !isPresent {
					raft.Replicas = append(raft.Replicas, replica___)
				}
			}
		}

		database.Replicas = raft.Replicas
		c.JSON(http.StatusOK, gin.H{})
	})
	router.GET("/instance", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"master":   raft.Master,
			"replicas": raft.Replicas,
		})
	})
	router.GET("/healtz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"master":   raft.Master,
			"replicas": raft.Replicas,
		})
	})
	router.GET("/", func(c *gin.Context) {
		key := c.Query("key")
		value, err, shardIndex := database.Get(key)
		fmt.Println("Key -", key, "shardIndex - ", shardIndex)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value, "shardIndex": shardIndex})
	})
	router.POST("/", func(c *gin.Context) {
		var keyValuePair KeyValue
		if err := c.BindJSON(&keyValuePair); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		status, err, shardIndex := database.Set(keyValuePair.Key, keyValuePair.Value)
		fmt.Println("Key -", keyValuePair.Key, "shardIndex - ", shardIndex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": status, "shardIndex": shardIndex})
	})
	router.POST("/electmaster", func(c *gin.Context) {
		var election Election
		if err := c.BindJSON(&election); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(election)
		if raft.Master != iam {
			raft.Master = election.Master
			raft.Replicas = election.Replicas
			database.IsReplica = true
			database.Replicas = election.Replicas

		} else {
			database.IsReplica = false
		}

		c.JSON(http.StatusOK, gin.H{})
	})
	router.DELETE("/:key", func(c *gin.Context) {
		key := c.Param("key")
		status, err, shardIndex := database.Delete(key)
		fmt.Println("Key -", key, "shardIndex - ", shardIndex)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": status, "shardIndex": shardIndex})

	})

	address := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting server on port %d...\n", *port)
	router.Run(address)
}
