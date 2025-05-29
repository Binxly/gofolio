//go:build !lambda
// +build !lambda

package main

import (
	"log"

	"github.com/gorilla/mux"
)

func startLambda(r *mux.Router) {
	log.Fatal("Required for compiler skill issues. This means a Lambda function called in non-lambda build - if you're seeing this, you probably fucked up")
}
