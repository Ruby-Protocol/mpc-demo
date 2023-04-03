package mpc

type SecretUnit struct {
	index     uint8
	threshold uint8
	value     []byte
}

type Poly func(index uint8) (unit uint8)
