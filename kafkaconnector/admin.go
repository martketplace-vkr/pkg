package kafkaconnector

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/IBM/sarama"
	"github.com/martketplace-vkr/pkg/logger/log"
)

type (
	AdminConfig struct {
		Username string
		Password string
		Brokers  []string
		SASL     *SASL
	}
	Admin struct {
		sarama.ClusterAdmin
	}
)

func NewAdmin(cfg AdminConfig) *Admin {
	config, err := getConfig(cfg)
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}

	admCtrl, err := sarama.NewClusterAdmin(cfg.Brokers, config)
	if err != nil {
		log.Fatalf("error creating cluster admin: %v", err)
	}

	return &Admin{
		ClusterAdmin: admCtrl,
	}
}

func getConfig(crds AdminConfig) (*sarama.Config, error) {
	config := sarama.NewConfig()

	if crds.SASL != nil {
		tlsConfig, err := setupTLS(crds.SASL)
		if err != nil {
			return nil, err
		}

		config.Net.TLS.Config = tlsConfig
	}

	config.Version = sarama.V2_6_0_0
	config.Net.TLS.Enable = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = crds.Username
	config.Net.SASL.Password = crds.Password
	config.Producer.Return.Successes = true

	return config, nil
}

func setupTLS(sasl *SASL) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	if sasl.CaPath != nil {
		caCert, err := os.ReadFile(*sasl.CaPath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.RootCAs = caCertPool
	}

	tlsConfig.InsecureSkipVerify = true

	return &tlsConfig, nil
}
