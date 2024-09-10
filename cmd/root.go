/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "img",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: img,
}

var (
	query string
	ext   string
	size  int
)

type Img struct {
	Name string `json:"name"`
}

type ImageListResponse struct {
	Icons []string `json:"icons"`
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.img.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&query, "query", "q", "", "image query you want to search (required)")
	rootCmd.MarkFlagRequired("name")

	rootCmd.Flags().StringVarP(&ext, "extension", "e", "", "extension of the image you want to search (required)")
	rootCmd.MarkFlagRequired("ext")

	rootCmd.Flags().IntVarP(&size, "size", "s", 0, "size of the image you want to search (optional)")
}

func img(cmd *cobra.Command, args []string) error {

	apiURL := "https://api.iconify.design"

	apiSearchURL, err := url.Parse(apiURL + "/search")
	if err != nil {
		return fmt.Errorf("can't parse url: %w", err)
	}

	params := url.Values{}
	params.Add("query", strings.ToLower(query))
	apiSearchURL.RawQuery = params.Encode()


	imgListResp, err := http.Get(apiSearchURL.String())
	if err != nil {
		return fmt.Errorf("can't get image list: %w", err)
	}

	defer imgListResp.Body.Close()

	if imgListResp.StatusCode != http.StatusOK {
		return fmt.Errorf("can't get image list: server error.  %w", err)
	}

	imgListBody, err := io.ReadAll(imgListResp.Body)
	if err != nil {
		return fmt.Errorf("can't read image list: %w", err)
	}

	var imgList ImageListResponse
	err = json.Unmarshal(imgListBody, &imgList)
	if err != nil {
		return fmt.Errorf("can't unmarshal image list: %w", err)
	}

	//	fmt.Println(imgList.Icons)

	icon := ""

	// check if exists a file with the name icon in the current working directory
	for _, img := range imgList.Icons {
		tmp := img + "." + ext
		if _, err := os.Stat(tmp); err == nil {
			continue
		}

		icon = img
		break
	}

	// split into prefix and name
	prefix := strings.Split(icon, ":")[0]
	name := strings.Split(icon, ":")[1]
	fmt.Println(prefix, name)

	retrieveIconURL := apiURL + "/" + prefix + "/" + name + ".svg"

	iconResp, err := http.Get(retrieveIconURL)
	if err != nil {
		return fmt.Errorf("can't get image: %w", err)
	}

	defer iconResp.Body.Close()

	if iconResp.StatusCode != http.StatusOK {
		return fmt.Errorf("can't get image: server error.  %w", err)
	}

	iconBody, err := io.ReadAll(iconResp.Body)
	if err != nil {
		return fmt.Errorf("can't read image: %w", err)
	}

	fmt.Println(string(iconBody))

	// create a file with the name icon in the current working directory

	


	//	fmt.Println("Image downloaded successfully")

	/*
	   pwd, err := os.Getwd()

	   	if err != nil {
	   		return fmt.Errorf("can't get current working directory: %w", err)
	   	}
	*/
	return nil
}
