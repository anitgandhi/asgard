// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package aesguard

import (
	"crypto/cipher"
	"errors"
	"strconv"

	"github.com/awnumar/memguard"
)

// BlockSize is the AES block size in bytes.
const BlockSize = 16

// A cipher is an instance of AES encryption using a particular key.
type aesCipher struct {
	enc             []uint32
	dec             []uint32
	encLockedBuffer *memguard.LockedBuffer
	decLockedBuffer *memguard.LockedBuffer
}

// KeySizeError is an error message for invalid AES key sizes
type KeySizeError int

func (k KeySizeError) Error() string {
	return "crypto/aes: invalid key size " + strconv.Itoa(int(k))
}

// NewCipher creates and returns a new cipher.Block.
// The key argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func NewCipher(key []byte) (cipher.Block, error) {
	k := len(key)
	switch k {
	default:
		return nil, KeySizeError(k)
	case 16, 24, 32:
		break
	}
	return newCipher(key)
}

// newCipherGeneric creates and returns a new cipher.Block
// implemented in pure Go.
func newCipherGeneric(key []byte) (cipher.Block, error) {
	keyBuffer, encScheduleBuffer, decScheduleBuffer, encScheduleUint32, decScheduleUint32, err := prepareMemguard(key)
	if err != nil {
		return nil, err
	}

	c := aesCipher{
		enc:             encScheduleUint32,
		dec:             decScheduleUint32,
		encLockedBuffer: encScheduleBuffer,
		decLockedBuffer: decScheduleBuffer,
	}
	expandKeyGo(keyBuffer.Buffer(), c.enc, c.dec)

	err = finalizeMemguard(keyBuffer, encScheduleBuffer, decScheduleBuffer)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *aesCipher) BlockSize() int { return BlockSize }

func (c *aesCipher) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("crypto/aes: input not full block")
	}
	if len(dst) < BlockSize {
		panic("crypto/aes: output not full block")
	}
	encryptBlockGo(c.enc, dst, src)
}

func (c *aesCipher) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("crypto/aes: input not full block")
	}
	if len(dst) < BlockSize {
		panic("crypto/aes: output not full block")
	}
	decryptBlockGo(c.dec, dst, src)
}

// Destroy destroys the encryption and decryption key schedule LockedBuffers
func (c *aesCipher) Destroy() {
	c.encLockedBuffer.Destroy()
	c.decLockedBuffer.Destroy()
}

// DestroyCipher is a helper function that
// calls the appropriate Destroy method of the given block
func DestroyCipher(block cipher.Block) error {

	switch block.(type) {
	case *aesCipher:
		block.(*aesCipher).Destroy()
	case *aesCipherAsm:
		block.(*aesCipherAsm).Destroy()
	case *aesCipherGCM:
		block.(*aesCipherGCM).Destroy()
	default:
		return errors.New("block is not an aesguard cipher")
	}

	return nil
}
