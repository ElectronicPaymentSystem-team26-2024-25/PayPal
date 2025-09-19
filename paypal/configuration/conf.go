package configuration

import (
	"fmt"

	"github.com/ArthurHlt/go-eureka-client/eureka"
)

func EurekaClientConfig() {

	client := eureka.NewClient([]string{"http://localhost:8761/eureka"})
	fmt.Println("Printing Client Details..")
	fmt.Println(client)
	instance := eureka.NewInstanceInfo(
		"localhost",
		"PAYPAL-SERVICE",
		"127.0.0.1",
		443,
		30,
		false,
	)

	client.RegisterInstance("PAYPAL-SERVICE", instance)
	client.SendHeartbeat(instance.App, instance.HostName)
	fmt.Println("Printing Instance Details...")
	fmt.Println(client.GetInstance(instance.App, instance.HostName))
}
