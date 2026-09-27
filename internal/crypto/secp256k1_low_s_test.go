package crypto

import (
	"encoding/hex"
	"math/big"
	"testing"
)

func TestVerifyECDSARejectsHighS(t *testing.T) {
	priv, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	pub, err := PublicKeyFromPrivate(priv)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := SignECDSA(priv, "phase9.3-low-s")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyECDSA(hex.EncodeToString(pub), "phase9.3-low-s", sig) {
		t.Fatal("canonical signature rejected")
	}
	s := new(big.Int).SetBytes(sig[32:])
	high := new(big.Int).Sub(secp256k1OrderForTest(), s)
	highBytes := high.FillBytes(make([]byte, 32))
	mal := append([]byte(nil), sig...)
	copy(mal[32:], highBytes)
	if VerifyECDSA(hex.EncodeToString(pub), "phase9.3-low-s", mal) {
		t.Fatal("high-S signature accepted")
	}
}

func secp256k1OrderForTest() *big.Int { return new(big.Int).Set(N()) }
