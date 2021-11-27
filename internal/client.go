package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

var ErrNoFormatSpecified = errors.New("no format has been specified. Use -o to specify a format")

type Client struct {
	secretsClient coreV1Types.SecretInterface
	dryRun        bool
	namespace     string
	secretName    string
	format        string
	data          []byte
	key           string
}

func NewClient(dryRun bool, namespace, secretName, format, key string, data []byte) (*Client, error) {
	// NewNonInteractiveDeferredLoadingClientConfig is being used, because this way we respect the KUBECONFIG
	// environment variable and the kubeConfig path. It also allows us to get the current selected namespace
	// from the kube configuration. With BuildConfigFromFlags this is not possible.
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// set namespace via namespace from kubeConfig if empty
	if namespace == "" {
		var err error
		namespace, _, err = kubeConfig.Namespace()
		if err != nil {
			return nil, err
		}
	}

	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	secretsClient := clientSet.CoreV1().Secrets(namespace)
	return &Client{
		secretsClient: secretsClient,
		dryRun:        dryRun,
		namespace:     namespace,
		secretName:    secretName,
		format:        format,
		data:          data,
		key:           key,
	}, nil
}

// Create either creates a new secret on the cluster or just prints the configuration in YAML or JSON
// if dry-run has been enabled
func (c *Client) Create() error {
	secret := coreV1.Secret{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      c.secretName,
			Namespace: c.namespace,
		},
		Data: map[string][]byte{c.key: c.data},
		Type: "Opaque",
	}
	if !c.dryRun {
		_, err := c.secretsClient.Create(context.Background(), &secret, metaV1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	switch c.format {
	case "json":
		result, err := json.MarshalIndent(secret, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(result))
		return nil
	case "yaml":
		result, err := yaml.Marshal(secret)
		if err != nil {
			return err
		}
		fmt.Println(string(result))
		return nil
	default:
		return ErrNoFormatSpecified
	}
}
