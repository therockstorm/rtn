# rtn

Download, parse, and store the [Fed routing directory](https://www.frbservices.org/EPaymentsDirectory/FedACHdir.txt) to S3.

## Install dependencies

`brew update && brew install golang && npm install -g serverless`  

Create an S3 bucket matching Updater.bucket

## Run tests

`go test -cover`  

## Build for production

This project uses [aws-lambda-go-shim](https://github.com/eawsy/aws-lambda-go-shim) which relies on Go plugins that are only supported on Linux. You can develop and run tests on a Mac, but building `cmd/handler.go` must happen in Linux. For [Vagrant](https://www.vagrantup.com/), add this to your Vagrantfile and reprovision, `config.vm.synced_folder "#{ENV['HOME']}/go", "/vagrant/go", create: true`. Once vagrant is up, run `vagrant ssh` and then,

`cd /vagrant/go/src/github.com/therockstorm/rtn && make`  

## Deploy

For Linux or Mac, make sure your AWS credentials are set locally and run,

`make deploy`  
