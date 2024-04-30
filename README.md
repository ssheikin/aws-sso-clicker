## aws-sso-clicker

### Background

We want to avoid leaving the terminal and opening yet another tab and having to click Next next next...

This is a different version of https://github.com/mziyabo/headless-sso when sso already logged-in 
within the default browser, which just clicks Next next next...

### Install

checkout and install
``` sh
 go install aws-sso-clicker.go
```

### Usage:

``` bash
aws sso login --sso-session ${SSO_SESSION_NAME} --no-browser | ~/go/bin/aws-sso-clicker
```


**Note:** Browser has to be started with `--remote-debugging-port=9222` 
Or if browser is started from intellij it already has it.
More details: https://github.com/go-rod/go-rod.github.io/blob/main/custom-launch.md#custom-browser-launch

### Manual testing

Redirection to file may be useful for manual testing. Later feed it as `/tmp/a` or (`/tmp/b` as more generalized) as 
Intellij -> Runtime Configuration -> Redirect input from 
`aws sso login --sso-session ${SSO_SESSION_NAME} --no-browser > /tmp/a | (while ! grep user_code /tmp/a ; do sleep 1; done && grep user_code /tmp/a > /tmp/b)`
