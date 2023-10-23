package security

import (
	"github.com/wumansgy/goEncrypt/aes"
)

func AesEncrypt(raw []byte) ([]byte, int64, error) {
	if !TokenEnabled || len(raw) == 0 {
		return raw, 0, nil
	}

	token, err := chooseToken()
	if err != nil {
		return nil, 0, err
	}
	cypher, err := aes.AesCbcEncryptBase64(raw, token.Raw, nil)
	if err != nil {
		return nil, 0, err
	}
	return []byte(cypher), token.Serial, nil
}

func AesDecrypt(cypher []byte, serial int64) ([]byte, error) {
	if !TokenEnabled || len(cypher) == 0 {
		return cypher, nil
	}
	token, err := matchToken(serial)
	if err != nil {
		return nil, err
	}
	plain, err := aes.AesCbcDecryptByBase64(string(cypher), token.Raw, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
