##

This projects aims to provide a simple and secure way to share files between producers (uploaders) and consumers (downloaders) using a simple REST API.
The links are secured in AWSs [S3 pre-signed URLs](https://docs.aws.amazon.com/AmazonS3/latest/userguide/ShareObjectPreSignedURL.html) style.
But can be easily stored in various storage backends (e.g. local, S3, GCS, Azure, etc). Typical usecases include onprem environments where you want to share files between different services or users but can't use public cloud services and SaaS solutions where you want to have more control over the data.

# Features

- Secure file sharing
- Simple REST API
- Easy to deploy
- Easy to use
- Supports multiple storage backends
- Supports multiple file types
- Supports file expiration
- Supports file size limits

a typical url looks like this:

```
https://files.acme.ai/f2d7edc7-aabd-4a07-9305-b6bc81925f86.jpg
?Fs-Algorithm=HMAC:SHA256
&Fs-Date=20240829T090533Z
&Fs-Expires=600  //10min
&Fs-Signature:<Method;Date;Expires;Resource>
```

# TODO

- [ ] Add TTL support for PUT calls
- [ ] k8s
- [ ] Add more tests
- [ ] add AWS S3 support
- [ ] add usage examples
- [ ] python examples
