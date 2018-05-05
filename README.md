# aesguard

This is a proof-of-concept to show how to use [memguard](https://github.com/awnumar/memguard) in tandem with Go's AES crypto primitives, in order to protect the key in memory.

The code is originally from the Go standard library generic and amd64 AES implementations, but modified to do the following when the cipher object is created:
  1. Creates an immutable buffer from the given key, which wipes the original key location
  2. Creates the encryption and decryption key schedule buffers, which are subsequently populated during the existing expansion call
  3. Destroys the key buffer, since it's no longer needed after expansion
  4. Marks the encryption and decryption key schedule buffers as immutable
  5. Disables core dumps on Unix

To understand how `memguard` actually protects values in memory, please see its documentation and README.

Since the API remains the same, and everything is based around the `cipher.Block` interface, this is largely a drop-in replacement for the existing `crypto/aes` package

## Tests

The unit tests from the standard library were copied here as well. The only modification needed was in `TestCipherEncrypt` and `TestCipherDecrypt` to use a local copy of the key to prevent wiping the test vectors struct.

## Notes

The amd64 optimized implementation could be ported because `golang.org/x/sys/cpu` provides the necessary CPU feature detection for AES-NI and PCMUL.

Unfortunately, it's hard to provide optimized implementations for the other platforms until `golang.org/x/sys/cpu` provides the CPU feature flags for them.


## TODO

CTR/GCM XORKeyStream protection?
