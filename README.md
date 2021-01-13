# pjl

Pretty JSON logger reads JSON logs and formats them in a more readable way.

## Usage

```sh
$ cat test.log
2021/01/13 09:00:23 Non-json log
{"command":"test-command","file":"/go/pkg/mod/github.com/sirupsen/logrus@v1.4.2/entry.go:285","func":"github.com/sirupsen/logrus.(*Entry).Info","level":"info","msg":"Test message","name": "Test","time":"2021-01-13T09:00:23.41Z"}

$ cat test.log | pjl -excludeJsonFields command,file,func
2021/01/13 09:00:23 Non-json log
INFO[2021-01-13T09:00:23.410Z] Test message                                  name=Test
```

## Install

```sh
go install github.com/bancek/pjl/cmd/pjl
```

## License

MIT
