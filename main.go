package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	Version = "0.0.1"

	// BuildTag is updated from the current git SHA when the binary is built.
	BuildTag = "HEAD"
)

func main() {
	credProfile := flag.String("credentials", "default", "AWS credentials profile; see https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html")
	topicArn := flag.String("topic", "", "SNS topic ARN (required)")
	version := flag.Bool("version", false, "Print version")
	flag.Parse()

	if *version {
		fmt.Printf("v%s-%s\n", Version, BuildTag)
		os.Exit(0)
	}

	if *topicArn == "" {
		fmt.Fprintln(os.Stderr, "SNS topic ARN is required\n")
		os.Exit(1)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           *credProfile,
		SharedConfigState: session.SharedConfigEnable,
	}))

	message, err := buildMessage()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client := sns.New(sess)
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(*topicArn),
	}

	// First return value is the MessageId returned from SNS on success.
	_, err = client.Publish(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Sent SMS: %s\n", message)
}

func buildMessage() (string, error) {
	name, ok := os.LookupEnv("TR_TORRENT_NAME")
	if !ok {
		return "", errors.New("TR_TORRENT_NAME not set")
	}
	return "Downloaded: " + name, nil
}
