# s3website-openid-proxy

This service provides [OpenID](https://openid.net/) authenticated access to a static website hosted in an s3 bucket.  

It is designed to be a simple way to add authentication to [static websites](https://en.wikipedia.org/wiki/Static_web_page) stored in [AWS S3](https://aws.amazon.com/s3/).

This service uses [AWS API Gateway](https://aws.amazon.com/api-gateway/) HTTP APIs and is powered by [AWS Lambda](https://aws.amazon.com/lambda/) .

# Goals

1. Provide a simple authentication access to static websites hosted in s3.
2. Utilise AWS lambda and API Gateway to enable low cost hosting.
3. Take advantage of the rate limiting provided by AWS API Gateway to ensure access isn't possible using [brute force attacks](https://en.wikipedia.org/wiki/Brute-force_attack).
4. Use existing opensource libraries to provide secure access via cookies.
5. Support OpenID authentication of users accessing the site.

# Deployment

You will need the following tools.

* [AWS cli](https://aws.amazon.com/cli/) 
* [SAM cli](https://github.com/aws/aws-sam-cli)

Also an aws profile setup with your [aws credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html).

Create an OpenID application in a service such as [Okta](https://www.okta.com/).

Create an .envrc file using [direnv](https://direnv.net/).

```bash
#!/bin/bash

export AWS_PROFILE=wolfeidau
export AWS_DEFAULT_PROFILE=wolfeidau
export AWS_REGION=ap-southeast-2

# these are provided by your OpenID provider 
export CLIENT_ID=xxxxxxxxx
export CLIENT_SECRET=xxxxxxxxx
export ISSUER=https://dev-xxxxxx.okta.com

export HOSTED_ZONE_ID=XXXXXXXXXX

# results in $SUBDOMAIN_NAME.$HOSTED_ZONE_NAME or something.wolfe.id.au
export HOSTED_ZONE_NAME=wolfe.id.au
export SUBDOMAIN_NAME=something
```

Run make.

```
make
```

# TODO

* [ ] Add support for [PKCE](https://oauth.net/2/pkce/)
* [ ] Containerise this service to enable running in [AWS fargate](https://aws.amazon.com/fargate/) or possibly [kubernetes](https://kubernetes.io/). 

# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).
