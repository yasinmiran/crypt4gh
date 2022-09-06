package headers

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/neicnordic/crypt4gh/keys"
)

const crypt4gh_x25519_sec = `-----BEGIN CRYPT4GH ENCRYPTED PRIVATE KEY-----
YzRnaC12MQAGc2NyeXB0ABQAAAAAbY7POWSS/pYIR8zrPQZJ+QARY2hhY2hhMjBfcG9seTEzMDUAPKc4jWLf1h2T5FsPhNUYMMZ8y36ESATXOuloI0uxKxov3OZ/EbW0Rj6XY0pd7gcBLQDFwakYB7KMgKjiCAAA
-----END CRYPT4GH ENCRYPTED PRIVATE KEY-----
`
const crypt4gh_x25519_pub = `-----BEGIN CRYPT4GH PUBLIC KEY-----
y67skGFKqYN+0n+1P0FyxYa/lHPUWiloN4kdrx7J3BA=
-----END CRYPT4GH PUBLIC KEY-----
`

const ssh_ed25519_sec_enc = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABCKYb3joJ
xaRg4JDkveDbaTAAAAEAAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIA65hmgJeJakva2c
tMpwAqifM/904s6O1zkwLeS5WiDDAAAAoLwLn+qb6fvbYvPn5VuK2IY94BGFlxPdsJElH0
qLE4/hhZiDTXKv7sxup9ZXeJ5ZS5pvFRFPqODCBG87JlbpNBra5pbywpyco89Gr+B0PHff
PR84IfM7rbdETegmHhq6rX9HGSWhA2Hqa3ntZ2dDD+HUtzdGi3zRPAFLCF0uy3laaiBItC
VgFxmKhQ85221EUcMSEk6ophcCe8thlrtxjZk=
-----END OPENSSH PRIVATE KEY-----
`

func TestHeaderMarshallingWithNonce(t *testing.T) {

	writerPrivateKey, err := keys.ReadPrivateKey(strings.NewReader(ssh_ed25519_sec_enc), []byte("123123"))
	if err != nil {
		panic(err)
	}

	readerPublicKey, err := keys.ReadPublicKey(strings.NewReader(crypt4gh_x25519_pub))
	if err != nil {
		panic(err)
	}
	var nonce = [12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	magic := [8]byte{}
	copy(magic[:], MagicNumber)
	header := Header{
		MagicNumber:       magic,
		Version:           1,
		HeaderPacketCount: 2,
		HeaderPackets: []HeaderPacket{{
			WriterPrivateKey:       writerPrivateKey,
			ReaderPublicKey:        readerPublicKey,
			PacketLength:           10,
			HeaderEncryptionMethod: X25519ChaCha20IETFPoly1305,
			Nonce:                  &nonce,
			EncryptedHeaderPacket: DataEncryptionParametersHeaderPacket{
				EncryptedSegmentSize: 10,
				PacketType:           PacketType{DataEncryptionParameters},
				DataEncryptionMethod: ChaCha20IETFPoly1305,
				DataKey:              [32]byte{},
			},
		},
			{
				WriterPrivateKey:       writerPrivateKey,
				ReaderPublicKey:        readerPublicKey,
				PacketLength:           10,
				HeaderEncryptionMethod: X25519ChaCha20IETFPoly1305,
				Nonce:                  &nonce,
				EncryptedHeaderPacket: DataEditListHeaderPacket{
					PacketType:    PacketType{DataEditList},
					NumberLengths: 3,
					Lengths:       []uint64{1, 2, 3},
				},
			},
		},
	}
	marshalledHeader, err := header.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if hex.EncodeToString(marshalledHeader) != "637279707434676801000000020000006c000000000000005ee4b32a4b0fb53dc04dcb02aea9d258afd07736e13522ccaaf4077e643c8d1b0102030405060708090a0b0c8f5854ea6eceff229d474a1f35af0c7b9813ccc1ff370a56a630018203f102d99e83bd6e6cad47cc6d8185d1fa9ea800aedad79f47042ca364000000000000005ee4b32a4b0fb53dc04dcb02aea9d258afd07736e13522ccaaf4077e643c8d1b0102030405060708090a0b0c8e5854ea6dceff229c474a1f35af0c7b9a13ccc1ff370a56a530018203f102d9bb97386e42d0695f862312bd04206bb6" {
		t.Fail()
	}
}

func TestNewHeader(t *testing.T) {
	decodedHeader, err := hex.DecodeString("637279707434676801000000020000006c000000000000005ee4b32a4b0fb53dc04dcb02aea9d258afd07736e13522ccaaf4077e643c8d1b0102030405060708090a0b0c8f5854ea6eceff229d474a1f35af0c7b9813ccc1ff370a56a630018203f102d99e83bd6e6cad47cc6d8185d1fa9ea800aedad79f47042ca364000000000000005ee4b32a4b0fb53dc04dcb02aea9d258afd07736e13522ccaaf4077e643c8d1b0102030405060708090a0b0c8e5854ea6dceff229c474a1f35af0c7b9a13ccc1ff370a56a530018203f102d9bb97386e42d0695f862312bd04206bb6")
	if err != nil {
		t.Error(err)
	}
	buffer := bytes.NewBuffer(decodedHeader)
	readerSecretKey, err := keys.ReadPrivateKey(strings.NewReader(crypt4gh_x25519_sec), []byte("password"))
	if err != nil {
		panic(err)
	}
	header, err := NewHeader(buffer, readerSecretKey)
	if err != nil {
		panic(err)
	}
	if fmt.Sprintf("%v", header) != "&{[99 114 121 112 116 52 103 104] 1 2 [{[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 108 0 <nil> {65564 {0} 0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}} {[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 100 0 <nil> {{1} 3 [1 2 3]}}]}" {
		t.Fail()
	}
}

func TestReadHeader(t *testing.T) {
	inFile, err := os.Open("../../test/sample.txt.enc")
	if err != nil {
		t.Error(err)
	}
	readerSecretKey, err := keys.ReadPrivateKey(strings.NewReader(crypt4gh_x25519_sec), []byte("password"))
	if err != nil {
		t.Error(err)
	}
	header, err := ReadHeader(inFile)
	if err != nil {
		t.Error(err)
	}
	buffer := bytes.NewBuffer(header)
	err = inFile.Close()
	if err != nil {
		t.Error(err)
	}
	inFile, err = os.Open("../../test/sample.txt.enc")
	if err != nil {
		t.Error(err)
	}
	header1, err := NewHeader(inFile, readerSecretKey)
	if err != nil {
		t.Error(err)
	}
	header2, err := NewHeader(buffer, readerSecretKey)
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%v", header1) != fmt.Sprintf("%v", header2) {
		t.Fail()
	}
}

func TestHeaderMarshallingWithoutNonce(t *testing.T) {

	writerPrivateKey, err := keys.ReadPrivateKey(strings.NewReader(ssh_ed25519_sec_enc), []byte("123123"))
	if err != nil {
		panic(err)
	}

	readerPublicKey, err := keys.ReadPublicKey(strings.NewReader(crypt4gh_x25519_pub))
	if err != nil {
		panic(err)
	}
	magic := [8]byte{}
	copy(magic[:], MagicNumber)
	header := Header{
		MagicNumber:       magic,
		Version:           1,
		HeaderPacketCount: 2,
		HeaderPackets: []HeaderPacket{{
			WriterPrivateKey:       writerPrivateKey,
			ReaderPublicKey:        readerPublicKey,
			PacketLength:           10,
			HeaderEncryptionMethod: X25519ChaCha20IETFPoly1305,
			EncryptedHeaderPacket: DataEncryptionParametersHeaderPacket{
				EncryptedSegmentSize: 10,
				PacketType:           PacketType{DataEncryptionParameters},
				DataEncryptionMethod: ChaCha20IETFPoly1305,
				DataKey:              [32]byte{},
			},
		},
			{
				WriterPrivateKey:       writerPrivateKey,
				ReaderPublicKey:        readerPublicKey,
				PacketLength:           10,
				HeaderEncryptionMethod: X25519ChaCha20IETFPoly1305,
				EncryptedHeaderPacket: DataEditListHeaderPacket{
					PacketType:    PacketType{DataEditList},
					NumberLengths: 3,
					Lengths:       []uint64{1, 2, 3},
				},
			},
		},
	}
	marshalledHeader, err := header.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if marshalledHeader == nil {
		t.Error(err)
	}
}

func TestHeader_GetDataEncryptionParameterHeaderPackets(t *testing.T) {
	header := Header{
		HeaderPackets: []HeaderPacket{
			{
				EncryptedHeaderPacket: DataEncryptionParametersHeaderPacket{
					EncryptedSegmentSize: 10,
					PacketType:           PacketType{DataEncryptionParameters},
					DataEncryptionMethod: ChaCha20IETFPoly1305,
					DataKey:              [32]byte{},
				},
			},
		},
	}
	packets, err := header.GetDataEncryptionParameterHeaderPackets()
	if err != nil {
		t.Error(err)
	}
	if len(*packets) != 1 {
		t.Fail()
	}
	packet := (*packets)[0]
	dataKey := [32]byte{}
	if packet.EncryptedSegmentSize != 10 ||
		packet.PacketType.PacketType != DataEncryptionParameters ||
		packet.DataEncryptionMethod != ChaCha20IETFPoly1305 ||
		!bytes.Equal(packet.DataKey[:], dataKey[:]) {
		t.Fail()
	}
}

func TestHeader_GetDataEditListHeaderPacket(t *testing.T) {
	header := Header{
		HeaderPackets: []HeaderPacket{
			{
				EncryptedHeaderPacket: DataEditListHeaderPacket{
					PacketType:    PacketType{DataEditList},
					NumberLengths: 2,
					Lengths:       []uint64{10, 100},
				},
			},
		},
	}
	packet := header.GetDataEditListHeaderPacket()
	if packet == nil {
		t.Fail()
	} else if packet.PacketType.PacketType != DataEditList ||
		packet.NumberLengths != 2 ||
		packet.Lengths[0] != 10 ||
		packet.Lengths[1] != 100 {
		t.Fail()
	}
}
