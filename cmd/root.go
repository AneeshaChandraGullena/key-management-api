// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.

package cmd

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	version "github.ibm.com/Alchemy-Key-Protect/go-hello-service"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/certificate"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/proxy"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/utils/logging"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-config"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
	"github.ibm.com/Alchemy-Key-Protect/protobuffer-go-spec/greeter"
	"golang.org/x/net/context"
)

// mainSemver is set by build to denote the semver numbering
var mainSemver string

// mainCommit is set by the build to denote the commit SHA1 of the build
var mainCommit string

var cfgFile string

func isDeployed() bool {
	return runtime.GOOS == constants.LinuxRuntime
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "key-management-api",
	Short: "IBM Key Protect API service",
	Long:  `IBM Key Protect API service provides access to all the microservices`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		config := configuration.Get()
		logger := logging.GlobalLogger()

		validateVersion(config)

		mainLogger := log.NewContext(logger).With("component", "main")
		mainLogger.Log("semver", mainSemver,
			"commit", mainCommit)

		httpAddr := flag.String("http.addr", ":"+config.GetString("host.http.port"), "HTTP listen address")
		grpcAddr := flag.String("grpc.addr", ":"+config.GetString("host.grpc.port"), "gRPC (HTTP) listen address")
		flag.Parse()

		ctx := context.Background()
		errs := make(chan error, 2)
		// Transport HTTP
		go func() {
			httpLogger := log.NewContext(logger).With("transport", "https")

			mux := http.NewServeMux()

			commonHandlers := middleware.BaseMiddlewareChain()

			apiHandlers := middleware.StdMiddlewareChain(commonHandlers)
			provHandlers := middleware.ProvisionMiddlewareChain(commonHandlers)

			http.Handle("/", commonHandlers.Then(mux))
			http.Handle("/api/", apiHandlers.Then(mux))

			// Make Reverse proxy to forward admin api requests
			adminReverseProxy := &proxy.Prox{}
			adminReverseProxy = proxy.New("http://" + config.GetString(constants.AdminServiceIpv4Address) + ":" + config.GetString(constants.AdminServicePort))
			mux.Handle("/admin/v1/", commonHandlers.ThenFunc(adminReverseProxy.Handle))

			// Make bluemix proxy to forward service api requests
			provisionReverseProxy := &proxy.Prox{}
			provisionReverseProxy = proxy.New("http://" + config.GetString(constants.BluemixServiceIpv4Address) + ":" + config.GetString(constants.BluemixServicePort))
			mux.Handle("/service/v1/", provHandlers.ThenFunc(provisionReverseProxy.Handle))

			// Make bluemix proxy to forward service api requests
			lifecycleReverseProxy := &proxy.Prox{}
			lifecycleReverseProxy = proxy.New("http://" + config.GetString(constants.LifecycleServiceIpv4Address) + ":" + config.GetString(constants.LifecycleServicePort))
			mux.Handle("/api/v2/", apiHandlers.ThenFunc(lifecycleReverseProxy.Handle))

			errHTTPLogger := httpLogger.Log("address", *httpAddr, "msg", "listening")
			if errHTTPLogger != nil {
				panic("Unable to Log HTTP transport")
			}

			// for local testing
			if isDeployed() == false {
				config.Set(constants.HTTPSCertBasePath, ".")
			}
			basePath := config.GetString(constants.HTTPSCertBasePath)

			certPath := basePath + config.GetString(constants.HTTPSCertPath)
			keyPath := basePath + config.GetString(constants.HTTPSKeyPath)
			err := certificate.Exists(certPath, keyPath)
			// If they are not available, generate new ones.
			if err != nil {
				httpLogger.Log("ERROR", "No certs found, so generated some",
					"certPath", certPath)
				err = certificate.Generate(certPath, keyPath, *httpAddr)
				if err != nil {
					panic("Error: Couldn't create https certs.")
				}
			}

			tlsConfig := &tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: config.GetString(constants.TLSServerName),

				// TODO once at go v1.7 can  uncomment other good ciphers [elo 09/02/2016]
				// Note: The cipher list must be ordered in preferences from most desired to least
				// Note: 128 bit ciphers are considered less secure than 256 bit ciphers so they have been disabled
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					//tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384,
					//tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
					//tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					//tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
					//tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					//tls.TLS_RSA_WITH_AES_256_CBC_SHA256,
					//tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
				},
				PreferServerCipherSuites: true,
			}
			readTimeout := config.GetInt("timeouts.readTimeout")
			writeTimeout := config.GetInt("timeouts.writeTimeout")
			// Header size is ~ 2x the typical header items we have: jwt, org, space, instance, content
			server := &http.Server{
				Addr:           *httpAddr,
				ReadTimeout:    time.Duration(readTimeout) * time.Second,
				WriteTimeout:   time.Duration(writeTimeout) * time.Second,
				MaxHeaderBytes: 1 << 11,
				TLSConfig:      tlsConfig,
				TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
			}

			errs <- server.ListenAndServeTLS(certPath, keyPath)
		}()

		// Transport gRPC
		go func() {
			grpcLogger := log.NewContext(logger).With("transport", "gRPC")

			ln, err := net.Listen("tcp", *grpcAddr)
			if err != nil {
				errs <- err
				return
			}

			s := grpc.NewServer()

			var versionEndpoint endpoint.Endpoint
			{
				versionLogger := log.NewContext(grpcLogger).With("method", "Version")
				versionEndpoint = version.MakeVersionEndpoint(version.NewBasicService())
				versionEndpoint = version.EndpointLoggingMiddleware(versionLogger)(versionEndpoint)
			}

			var runtimeVersionEndpoint endpoint.Endpoint
			{
				runtimeVersionLogger := log.NewContext(grpcLogger).With("method", "RuntimeVersion")
				runtimeVersionEndpoint = version.MakeRuntimeVersionEndpoint(version.NewBasicService())
				runtimeVersionEndpoint = version.EndpointLoggingMiddleware(runtimeVersionLogger)(runtimeVersionEndpoint)
			}

			endpointsGreeter := version.Endpoints{
				VersionEndpoint:        versionEndpoint,
				RuntimeVersionEndpoint: runtimeVersionEndpoint,
			}

			svcVersion := version.MakeGRPCServer(ctx, endpointsGreeter, logger)
			greeter.RegisterGreeterServer(s, svcVersion)

			errGRPCLogger := grpcLogger.Log("address", *grpcAddr, "msg", "listening")
			if errGRPCLogger != nil {
				panic("Unable to Log HTTP transport")
			}
			errs <- s.Serve(ln)
		}()

		go func() {
			signalChan := make(chan os.Signal)
			signal.Notify(signalChan, syscall.SIGINT)
			errs <- fmt.Errorf("%s", <-signalChan)
		}()

		errSigLog := logger.Log("terminated", <-errs)
		if errSigLog != nil {
			panic("cannot log basic server info")
		}
	},
}

func validateVersion(config configuration.Configuration) {
	// ensure that the configuration file and binary file were built together
	configVersion := config.GetString("version.semver")
	if isDeployed() && mainSemver != configVersion {
		panic(fmt.Sprintf("Version mismatch enabled on %s: expected %s have %s ", runtime.GOOS, configVersion, mainSemver))
	}

	configCommit := config.GetString("version.commit")
	if isDeployed() && mainCommit != mainCommit {
		panic(fmt.Sprintf("Commit mismatch enabled on %s: expected %s have %s ", runtime.GOOS, configCommit, mainCommit))
	}
}

// SetVersion needs to be called by main.main() to set build version, so that the version commmand returns the value matching the build
func SetVersion(version string, commit string) {
	if version == "" {
		mainSemver = "0.0.0"
	} else {
		mainSemver = version
	}
	if commit == "" {
		mainCommit = "0000"
	} else {
		mainCommit = commit
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.key-management-api.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(constants.APIServerConfigFile) // name of config file (without extension)
	viper.AddConfigPath(constants.ServerConfigPath)    // adding home directory as first search path
	viper.AutomaticEnv()                               // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
