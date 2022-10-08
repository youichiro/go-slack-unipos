package util

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
)

func VerifySlackSigningSecret(c *gin.Context, signingSecret string) error {
	verifier, err := slack.NewSecretsVerifier(c.Request.Header, signingSecret)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	verifier.Write(body)
	if err = verifier.Ensure(); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
