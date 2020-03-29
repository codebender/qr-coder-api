/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	qrcodev1 "github.com/codebender/qrcode-api/proto/codebender/qrcode/v1"
	"github.com/codebender/qrcode-api/service"
	"github.com/codebender/qrcode-api/transport"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/oklog/oklog/pkg/group"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qrcode-api",
	Short: "Starts the QR Code API Service",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func start() {
	var (
		grpcAddr = fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "service", "qrcode-api")
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	server := grpc.NewServer()

	service := service.New(logger)
	endpoints := transport.New(service, logger)
	qrCodeAPIServer := transport.NewGRPCServer(endpoints, logger)

	var g group.Group
	{
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			level.Error(logger).Log("listener", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}

		g.Add(func() error {
			logger.Log("listener", "gRPC", "addr", grpcAddr)

			qrcodev1.RegisterQRCodeAPIServer(server, qrCodeAPIServer)

			reflection.Register(server)
			return server.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	logger.Log("exit", g.Run())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qrcode-api.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".qrcode-api" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".qrcode-api")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
