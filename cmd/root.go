package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/idna"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Error struct {
	Code  int `json:"code"`
	Message string `json:"message"`
	Status string `json:"status"`
}
type Result struct {
	Error *Error `json:"error,omitempty"`
	ShortLink string `json:"shortLink,omitempty"`
}

type SocialMetaTagInfo struct {
	SocialTitle string `json:"socialTitle,omitempty"`
	SocialDescription string `json:"socialDescription,omitempty"`
	SocialImageLink string `json:"socialImageLink,omitempty"`
}

type DynamicLinkInfo struct {
	DomainUriPrefix string `json:"domainUriPrefix,omitempty"`
	Link string `json:"link,omitempty"`
	SocialMetaTagInfo  *SocialMetaTagInfo `json:"socialMetaTagInfo,omitempty"`
}

var cfgFile string
var webApiKey string

var dynamicLinkInfo = DynamicLinkInfo{
	SocialMetaTagInfo: &SocialMetaTagInfo{},
}

var rootCmd = &cobra.Command{
	Use:   "firelink",
	Short: "Create shortlink by Firebase Dynamic Link",
	Run: handle,
}

func handle(cmd *cobra.Command, args []string) {

	if dynamicLinkInfo.Link == "" {
		clipboardContent, err := clipboard.ReadAll()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		dynamicLinkInfo.Link = strings.TrimSpace(clipboardContent)
	}

	if _, err := url.ParseRequestURI(dynamicLinkInfo.Link); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	dynamicLinkInfo.DomainUriPrefix = viper.GetString("domainUriPrefix")

	values := map[string]interface{}{
		"dynamicLinkInfo": dynamicLinkInfo,
		"suffix": map[string]string{
			"option": "SHORT",
		},
	}

	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post("https://firebasedynamiclinks.googleapis.com/v1/shortLinks?key=" + viper.GetString("webApiKey"), "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if result.Error != nil {
		fmt.Print(result.Error)
		os.Exit(1)
	}

	result.ShortLink, err = convertLinkToUnicode(result.ShortLink)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Print(result.ShortLink)
	os.Exit(0)
}

func convertLinkToUnicode(originLink string)(string, error){
	parsedShortLink, err := url.Parse(originLink)
	if err != nil {
		return "", err
	}
	p := idna.New(idna.ValidateForRegistration())
	unicodeHost, err := p.ToUnicode(parsedShortLink.Host)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(originLink, parsedShortLink.Host, unicodeHost), nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c","", "config file (default is $HOME/.firelink.yaml)")
	rootCmd.PersistentFlags().StringVarP(&webApiKey, "key", "k", "", "Web API Key")
	rootCmd.PersistentFlags().StringVarP(&dynamicLinkInfo.DomainUriPrefix, "prefix", "p", "", "domainUriPrefix")

	rootCmd.Flags().StringVarP(&dynamicLinkInfo.Link, "link", "l", "", "Long-link to be shorten, if this flag is empty, we will try to get from clipboard")
	rootCmd.Flags().StringVarP(&dynamicLinkInfo.SocialMetaTagInfo.SocialTitle, "title", "t", "", "Social title")
	rootCmd.Flags().StringVarP(&dynamicLinkInfo.SocialMetaTagInfo.SocialDescription, "description", "d", "", "Social description")
	rootCmd.Flags().StringVarP(&dynamicLinkInfo.SocialMetaTagInfo.SocialImageLink, "imageLink", "i", "", "Social image link")

	if err := viper.BindPFlag("webApiKey", rootCmd.PersistentFlags().Lookup("key")); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("domainUriPrefix", rootCmd.PersistentFlags().Lookup("prefix")); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
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
			fmt.Print(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".firelink" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".firelink")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if !viper.IsSet("webApiKey") || !viper.IsSet("domainUriPrefix") {
			fmt.Print(err)
			os.Exit(1)
		}

	}
}
