package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
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
	logFilename := flag.String("log", "/tmp/media-notifier.log", "File to write logs")
	version := flag.Bool("version", false, "Print version")
	flag.Parse()

	if *version {
		fmt.Println(versionString())
		os.Exit(0)
	}

	if *topicArn == "" {
		fmt.Fprintln(os.Stderr, "SNS topic ARN is required\n")
		os.Exit(1)
	}

	logFile, err := os.OpenFile(*logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer logFile.Close()

	logger := log.New(logFile, versionString()+": ", log.LstdFlags)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           *credProfile,
		SharedConfigState: session.SharedConfigEnable,
	}))

	message, err := buildMessage()
	if err != nil {
		logger.Fatal(err)
	}

	client := sns.New(sess)
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(*topicArn),
	}

	// First return value is the MessageId returned from SNS on success.
	_, err = client.Publish(input)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Sent SMS: %s\n", message)
}

func buildMessage() (string, error) {
	name, ok := os.LookupEnv("TR_TORRENT_NAME")
	if !ok {
		return "", errors.New("TR_TORRENT_NAME not set")
	}
	return "Downloaded: " + name, nil
}

func versionString() string {
	return fmt.Sprintf("v%s-%s", Version, BuildTag)
}
