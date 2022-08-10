package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const LND_NODE = "lnd"

type Credentials struct {
	MacPath string `yaml:"macPath"`
	TlsPath string `yaml:"tlsPath"`
}

type Configuration struct {
	Name        string      `yaml:"name"`
	NodeType    string      `yaml:"nodeType"`
	Address     string      `yaml:"address"`
	Network     string      `yaml:"network"`
	Credentials Credentials `yaml:"credentials"`
}

type Config struct {
	Active  string                   `yaml:"active"`
	Configs map[string]Configuration `yaml:"configs"`
}

func getHomeDir() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return home
}

func getFuseDir() string {
	home := getHomeDir()
	dir := filepath.Join(home, ".fuse")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		cobra.CheckErr(err)
	}

	return dir
}

func generateConfig() {

	fuseDir := getFuseDir()
	credDir := filepath.Join(fuseDir, "credentials")

	credentials := Credentials{
		MacPath: filepath.Join(credDir, "admin.macaroon"),
		TlsPath: filepath.Join(credDir, "tls.cert"),
	}

	configuration := Configuration{
		Name:        "fuse-ln",
		NodeType:    LND_NODE,
		Address:     "localhost:1000",
		Network:     "regtest",
		Credentials: credentials,
	}

	configs := map[string]Configuration{}
	configs[configuration.Name] = configuration

	config := Config{
		Active:  configuration.Name,
		Configs: configs,
	}

	data, err := yaml.Marshal(config)
	cobra.CheckErr(err)

	err = viper.MergeConfig(bytes.NewBuffer(data))
	cobra.CheckErr(err)

	viper.SafeWriteConfig()
}

func initConfig() {

	viper.AddConfigPath(getFuseDir())
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.SetEnvPrefix("FUSE_")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	// TODO: Check specifically for no config error found
	if err != nil {
		generateConfig()
	}
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage config file for fusecli",
}

var getConfigCmd = &cobra.Command{
	Use:   "get {config_name}",
	Short: "Get config information",
	Long:  `Reads config file and retrieves information about the specified config`,
	Run: func(cmd *cobra.Command, args []string) {
		var config Config
		err := viper.Unmarshal(&config)
		cobra.CheckErr(err)

		var configName string
		if len(args) == 0 {
			configName = config.Active
		} else {
			configName = args[0]
		}

		configuration := config.Configs[configName]

		s, _ := json.MarshalIndent(configuration, "", "\t")
		fmt.Println(string(s))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(getConfigCmd)
}
