package encrypt

import (
	"encoding/hex"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func TestRandString(t *testing.T) {
	str, _ := RandString(12)
	log.Info().Str("rand string", str).Send()
}

func TestSymmetricEncryptSB(t *testing.T) {
	key, _ := RandString(32)
	content := "hello world."
	en := SymmetricEncryptSB(key, []byte(content))
	de, _ := SymmetricDecryptSB(key, en)
	log.Info().Bytes("de", de).Send()
}

func TestAsymmetricEncryptSB(t *testing.T) {
	pub, pri := RSACreate()
	content := "hello world."
	en, _ := AsymmetricEncryptSB(pub, []byte(content))
	de, _ := AsymmetricDecryptSB(pri, en)
	log.Info().Bytes("de", de).Send()
}

func TestProd(t *testing.T) {
	pub, pri := RSACreate()
	realKey, _ := RandString(32)

	log.Info().Str("real key", realKey).Send()
	content := "hello world."
	enKeyBytes, _ := AsymmetricEncryptSB(pub, []byte(realKey))
	accessKey := hex.EncodeToString(enKeyBytes)


	bk, _ := hex.DecodeString(accessKey)
	realKeyBytes, _ := AsymmetricDecryptSB(pri, bk)
	log.Info().Bytes("real key", realKeyBytes).Send()

	log.Info().Str("aes key", accessKey).Send()
	en := SymmetricEncryptBB(realKeyBytes, []byte(content))

	de, _ := SymmetricDecryptBB(realKeyBytes, en)
	log.Info().Bytes("de", de).Send()
}
/**
 	1. create rsa pub pri
	2. get pub
	3. make rand string (realKey)
	4. encrypt realKey send to server.
	5. used aes en/de content
 */