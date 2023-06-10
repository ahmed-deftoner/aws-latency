package awsping

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Node struct {
	Region     string                 `json:"region"`
	Properties map[string]interface{} `json:"properties"`
}

type Graph struct {
	Adjacency map[string][]string `json:"adjacency"`
	Nodes     map[string]Node     `json:"nodes"`
	Edges     map[string]Edge     `json:"edges"`
}

func DecodeJSON(fileNmae string) Graph {
	jsonFile, err := os.Open(fileNmae)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result Graph
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

func ReadIntoGraph() (map[string]map[string]float64, error) {
	var regions [22][22]float64
	graph := make(map[string]map[string]float64)

	var regionNames [22]string
	txtFile, err := os.Open("cross-region.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened cross-regio.txt")
	defer txtFile.Close()
	body, err := ioutil.ReadAll(txtFile)
	if err != nil {
		return map[string]map[string]float64{}, err
	}
	raw := strings.Split(string(body), "\r\n")
	count := 0
	for i := 0; i < len(raw); i++ {
		count++
		value := raw[i]
		x := strings.Split(string(value), " ")
		regionNames[i] = x[0]
		for j := 0; j < len(x)-1; j++ {
			regions[i][j], err = strconv.ParseFloat(x[j+1], 64)
			if err != nil {
				return map[string]map[string]float64{}, err
			}
		}
	}

	for i, v := range regionNames {
		graph[v] = make(map[string]float64)
		for j := 0; j < 22; j++ {
			graph[v][regionNames[j]] = regions[i][j]
		}
	}
	return graph, nil
}
