package configuration

import (
	"fmt"
	"log"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
)

func EurekaClientConfig() {

	client := eureka.NewClient([]string{"http://localhost:8761/eureka"})
	fmt.Println("Printing Client Details..")
	fmt.Println(client)
	instance := eureka.NewInstanceInfo(
		"paypal-host",
		"PAYPAL-SERVICE",
		"127.0.0.1",
		8443,
		30,
		true,
	)

	client.RegisterInstance("PAYPAL-SERVICE", instance)
	go func() {
		ticker := time.NewTicker(25 * time.Second) // send before lease expires
		defer ticker.Stop()

		for range ticker.C {
			err := client.SendHeartbeat(instance.App, instance.HostName)
			if err != nil {
				log.Printf("Heartbeat failed: %v", err)
			} else {
				log.Printf("Heartbeat sent for %s (%s)", instance.App, instance.HostName)
			}
		}
	}()

	fmt.Println("Printing Instance Details...")
	fmt.Println(client.GetInstance(instance.App, instance.InstanceID))
}
