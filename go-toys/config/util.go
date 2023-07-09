package config

func xorItemHash(checksum []byte, itemHash []byte) {
	for i := 0; i < 32; i++ {
		checksum[i] = checksum[i] ^ itemHash[i]
	}
}
