package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	host       string
	port       int
	retries    int
	retryDelay int
)

var (
	rootCmd = &cobra.Command{
		Use:   "ngrokcd",
		Short: "Attaches multiple subdomains to tunnels",
		Long: `This requires modification to ngrok.yml configuration file.

The following is the format for ngrok.yml file:

authtoken: [your ngrok authtoken]
tunnels:
  plex:
    proto: http
    addr: 32400
    bind_tls: false
    records:
	- movies
	- media
records:
  movies:
    dns: personal-domain
    cname: movies
  media:
    dns: personal-domain
    cname: media
dns:
  personal-domain:
    service: godaddy
    domain: aliyousuf.com
    key: [your godaddy key]
    secret: [your godaddy secret]
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ngrok := getNgrok(host, port)

			configPath := cmd.Flag("config").Value.String()
			config, err := loadConfiguration(configPath)
			if err != nil {
				return err
			}

			return performMatchingWithRetries(ngrok, config, retries, retryDelay)
		},
	}
	attachFirstCmd = &cobra.Command{
		Use:          "attach-first [subdomain] [domain] [service=godaddy] [key] [secret]",
		Short:        "Attaches a subdomain to the first ngrok tunnel it finds",
		Args:         cobra.ExactArgs(5),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ngrok := getNgrok(host, port)
			tunnels, err := ngrok.Tunnels()
			if err != nil {
				return err
			}

			tunnel := tunnels.First()
			config := newConfiguration(tunnel.Name, args[0], args[1], args[2], args[3], args[4])

			return performMatchingWithRetries(ngrok, config, retries, retryDelay)
		},
	}
	attachSpecificCmd = &cobra.Command{
		Use:          "attach [ngrok tunnel name] [subdomain] [domain] [service=godaddy] [key] [secret]",
		Short:        "Attaches a subdomain to a specific ngrok tunnel",
		Args:         cobra.ExactArgs(6),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ngrok := getNgrok(host, port)
			config := newConfiguration(args[0], args[1], args[2], args[3], args[4], args[5])

			return performMatchingWithRetries(ngrok, config, retries, retryDelay)
		},
	}
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	defaultConfig := fmt.Sprintf("%s/.ngrok2/ngrok.yml", home)

	rootCmd.AddCommand(attachFirstCmd)
	rootCmd.AddCommand(attachSpecificCmd)
	rootCmd.Flags().StringP("config", "c", defaultConfig, "ngrok YAML config file")

	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "localhost", "ngrok hostname")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 4040, "ngrok port")
	rootCmd.PersistentFlags().IntVarP(&retries, "retries", "r", 5, "retry count")
	rootCmd.PersistentFlags().IntVarP(&retryDelay, "delay", "d", 1000, "retry delay (ms)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
