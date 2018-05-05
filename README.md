# aesguard

This is a proof-of-concept to show how to use [memguard](https://github.com/awnumar/memguard) in tandem with Go's AES crypto primitives, in order to protect the key in memory.

The code is originally from the Go standard library generic and amd64 AES implementations, but modified to do the following when the cipher object is created:
  1. Creates an immutable buffer from the given key, which wipes the original key location
  2. Creates the encryption and decryption key schedule buffers, which are subsequently populated during the existing expansion call
  3. Destroys the key buffer, since it's no longer needed after expansion
  4. Marks the encryption and decryption key schedule buffers as immutable
  5. Disables core dumps on Unix

To understand how `memguard` actually protects values in memory, please see its documentation and README.

Since the API largely remains the same, and the returned cipher object continues to meet the `cipher.Block` interface, this is a drop-in replacement for the existing `crypto/aes` package. 

The package name has been changed to `aesguard` on purpose, since there are some additional APIs, and to avoid confusion.

## Tests

The unit tests from the standard library were copied here as well. The only modification needed was in `TestCipherEncrypt` and `TestCipherDecrypt` to use a local copy of the key to prevent wiping the test vectors struct.

## Usage

aesguard is a `package aes` drop-in replacement, so you can just replace `crypto/aes` imports with `github.com/anitgandhi/aesguard`, and `aes.NewCipher` function calls with `aesguard.Nipher`

Additionally, the returned cipher object has a method `Destroy()`, which will destroy the enc/dec schedule buffers. This method can't be reached since it's hidden behind the unexported concrete types. So, there's an exported function `aesguard.DestroyCipher` you can call after you're completely done using the AES block.

## Notes

The amd64 optimized implementation could be ported because `golang.org/x/sys/cpu` provides the necessary CPU feature detection for AES-NI and PCMUL.

Unfortunately, it's hard to provide optimized implementations for the other platforms until `golang.org/x/sys/cpu` provides the CPU feature flags for them.

## TODO

CTR/GCM XORKeyStream protection?
