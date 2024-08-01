package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Raft struct {
	Master   string
	Replicas []string
	WhoAm    string
}

func GetRaft(master string, iam string) *Raft {
	return &Raft{
		Master:   master,
		Replicas: make([]string, 0),
		WhoAm:    iam,
	}
}

func (raft *Raft) HeartBeat() {

	for {
		if raft.Master == "" {
			// master get lost  first replica become the master
			if len(raft.Replicas) > 0 {
				raft.Master = raft.Replicas[0]

			}
			var replicas []string
			for _, replica := range raft.Replicas {
				if raft.Master != replica {
					replicas = append(replicas, replica)
				}
			}
			raft.Replicas = replicas

			jsonData, _ := json.Marshal(Raft{
				Replicas: raft.Replicas,
				Master:   raft.Master,
			})
			for _, x := range raft.Replicas {
				fmt.Println(x + "/electmaster")
				http.Post(x+"/electmaster", "application/json", bytes.NewBuffer(jsonData))

			}

		} else if raft.WhoAm == raft.Master {
			// check replicas can be reachable
			for index, raplica := range raft.Replicas {
				jsonData, _ := json.Marshal(raft.Replicas)
				_, err := http.Post(raplica+"/replica", "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					if len(raft.Replicas) != 0 {
						if index+1 != len(raft.Replicas) {
							raft.Replicas = append(raft.Replicas[:index], raft.Replicas[index+1])
						} else {
							raft.Replicas = raft.Replicas[:index]
						}
					}
				}

			}
			time.Sleep(100 * time.Millisecond)

		} else {
			// check master can be reachable
			jsonData, _ := json.Marshal(raft.Replicas)
			_, err := http.Post(raft.Master+"/replica?replica="+raft.WhoAm, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				raft.Master = ""
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}

}
