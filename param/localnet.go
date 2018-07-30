package param

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/EducationEKT/EKT/p2p"
)

var LocalNet = []p2p.Peer{
	{"8d9ba40d8d2638386fb9bcea9e92cf636b044bb409486dd88c2ae104a42de348", "127.0.0.1", 19951, 4, ""},
	{"053d2c7c6c29571916264261ffa53d74a2e44753ee029f6fbb9ac43d953bfebe", "127.0.0.1", 19952, 4, ""},
	{"416d5a9691c1893fec49189ca2b94b773cd297448100557b2f85dd1c4cedb230", "127.0.0.1", 19953, 4, ""},
	{"993ab6fe4eedcae2723381f01c89766e6dc4c3c2f5f66aba47b4d336d16ff32b", "127.0.0.1", 19954, 4, ""},
	{"4fb83a25a798cfc1c5d5c55c7a06bedf02c5033b4bacf8bb9ab3b2f95462d331", "127.0.0.1", 19955, 4, ""},
	{"46675f25112ad2b2008fbdf2c66d6403edc27308efba8bdf7c4d68fdc0875ea6", "127.0.0.1", 19956, 4, ""},
	{"3d26f85cb8cd45e6e89b3ef1574982e75a7a8f8e29206e03eab041f92754306f", "127.0.0.1", 19957, 4, ""},
}

func init() {
	loadLocalNet()
}

func loadLocalNet() {
	cfg := "localnet.json"
	data, err := ioutil.ReadFile(cfg)
	if err != nil {
		return
	}
	log.Println("Found localnet.json, loading it")
	peers := [][]interface{}{}
	err = json.Unmarshal(data, &peers)
	if err != nil {
		log.Println("Invalid localnet.json format, ingore it")
		return
	}
	net := []p2p.Peer{}
	for _, peer := range peers {
		if len(peer) != 5 {
			fmt.Println("Invalid localnet.json format, ingore it")
			return
		}
		net = append(net, p2p.Peer{peer[0].(string), peer[1].(string), int32(peer[2].(float64)), int(peer[3].(float64)), peer[4].(string)})
	}
	LocalNet = net
	log.Println("Using localnet.json")
}
