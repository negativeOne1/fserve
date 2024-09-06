package client

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/martin.kluth1/fserve/signature"
)

type getSignatureOptions struct {
	Algorithm string
	Date      string
	Expires   string
	Method    string
	Resource  string
}

var DefaultAlgorithm = "HMAC:SHA256"

func addSignatureFlags(cmd *cobra.Command, o *getSignatureOptions) {
	cmd.Flags().StringVar(&o.Algorithm, "algorithm", DefaultAlgorithm, "The signature algorithm.")
	cmd.Flags().StringVar(&o.Date, "date", "", "The date and time the request was signed.")
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
	fmt.Println(opt)

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
	println(hex.EncodeToString(s))
}
