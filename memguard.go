package aesguard

import (
	"github.com/awnumar/memguard"
)

// prepareMemguard creates key schedule buffers and uint32 arrays using memguard for use by newCipher
func prepareMemguard(key []byte) (*memguard.LockedBuffer, *memguard.LockedBuffer, *memguard.LockedBuffer, []uint32, []uint32, error) {
	// this is how many bytes it takes to represent the enc/dec key schedules
	// each schedule contains n uint32s
	lenKeyScheduleBytes := 4 * (len(key) + 28)

	// create an immutable locked bucffer from the key
	keyBuffer, err := memguard.NewImmutableFromBytes(key)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	encScheduleBuffer, err := memguard.NewMutable(lenKeyScheduleBytes)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	decScheduleBuffer, err := memguard.NewMutable(lenKeyScheduleBytes)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	encScheduleUint32, err := encScheduleBuffer.Uint32()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	decScheduleUint32, err := decScheduleBuffer.Uint32()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// disable core dumps just because
	memguard.DisableUnixCoreDumps()

	return keyBuffer, encScheduleBuffer, decScheduleBuffer, encScheduleUint32, decScheduleUint32, nil

}

// finalizeMemguard destroys the given keyBuffer, and marks the enc and dec schedule buffers as immutable for use by newCipher
func finalizeMemguard(keyBuffer, encScheduleBuffer, decScheduleBuffer *memguard.LockedBuffer) error {
	// destroy key buffer, the original key is no longer needed

	keyBuffer.Destroy()

	// make the key schedules as immutable now that they have been populated
	err := encScheduleBuffer.MakeImmutable()
	if err != nil {
		return err
	}

	err = decScheduleBuffer.MakeImmutable()
	if err != nil {
		return err
	}

	return nil
}
