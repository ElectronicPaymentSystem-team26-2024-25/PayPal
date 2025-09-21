package configuration

import (
	"fmt"
	"log"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
)

func EurekaClientConfig() {
	client := eureka.NewClient([]string{"http://localhost:8761/eureka"})

	// Configure instance with HTTPS
	instance := eureka.NewInstanceInfo(
		"127.0.0.1",      // hostname
		"PAYPAL-SERVICE", // app name
		"127.0.0.1",      // IP address
		8443,             // port (HTTPS)
		30,               // lease duration
		true,             // isSSL
	)

	// Override URLs so Eureka publishes HTTPS info
	// Ensure SecurePort struct exists
	if instance.SecurePort == nil {
		instance.SecurePort = &eureka.Port{}
	}
	instance.SecurePort.Port = 8443
	instance.SecurePort.Enabled = true

	// Disable non-secure port
	if instance.Port != nil {
		instance.Port.Enabled = false
	}

	instance.SecureVipAddress = "PAYPAL-SERVICE"
	instance.VipAddress = "PAYPAL-SERVICE"

	instance.StatusPageUrl = fmt.Sprintf("https://%s:%d/info", instance.IpAddr, instance.SecurePort.Port)
	instance.HealthCheckUrl = fmt.Sprintf("https://%s:%d/health", instance.IpAddr, instance.SecurePort.Port)
	instance.HomePageUrl = fmt.Sprintf("https://%s:%d/", instance.IpAddr, instance.SecurePort.Port)

	// Register instance
	client.RegisterInstance("PAYPAL-SERVICE", instance)

	// Send heartbeat before lease expires
	go func() {
		ticker := time.NewTicker(25 * time.Second)
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
