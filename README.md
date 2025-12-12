# aws-schedule-expressions

Validator for AWS EventBridge Scheduler schedule expressions

## Features

- Validate `at(...)`, `cron(...)` and `rate(...)` schedule expressions
- No third-party dependencies

## Usage

```go
package main

import (
	"fmt"
  "os"

	"github.com/robdasilva/aws-schedule-expressions/validator"
)

func main() {
	if ok := validator.Validate("at(2025-12-12T09:57:32)"); !ok {
    fmt.Println("The given `at` expression is invalid.")
    os.Exit(1)
  }

  if ok := validator.Validate("cron(57 9 ? */3 FRI#2 *)"); !ok {
    fmt.Println("The given `cron` expression is invalid.")
    os.Exit(1)
  }

  if ok := validator.Validate("rate(90 days)"); !ok {
    fmt.Println("The given `rate` expression is invalid.")
    os.Exit(1)
  }
}
```
