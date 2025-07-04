package tcrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"testing"
)

func TestECDH(t *testing.T) {
	fmt.Printf("--ECC Parameters--\n")
	fmt.Printf(" Name: %s\n", elliptic.P256().Params().Name)
	fmt.Printf(" N: %x\n", elliptic.P256().Params().N)
	fmt.Printf(" P: %x\n", elliptic.P256().Params().P)
	fmt.Printf(" Gx: %x\n", elliptic.P256().Params().Gx)
	fmt.Printf(" Gy: %x\n", elliptic.P256().Params().Gy)
	fmt.Printf(" Bitsize: %x\n\n", elliptic.P256().Params().BitSize)

	priva, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privb, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	puba := priva.PublicKey
	pubb := privb.PublicKey

	fmt.Printf("\nPrivate key (Alice) %x", priva.D)
	fmt.Printf("\nPrivate key (Bob) %x\n", privb.D)

	fmt.Printf("\nPublic key (Alice) (%x,%x)", puba.X, puba.Y)
	fmt.Printf("\nPublic key (Bob) (%x %x)\n", pubb.X, pubb.Y)

	a, _ := puba.Curve.ScalarMult(puba.X, puba.Y, privb.D.Bytes())
	shared1 := sha256.Sum256(a.Bytes())

	b, _ := pubb.Curve.ScalarMult(pubb.X, pubb.Y, priva.D.Bytes())

	shared2 := sha256.Sum256(b.Bytes())

	fmt.Printf("\nShared key (Alice) %x\n", shared1)
	fmt.Printf("\nShared key (Bob)  %x\n", shared2)
}

func TestEncrypt(t *testing.T) {
	priva, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	cipher, _ := encrypt([]byte("Hello world!"), &priva.PublicKey, rand.Reader)
	fmt.Printf("ciphertext len= %d\n", len(cipher))
	plain, _ := decrypt(cipher, priva)
	fmt.Println(string(plain))
}

func parsePKCS8b64(pkcs8B64 string) *ecdsa.PrivateKey {
	x, _ := base64.StdEncoding.DecodeString(pkcs8B64)
	y, _ := x509.ParsePKCS8PrivateKey(x)
	if privKey, ok := y.(*ecdsa.PrivateKey); ok {
		return privKey
	}
	panic("Cant parse common key")
}

func encrypt(plain []byte, key *ecdsa.PublicKey, rand io.Reader) ([]byte, error) {
	tmpKey, err := ecdsa.GenerateKey(key.Curve, rand)
	if err != nil {
		return nil, err
	}
	a, _ := key.Curve.ScalarMult(key.X, key.Y, tmpKey.D.Bytes())
	shared1 := sha256.Sum256(a.Bytes())
	block, err := aes.NewCipher(shared1[:])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plain, nil)

	tmpPubKeyBytes, err := x509.MarshalPKIXPublicKey(tmpKey.Public())
	if err != nil {
		return nil, err
	}

	var outBuff bytes.Buffer
	lenBuff := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenBuff, uint16(len(tmpPubKeyBytes)))
	outBuff.Write(lenBuff)
	outBuff.Write(tmpPubKeyBytes)
	binary.LittleEndian.PutUint16(lenBuff, uint16(len(nonce)))
	outBuff.Write(lenBuff)
	outBuff.Write(nonce)
	outBuff.Write(ciphertext)
	return outBuff.Bytes(), nil
}

func decrypt(encData []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	inBuff := bytes.NewBuffer(encData)
	lenBuff := make([]byte, 2)
	if _, err := inBuff.Read(lenBuff); err != nil {
		return nil, err
	}
	tmpPubKeyBytes := make([]byte, binary.LittleEndian.Uint16(lenBuff))
	if _, err := inBuff.Read(tmpPubKeyBytes); err != nil {
		return nil, err
	}
	if _, err := inBuff.Read(lenBuff); err != nil {
		return nil, err
	}
	nonce := make([]byte, binary.LittleEndian.Uint16(lenBuff))
	if _, err := inBuff.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := make([]byte, inBuff.Len())
	if _, err := inBuff.Read(ciphertext); err != nil {
		return nil, err
	}

	pubKey, err := x509.ParsePKIXPublicKey(tmpPubKeyBytes)
	if err != nil {
		return nil, err
	}
	ecPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid pub key type: %T", pubKey)
	}

	b, _ := key.Curve.ScalarMult(ecPubKey.X, ecPubKey.Y, key.D.Bytes())
	shared := sha256.Sum256(b.Bytes())

	block, err := aes.NewCipher(shared[:])
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
