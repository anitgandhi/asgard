# asgard

This is a proof-of-concept to show how to use [memguard](https://github.com/awnumar/memguard) in tandem with Go's AES crypto primitives, in order to protect the key in memory.

The code is originally from the Go standard library generic and amd64 AES implementations, but modified to do the following when the cipher object is created:
  1. Creates an immutable buffer from the given key, which wipes the original key location
  2. Creates the encryption and decryption key schedule buffers, which are subsequently populated during the existing expansion call
  3. Destroys the key buffer, since it's no longer needed after expansion
  4. Marks the encryption and decryption key schedule buffers as immutable
  5. Disables core dumps on Unix

To understand how `memguard` actually protects values in memory, please see its documentation and README.

Since the API largely remains the same, and the returned cipher object continues to meet the `cipher.Block` interface, this is a drop-in replacement for the existing `crypto/aes` package. 

The package name has been changed to `asgard` on purpose, since there are some additional APIs, to avoid confusion, and for fun. `asgard ~= aesguard = aes + memguard`.

## Tests

The unit tests from the standard library were copied here as well. The only modification needed was in `TestCipherEncrypt` and `TestCipherDecrypt` to use a local copy of the key to prevent wiping the test vectors struct.

## Usage

asgard is a `package aes` drop-in replacement, so you can just replace `crypto/aes` imports with `github.com/anitgandhi/asgard`, and `aes.NewCipher` function calls with `asgard.NewCipher`

Additionally, the returned cipher object has a method `Destroy()`, which will destroy the enc/dec schedule buffers. This method can't be reached since it's hidden behind the unexported concrete types. So, there's an exported function `asgard.DestroyCipher` you can call after you're completely done using the AES block.

See an example here: https://github.com/anitgandhi/fpe-fun/blob/278ce577c7587df60fa7575b36985b9c5261e8e4/cmd/fpe-asgard/main.go

## Notes 

The amd64 optimized implementation could be ported because `golang.org/x/sys/cpu` provides the necessary CPU feature detection for AES-NI and PCMUL.

Unfortunately, it's hard to provide optimized implementations for the other platforms until `golang.org/x/sys/cpu` provides the CPU feature flags for them.

The GCM key schedule/stream is also protected, since it uses the same encryption key schedule already generated.

## Disclaimers

Important: the original key slice (`key []byte`) that is passed comes from memory managed by the Go runtime. Even after it's wiped, there's always a chance there is a dangling copy of it somewhere in memory or swap files, since it exists prior to involving `memguard`.

Using `memguard` does offer some benefits for "ongoing" protection, but it's not a generic gaurantee, given the nature of Go's memory model. As soon as you have a `[]byte` going into or coming out of `memguard` that results in a copy of the underlying data, it's owned by the Go runtime, and is no longer subject to `memguard` guarantees. [This comment](https://github.com/hashicorp/vault/issues/540#issuecomment-350757998) explains the issue well.

You can _possibly_ (definitely not guaranteed) reduce risk of exposure by somehow ensuring the fixed-size key array remains on the calling function's stack, but even then, stacks can be moved around by the Go runtime transparently. See the example link above.

Of course, if you're reading your key from environment variables, config files, or something else, that's always another point of possible exposure.

tl;dr the key is still coming from a Go runtime managed `[]byte` slice, so that's always going to be a point of possible exposure.

## TODO

Add a `NewCipherWithLockedBuffer` function to take an existing `*memguard.LockedBuffer` that's separately populated with a key by some other means. 
