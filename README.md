# engine

Based off [gin](github.com/gin-gonic/gin), **engine** is a collection of helpers to kick start for web projects in go.
Not sure how fast it is - just works for me.

## Installation

1. Get it.

```sh
$ go get github.com/mnbbrown/engine
```

2. Import and use.

```go
package main

import (
    "fmt"
    "github.com/mnbbrown/engine"
    "net/http"
)

func main() {
    e := engine.NewRouter()
    e.Get("/", func(rw *http.ResponseWriter, req *http.Request) {
        fmt.Fprint(rw, "Hello World!")
    })
    http.ListenAndServe(":8080", e)
}

```

## Licence

The MIT License (MIT)

Copyright (c) 2015- Matthew Brown <me@matthewbrown.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
