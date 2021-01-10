# aws-openid-proxy

This service provides a proxy to backend services/content with [OpenID](https://openid.net/) used for authentication, it is designed to be a simple way to add authentication to [static websites](https://en.wikipedia.org/wiki/Static_web_page) stored in [AWS s3](https://aws.amazon.com/s3/) or containerised applications running in [AWS fargate](https://aws.amazon.com/fargate/) or possibly kubernetes.

This service uses [AWS API Gateway](https://aws.amazon.com/api-gateway/) HTTP APIs with an OpenID authoriser to ensure users are authenticated before accessing content.

# Goals

1. Provide a secure audited access layer to static websites hosted in s3.
2. Utilise [AWS lambda](https://aws.amazon.com/lambda/) to enable low to no cost of hosting.
3. Take advantage of the rate limiting provided by AWS API Gateway to ensure access isn't possible using [brute force attacks](https://en.wikipedia.org/wiki/Brute-force_attack).
4. Use the authoriser lambdas to ensure there is a clear separation between access and authentication logic.
5. Use existing opensource libraries to provide secure access via cookies.

# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).