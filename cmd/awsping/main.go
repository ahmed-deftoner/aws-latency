package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strings"

	"time"

	"math/rand"

	"fmt"
	"log"
	"strconv"

	"github.com/ekalinin/awsping"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func GetLatency(g awsping.Graph) (float64, error) {

	cross, err := awsping.ReadIntoGraph()
	total := 0.0
	if err != nil {
		log.Fatal(err)
	}

	max := 0
	for k := range g.Edges {
		x, err := strconv.Atoi(k)
		if err != nil {
			log.Fatal(err)
			break
		}
		if x >= max {
			max = x
		}
	}

	for k, v := range g.Edges {
		_, err := strconv.Atoi(k)
		if err != nil {
			log.Fatal(err)
			return 0, err
		}
		to := strings.Split(v.To, "#")
		if v.From == "*" {
			service := strings.ToLower(to[0])
			r := g.Nodes[v.To].Region
			flag.Parse()

			regions := awsping.GetRegions()

			rand.Seed(time.Now().UnixNano())

			awsping.CalcLatency(regions, 1, true, false, service)
			for _, i := range regions {
				if i.Code == r {
					total += i.GetLatency()
				}
			}
		} else {
			from := strings.Split(v.From, "#")
			reg1 := g.Nodes[from[0]].Region
			reg2 := g.Nodes[to[0]].Region
			total += cross[reg1][reg2]
		}
	}
	fmt.Println(total)
	return total, nil
}

func gethandler(w http.ResponseWriter, r *http.Request) {
	g := awsping.Graph{}
	w.Header().Set("Content-Type", "application/json")
	res := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		res["status"] = "Bad Request"
		res["msg"] = "Unable to parse json"
		b, _ := json.Marshal(res)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
	}
	latency, _ := GetLatency(g)
	res["status"] = "OK"
	res["msg"] = latency
	b, _ := json.Marshal(res)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(b)
}

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,                                                        // All origins
		AllowedMethods:   []string{"GET", "OPTIONS", "POST", "DELETE", "HEAD", "PUT"}, // Allowing only get, just an example
	})

	r := mux.NewRouter()
	r.HandleFunc("/", gethandler).Methods("POST")

	log.Println("Server started on 8080")
	http.ListenAndServe(":8080", c.Handler(r))
}
