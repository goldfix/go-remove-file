[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/goldfix/go-remove-file/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/goldfix/go-remove-file.svg?branch=master)](https://travis-ci.org/goldfix/go-remove-file)

# go-remove-file
Possible substitute at 'rm' command. Similar functions plus recovery and trash functions.

<pre><code>
Usage of grm:
  -d        move files into recycle
  -e        empty recycle folder
  -f        ignore nonexistent files, never prompt
  -i        prompt before every files
  -ls       list files into recycle
  -r        remove contents recursively
  -u        recover list files from recycle
  -v        explain what is being done
  -version  output version information and exit
  -x        force delete files
  -t 		destination folder to restore
</code></pre>

#### Thanks to:
* https://github.com/satori/go.uuid

# Installation

To install go-remove-file, you can download a [prebuilt binary](https://github.com/goldfix/go-remove-file/releases), or you can build it from source.

### Prebuilt binaries

All you need to install go-remove-file is one file, the binary itself.

Download the binary from the [releases](https://github.com/goldfix/go-remove-file/releases) page.

### Building from source

If your operating system does not have binary, but does run Go, you can build from source.

Make sure that you have Go version 1.6 or greater.

```sh
go get -u github.com/goldfix/go-remove-file/...
```

...
