package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"

	"telecomx-provisioning-service/internal/application/service"
	"telecomx-provisioning-service/internal/config"
	"telecomx-provisioning-service/internal/infrastructure/adapter/event/listener"
	"telecomx-provisioning-service/internal/infrastructure/adapter/repository"
	"telecomx-provisioning-service/internal/infrastructure/adapter/rest"
)

func main() {
	cfg := config.InstanceConfig()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("telecomx_provisioning")

	repo := repository.NewMongoRepository(db)
	svc := service.NewProvisioningService(repo)

	listener.StartKafkaListener(svc, cfg.Brokers, cfg.Topic, cfg.Group, cfg.Client)

	mux := http.NewServeMux()
	rest.NewProvisioningHandler(svc).RegisterRoutes(mux)

	fmt.Printf("REST server running on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
