package main

import (
	"flag"
	"log"

	"github.com/luisfernandogaido/short/api"
)

func main() {
	var (
		addr      *string
		mongoURI  *string
		tokenRoot *string
		redisURI  *string
		domain    *string
	)
	addr = flag.String("addr", ":4018", "porta local")
	mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "string de conexão mongoDB")
	tokenRoot = flag.String("token_root", "MDKTQuSFiUleok6GA0MF", "token administrativo, crud usuários autorizados")
	redisURI = flag.String("redis_uri", "redis://localhost:6379", "string de conexão Redis")
	domain = flag.String("domain", "http://localhost:4018", "domínio que hospeda o serviço")
	flag.Parse()
	if err := api.Serve(*addr, *mongoURI, *tokenRoot, *redisURI, *domain); err != nil {
		log.Fatal(err)
	}
}
