package client

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/martin.kluth1/fserve/signature"
)

type getSignatureOptions struct {
	Algorithm string
	Date      string
	Expires   string
	Method    string
	Resource  string
	Host      string
	Port      string
}

var DefaultAlgorithm = "HMAC-SHA256"

func addSignatureFlags(cmd *cobra.Command, o *getSignatureOptions) {
	cmd.Flags().StringVar(&o.Host, "host", "localhost", "The host to connect to.")
	cmd.Flags().StringVar(&o.Port, "port", "8080", "The port to connect to.")
	cmd.Flags().StringVar(&o.Algorithm, "algorithm", DefaultAlgorithm, "The signature algorithm.")
	cmd.Flags().
		StringVar(&o.Date, "date", time.Now().UTC().Format(time.RFC3339),
			"The date and time the request was signed.")
	cmd.Flags().StringVar(&o.Expires, "expires", "", "The date and time the request expires.")
	cmd.Flags().StringVar(&o.Method, "method", "", "The HTTP method.")
	cmd.Flags().StringVar(&o.Resource, "resource", "", "The resource being requested.")
}

var createSignatureCmd = &cobra.Command{
	Use:   "create-signature",
	Short: "create-signature",
	Run:   createSignature,
}

func createSignature(cmd *cobra.Command, args []string) {
	s, err := signature.CreateSignature(
		opt.Algorithm,
		opt.Date,
		opt.Expires,
		opt.Method,
		opt.Resource,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	params := url.Values{
		"Fs-Algorithm": []string{opt.Algorithm},
		"Fs-Date":      []string{opt.Date},
		"Fs-Expires":   []string{opt.Expires},
		"Fs-Signature": []string{hex.EncodeToString(s)},
	}

	// q := fmt.Sprintf(
	// 	"%s?Fs-Algorithm=%s\\&Fs-Date=%s\\&Fs-Expires=%s&Fs-Signature=%s",
	// 	opt.Resource,
	// 	opt.Algorithm,
	// 	opt.Date,
	// 	opt.Expires,
	// 	hex.EncodeToString(s),
	// )

	fmt.Printf("http://%s:%s/v1/%s?%s\n", opt.Host, opt.Port, opt.Resource, params.Encode())
}
