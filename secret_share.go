package mpc

import (
	"math/rand"
	"time"
)

func SplitSecret(secret []byte, n, t uint8) []SecretUnit {
	secretLen := len(secret)
	secretUnits := make([]SecretUnit, n)
	for i := uint8(0); i < n; i++ {
		secretUnits[i] = SecretUnit{
			index:     i + 1,
			threshold: t,
			value:     make([]byte, secretLen),
		}
	}

	for i := 0; i < secretLen; i++ {
		parts := split(secret[i], n, t)

		for j := uint8(0); j < n; j++ {
			secretUnits[j].value[i] = parts[j]
		}
	}
	return secretUnits
}

func CombainParts(parts []SecretUnit) []byte {
	threshold := int(parts[0].threshold)
	secretLen := len(parts[0].value)
	coeffs, values := makeMatrix(parts, threshold, secretLen)
	coeffs, values = extinction(coeffs, values)
	return remainder(coeffs, values)
}

func makeMatrix(parts []SecretUnit, threshold, secretLen int) (coeffs [][]int, values [][]int) {
	coeffs = make([][]int, threshold)
	values = make([][]int, threshold)
	for i := 0; i < threshold; i++ {
		// make coefficient matrix
		coeffsRow := make([]int, threshold)
		coeffsRow[0] = 1
		index := int(parts[i].index)
		mult := index

		for j := 1; j < threshold; j++ {
			coeffsRow[j] = mult
			mult *= index
		}
		coeffs[i] = coeffsRow

		// make augmented matrix
		values[i] = make([]int, secretLen)
		for j := 0; j < secretLen; j++ {
			values[i][j] = int(parts[i].value[j])
		}
	}
	return
}

func split(M, n, t uint8) []byte {
	parts := make([]byte, n)
	poly := makePoly(M, t)

	for i := uint8(0); i < n; i++ {
		parts[i] = poly(i + 1)
	}

	return parts
}

func makePoly(M, t uint8) Poly {
	coeffs := make([]int, t)
	coeffs[0] = int(M)

	for i := 1; i < int(t); i++ {
		coeffs[i] = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(256)
	}

	return func(index uint8) uint8 {
		part := int(M)
		mult := int(index)

		for i := uint8(1); i < t; i++ {
			part += mult * coeffs[i]
			mult *= int(index)
		}
		return uint8(part % 256)
	}
}

func extinction(coeffs [][]int, values [][]int) ([][]int, [][]int) {
	if len(coeffs) == 1 {
		return coeffs, values
	}

	valueLen := len(values[0])
	newCoeffs := make([][]int, len(coeffs)-1)
	newValues := make([][]int, len(coeffs)-1)

	for i := 1; i < len(coeffs); i++ {
		coeffsRow := make([]int, len(coeffs)-1)
		base := coeffs[0][len(coeffs)-1]
		mult := coeffs[i][len(coeffs)-1]
		for j := 0; j < len(coeffs)-1; j++ {
			coeffsRow[j] = coeffs[i][j]*base - coeffs[0][j]*mult
		}
		newCoeffs[i-1] = coeffsRow

		valueRow := make([]int, valueLen)
		for j := 0; j < valueLen; j++ {
			valueRow[j] = values[i][j]*base - values[0][j]*mult
		}
		newValues[i-1] = valueRow
	}

	return extinction(newCoeffs, newValues)
}

func remainder(coeffs [][]int, values [][]int) []byte {
	if len(coeffs) != 1 || len(values) != 1 {
		return []byte{}
	}
	r := coeffs[0][0]
	value := values[0]
	result := make([]byte, len(value))
	for i := 0; i < len(value); i++ {
		v := value[i] / r
		v = v % 256
		if v < 0 {
			v += 256
		}
		result[i] = uint8(v)
	}
	return result
}
