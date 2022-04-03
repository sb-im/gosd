package read

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"sb.im/gosd/app/model"
)

const (
	nodesName = "nodes"
	pointsDir = "points"
)

type dataNode struct {
	ID     uint     `json:"id"`
	Name   string   `json:"name"`
	Points []string `json:"points"`
}

type dataPoints map[string]json.RawMessage

func getNode(path string) []dataNode {
	f, err := os.Open(path)
	if err != nil {
	}

	fNodes := []dataNode{}
	if err := json.NewDecoder(f).Decode(&fNodes); err != nil {
	}
	return fNodes
}

func getPoints(path string) dataPoints {
	filePoints, _ := ioutil.ReadDir(path)
	points := make(dataPoints, len(filePoints))

	for _, v := range filePoints {
		f, err := os.Open(filepath.Join(path, v.Name()))
		if err != nil {
			panic(err)
		}
		d, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		name := strings.TrimSuffix(v.Name(), ".json")
		points[name] = json.RawMessage(d)
	}

	return points
}

func ParseNode(path string) []model.Node {
	// load Node
	fNodes := getNode(filepath.Join(path, nodesName+".json"))
	// load points
	points := getPoints(filepath.Join(path, pointsDir))

	nodes := make([]model.Node, len(fNodes))

	for i, n := range fNodes {

		ps := make([]json.RawMessage, len(n.Points))
		for j, p := range n.Points {
			ps[j] = points[p]
		}
		data, _ := json.Marshal(ps)

		nodes[i].ID = n.ID
		nodes[i].Name = n.Name
		nodes[i].Points = data
	}
	return nodes
}
