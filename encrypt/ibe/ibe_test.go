package ibe

import (
	"testing"

	"github.com/drand/kyber"
	bls "github.com/drand/kyber-bls12381"
	"github.com/drand/kyber/pairing"
	"github.com/drand/kyber/util/random"
	"github.com/stretchr/testify/require"
)

func newSetting() (pairing.Suite, kyber.Point, []byte, kyber.Point) {
	suite := bls.NewBLS12381Suite()
	P := suite.G1().Point().Base()
	s := suite.G1().Scalar().Pick(random.New())
	Ppub := suite.G1().Point().Mul(s, P)

	ID := []byte("passtherand")
	IDP := suite.G2().Point().(kyber.HashablePoint)
	Qid := IDP.Hash(ID)     // public key
	sQid := Qid.Mul(s, Qid) // secret key
	return suite, Ppub, ID, sQid
}

func TestTimelock(t *testing.T) {
	suite, Ppub, ID, sQid := newSetting()
	msg := []byte("Hello World\n")

	c, err := Encrypt(suite, Ppub, ID, msg)
	require.NoError(t, err)
	msg2, err := Decrypt(suite, Ppub, sQid, c)
	require.NoError(t, err)
	require.Equal(t, msg, msg2)
}

/*func TestAggregateTimelock(t *testing.T) {*/
//suite, Ppub, ID, sQid := newSetting()

//n := 6 // number of proofs
//msgSize := 32
//ciphers := make([]Ciphertext, n)
//plains := make([][]byte, n)
//for i := 0; i < n; i++ {
//msg := make([]byte, msgSize)
//rand.Read(msg)
//plains[i] = msg
//c, err := Encrypt(suite, Ppub, ID, msg)
//require.NoError(t, err)
//ciphers[i] = *c
//}
//aggregated, err := AggregateEncrypt(suite, Ppub, sQid, ciphers)
//require.NoError(t, err)

//plaintexts, err := DecryptAggregate(suite, Ppub, sQid, aggregated)
//require.NoError(t, err)
//require.Equal(t, n, len(plaintexts))
//for i := 0; i < n; i++ {
//require.Equal(t, plains[i], plaintexts[i])
//}
//}

//func BenchmarkTimelockAggregateVerify(b *testing.B) {
//suite, Ppub, ID, sQid := newSetting()
//n := 512
//msgSize := 32
//ciphers := make([]Ciphertext, n)
//plains := make([][]byte, n)
//for i := 0; i < n; i++ {
//msg := make([]byte, msgSize)
//rand.Read(msg)
//plains[i] = msg
//c, err := Encrypt(suite, Ppub, ID, msg)
//if err != nil {
//panic(err)
//}
//ciphers[i] = *c
//}
//aggregated, err := AggregateEncrypt(suite, Ppub, sQid, ciphers)
//if err != nil {
//panic(err)
//}
//b.ResetTimer()
//for i := 0; i < b.N; i++ {
//_, err := DecryptAggregate(suite, Ppub, sQid, aggregated)
//if err != nil {
//panic(err)
//}
//}

//}

//func BenchmarkTimelockAggregateAggregate(b *testing.B) {
//suite, Ppub, ID, sQid := newSetting()
//n := 512
//msgSize := 32
//ciphers := make([]Ciphertext, n)
//plains := make([][]byte, n)
//for i := 0; i < n; i++ {
//msg := make([]byte, msgSize)
//rand.Read(msg)
//plains[i] = msg
//c, err := Encrypt(suite, Ppub, ID, msg)
//if err != nil {
//panic(err)
//}
//ciphers[i] = *c
//}

//b.ResetTimer()
//for i := 0; i < b.N; i++ {
//_, err := AggregateEncrypt(suite, Ppub, sQid, ciphers)
//if err != nil {
//panic(err)
//}
//}

//}

//func BenchmarkTimelockSingleVerify(b *testing.B) {
//suite, Ppub, ID, sQid := newSetting()
//n := 512
//msgSize := 32
//ciphers := make([]Ciphertext, n)
//plains := make([][]byte, n)
//for i := 0; i < n; i++ {
//msg := make([]byte, msgSize)
//rand.Read(msg)
//plains[i] = msg
//c, err := Encrypt(suite, Ppub, ID, msg)
//if err != nil {
//panic(err)
//}
//ciphers[i] = *c
//}
//b.ResetTimer()
//for i := 0; i < b.N; i++ {
//for _, c := range ciphers {
//_, err := Decrypt(suite, Ppub, sQid, &c)
//if err != nil {
//panic(err)
//}
//}

//}

//}

//func BenchmarkTimelockAggregateNoSum(b *testing.B) {
//suite, Ppub, ID, sQid := newSetting()
//n := 512
//msgSize := 32
//ciphers := make([]Ciphertext, n)
//plains := make([][]byte, n)
//for i := 0; i < n; i++ {
//msg := make([]byte, msgSize)
//rand.Read(msg)
//plains[i] = msg
//c, err := Encrypt(suite, Ppub, ID, msg)
//if err != nil {
//panic(err)
//}
//ciphers[i] = *c
//}
//aggregated, err := AggregateEncrypt(suite, Ppub, sQid, ciphers)
//if err != nil {
//panic(err)
//}
//aggregated.skipSum = true
//b.ResetTimer()
//for i := 0; i < b.N; i++ {
//_, err := DecryptAggregate(suite, Ppub, sQid, aggregated)
//if err != nil {
//panic(err)
//}
//}

//}

////
//// BenchmarkTimelockAggregateVerify-8                     2         733865201 ns/op
//// BenchmarkTimelockAggregateAggregate-8                  1        1643196283 ns/op
//// BenchmarkTimelockSingleVerify-8                        2         913203290 ns/op
//// BenchmarkTimelockAggregateNoSum-8                      7         158552998 ns/op

//// aggregate gives 20% verification speed up but GT operation still takes 80% of
/*// the time*/
