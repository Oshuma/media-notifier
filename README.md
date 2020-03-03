# media-notifier
Simple program that's triggered from Transmission to send an [Amazon SNS](https://aws.amazon.com/sns/) message on torrent complete.

Due to the way Transmission triggers 'done' scripts, this must be wrapped in a shell script to work properly.

`notifier.sh`:
```bash
#!/usr/bin/env bash

/path/to/media-notifier -topic "arn:aws:sns:us-east-1:xxxxxxxxxxxx:Example"
```

Make sure that script is executable, then update Transmission's `settings.json` to use the script:

```json
"script-torrent-done-enabled": true,
"script-torrent-done-filename": "/path/to/notifier.sh",
```

### Usage

```shell
$ media-notifier -help
Usage of media-notifier:
  -credentials string
        AWS credentials profile; see https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html (default "default")
  -log string
        File to write logs (default "/tmp/media-notifier.log")
  -topic string
        SNS topic ARN (required)
  -version
        Print version
```
