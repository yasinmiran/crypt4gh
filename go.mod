module github.com/neicnordic/crypt4gh

go 1.20

require (
	filippo.io/edwards25519 v1.1.0
	github.com/dchest/bcrypt_pbkdf v0.0.0-20150205184540-83f37f9c154a
	github.com/hashicorp/go-version v1.6.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/logrusorgru/aurora/v4 v4.0.0
	golang.org/x/crypto v0.18.0
	golang.org/x/term v0.16.0
)

require golang.org/x/sys v0.16.0 // indirect

retract v1.8.7 // has a bug related to file decryption that ends up in loop.
