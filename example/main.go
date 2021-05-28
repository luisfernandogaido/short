package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/luisfernandogaido/short/api"
)

func main() {
	var (
		addr              *string
		mongoURI          *string
		tokenRoot         *string
		redisURI          *string
		domain            *string
		authorizedDomains []string
	)
	addr = flag.String("addr", ":4018", "porta local")
	mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "string de conexão mongoDB")
	tokenRoot = flag.String("token_root", "MDKTQuSFiUleok6GA0MF", "token administrativo, crud usuários autorizados")
	redisURI = flag.String("redis_uri", "redis://localhost:6379", "string de conexão Redis")
	domain = flag.String("domain", "http://localhost:4018", "domínio que hospeda o serviço")
	flag.Parse()
	b, err := os.ReadFile("./authorized-domains.json")
	if err != nil {
		log.Fatal("./authorized-domains.json: ", err)
	}
	if err := json.Unmarshal(b, &authorizedDomains); err != nil {
		log.Fatal("./authorized-domains.json: ", err)
	}
	if err := api.Serve(*addr, *mongoURI, *tokenRoot, *redisURI, *domain, authorizedDomains); err != nil {
		log.Fatal(err)
	}
}
