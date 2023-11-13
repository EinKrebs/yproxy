package main

import (
	"io"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/yezzey-gp/yproxy/config"
	"github.com/yezzey-gp/yproxy/pkg/proc"
	"github.com/yezzey-gp/yproxy/pkg/ylogger"
)

var cfgPath string
var logLevel string
var decrypt bool

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := config.LoadInstanceConfig(cfgPath)
		if err != nil {
			return err
		}

		instanceCnf := config.InstanceConfig()

		con, err := net.Dial("unix", instanceCnf.SocketPath)

		if err != nil {
			return err
		}

		defer con.Close()
		msg := proc.ConstructMessage(args[0], decrypt)
		_, err = con.Write(msg)
		if err != nil {
			return err
		}

		ylogger.Zero.Debug().Bytes("msg", msg).Msg("constructed message")

		_, err = io.Copy(os.Stdout, con)
		if err != nil {
			return err
		}

		// fmt.Println("reply:", string(reply))
		//
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "/etc/yproxy/yproxy.yaml", "path to yproxy config file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "", "log level")
	rootCmd.PersistentFlags().BoolVarP(&decrypt, "decrypt", "d", false, "decrypt external object or not")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		ylogger.Zero.Fatal().Err(err).Msg("")
	}
}
