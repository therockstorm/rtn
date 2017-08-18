package main

import (
	"encoding/json"
	"log"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/therockstorm/rtn"
)

// Handle Lambda invocations via scheduled job.
func Handle(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	err := rtn.NewUpdater().Update()
	if err != nil {
		log.Println(err.Error())
	}

	return nil, err
}
