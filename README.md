# goMicroserviceTemplate

A go microservice that uses gin and upstart with two deploy options: directly to target VM via ssh (for small deployments), or to S3 then down to VMs via pull an automated service (for large deployments).

## Deployment

#### Developer machine

##### Setup

Install pip3

`sudo apt-get update`

`sudo apt-get install -y python3-pip`

Install awscli

`sudo pip3 install awscli --force-reinstall --upgrade`

##### Single-VM deployment

(VM domains/IP address and ssh key locations must first be configured in simpleDeploy.sh for all but local deployments.)

Running `bash simpleDeploy.sh [env]` will run `go install *.go`, gzip the resulting binary, and transfer the resulting file to a single vm, and then restart the goMicroserviceTemplate service.

##### Multi-VM deployment

(S3 bucket names and prefixes must first be configured in s3Deploy.sh and updateGoMicroserviceTemplate.conf)

Running `bash s3Deploy.sh [env]` will run `go install *.go`, gzip the resulting binary, and then upload the resulting file to S3 (with bucket and object prefix determined by the value of [env]). [env] can be `local`, `dev`, or `prod`.

#### Destination machine

##### Setup

Install pip3

`sudo apt-get update`

`sudo apt-get install -y python3-pip`

Install awscli

`sudo pip3 install awscli --force-reinstall --upgrade`

##### Operation

If the `updateGoMicroserviceTemplate.conf` file is placed in /etc/init/, then every 10 seconds the MD5 hash of the `/usr/local/bin/goMicroserviceTemplate.gz` file will be compared with the MD5 hash of the source object in S3 (where the bucket name and object prefix is determined by value of ENVTYPE in /etc/environment). If the hashes are different, the new binary is downloaded, uncompressed, and the `goMicroserviceTemplate` service is restarted.

The destination machine must have read access to the S3 bucket and object prefix where the source binary will be located.

## Endpoints

##### GET /status

Returns
```js
{
    "status": "ok"
}
```
