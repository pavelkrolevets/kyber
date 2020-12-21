package timelock

import (
	"errors"

	"golang.org/x/crypto/blake2s"

	"github.com/drand/drand/key"
	"github.com/drand/kyber"
	"github.com/drand/kyber/pairing"
	"github.com/drand/kyber/util/random"
)

type Ciphertext struct {
	//
	RP kyber.Point
	// ciphertext
	C []byte
}

// SigGroup = G2
// KeyGroup = G1
// P random generator of G1
// dist key: s, Ppub = s*P \in G1
// H1: {0,1}^n -> G1
// H2: GT -> {0,1}^n
// ID: Qid = H1(ID) = xP \in G2
// 	secret did = s*Qid \in G2
// Encrypt:
// - random r scalar
// - Gid = e(Ppub, r*Qid) = e(P, P)^(x*s*r) \in GT
// 		 = GidT
// - U = rP \in G1,
// - V = M XOR H2(Gid)) = M XOR H2(GidT)  \in {0,1}^n
// Decrypt:
// - V XOR H2(e(U, did)) = V XOR H2(e(rP, s*Qid))
//   = V XOR H2(e(P, P)^(r*s*x))
//   = V XOR H2(GidT) = M
func Encrypt(s pairing.Suite, basePoint, public kyber.Point, ID, msg []byte) (*Ciphertext, error) {
	if len(msg)>>16 > 0 {
		// we're using blake2 as XOF which only outputs 2^16-1 length
		return nil, errors.New("ciphertext too long")
	}
	hashable, ok := s.G2().Point().(kyber.HashablePoint)
	if !ok {
		return nil, errors.New("point needs to implement hashablePoint")
	}
	Qid := hashable.Hash(ID)
	r := s.G2().Scalar().Pick(random.New())
	rP := s.G1().Point().Mul(r, basePoint)

	// e(Qid, Ppub) = e( H(round), s*P) where s is dist secret key
	Ppub := public
	rQid := key.SigGroup.Point().Mul(r, Qid)
	GidT := key.Pairing.Pair(Ppub, rQid)
	// H(gid)
	hGidT, err := gtToHash(GidT, uint16(len(msg)))
	if err != nil {
		return nil, err
	}
	xored := xor(msg, hGidT)

	return &Ciphertext{
		RP: rP,
		C:  xored,
	}, nil
}

func Decrypt(s pairing.Suite, private kyber.Point, c *Ciphertext) ([]byte, error) {
	gidt := key.Pairing.Pair(c.RP, private)
	hgidt, err := gtToHash(gidt, uint16(len(c.C)))
	if err != nil {
		return nil, err
	}
	return xor(c.C, hgidt), nil
}

func gtToHash(gt kyber.Point, length uint16) ([]byte, error) {
	xof, err := blake2s.NewXOF(length, nil)
	if err != nil {
		return nil, err
	}
	gt.MarshalTo(xof)
	var b = make([]byte, length)
	n, err := xof.Read(b)
	if uint16(n) != length || err != nil {
		return nil, errors.New("couldn't read from xof")
	}
	return b[:], nil
}

func xor(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("wrong xor input")
	}
	res := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = a[i] ^ b[i]
	}
	return res
}
