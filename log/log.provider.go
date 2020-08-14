package log

import(
	"github.com/micro/go-micro/v2/config"
)

/**
 * fx.Invoke(InitLog)
 */
func InitLog(config config.Config) {
	env := config.Get("log", "env").String("dev")

	Init(env)

	// go watchConfig()
}

func watchConfig() {

}
