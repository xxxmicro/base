package source

import(
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/encoder/yaml"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/xxxmicro/go-micro-apollo-plugin"
	"github.com/xxxmicro/base/log"
)

func NewSourceProvider(c *cli.Context) source.Source {
	address := c.String("apollo_address")
	if len(address) == 0 {
 		address = Env("APOLLO_ADDRESS", "")
	}

	if len(address) == 0 {
 		log.Fatal("need config address")
	 	return nil
	}
		
	namespace := c.String("namespace")
	if len(namespace) == 0 {
 		namespace = Env("APOLLO_NAMESPACE", "application")
	}

	appId := c.String("apollo_app_id")
	if len(appId) == 0 {
 		appId = Env("APOLLO_APPID", "xpay-api")
	}

	cluster := c.String("apollo_cluster")
	if len(cluster) == 0 {
 		cluster = Env("APOLLO_CLUSTER", "dev")
	}

	backupConfigPath := Env("BACKUP_CONFIG_PATH", "./")

	e := yaml.NewEncoder()
	return apollo.NewSource(
		apollo.WithAddress(address),
		apollo.WithNamespace(namespace),
		apollo.WithAppId(appId),
		apollo.WithCluster(cluster),
		apollo.WithBackupConfigPath(backupConfigPath),
		source.WithEncoder(e),
	)
}
