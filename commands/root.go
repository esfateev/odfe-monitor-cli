package commands

import (
	"github.com/json-iterator/go"
	"net/http"
	"os"
	"strings"

	"github.com/esfateev/odfe-monitor-cli/es"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//Verbose logging if it is true, default to false
var Verbose bool

// ESConfig holds the for ES configuration
var esClient es.Client

var esURL string
var userName string
var password string
var rootDir string
var odVersion int

// RootCmd asd
var rootCmd = &cobra.Command{
	Use:   "odfe-monitor-cli",
	Short: "Manage opendistro alerting monitors.",
	Long:  `This application will help you to manage the Opendistro alerting monitors using YAML files.`,
}

func init() {
	cobra.OnInitialize(setup)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get CWD", err)
	}
	rootCmd.PersistentFlags().StringVarP(&rootDir, "rootDir", "r", dir, "Root directory where monitors yml files")
	rootCmd.PersistentFlags().StringVarP(&esURL, "esUrl", "e", "https://localhost:9200/", "URL to connect to Elasticsearch")
	rootCmd.PersistentFlags().StringVarP(&userName, "username", "u", "admin", "Username for opendistro Elasticsearch")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "admin", "Password for opendistro Elasticsearch")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().IntVarP(&odVersion, "odVersion", "", 0, "Major opendistro version")
}

func setup() {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if esURL != "" {
		//Validate URL
		if isURL(esURL) {
			// Validate ES is running?
			trailing := strings.HasSuffix(esURL, "/")
			if trailing {
				esURL = strings.TrimSuffix(esURL, "/")
			}
			esClient = es.Client{URL: esURL, Username: userName, Password: password, OdVersion: odVersion}
			resp, err := esClient.MakeRequest(http.MethodGet, "", nil, nil)
			indentJSON, _ := json.MarshalIndent(resp, "", "\t")
			check(err)
			if resp.Status != 200 {
				log.Fatal("Unable to connect to elasticsearch \n", string(indentJSON))
			}
		} else {
			log.WithFields(log.Fields{"elasticsearch-url": esURL}).Fatal("Elasticsearch url is invalid")
		}
	} else {
		// Solve with required flags
		log.Fatal("Ensure esURL is provided")
	}

	if Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

//Execute initiate the program and let cobra handles the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
