package aes

import "testing"

func TestAES(t *testing.T) {

	key := DEFAULT_KEY //AESKey()

	msg := "hello world,你好 世界!"
	encrypted, err := EncryptToString(msg, key)
	if err != nil {
		t.Error(err)
	}

	decrypted, err := DecryptToString(encrypted, key)
	if err != nil {
		t.Error(err)
	}
	//plain := string(decrypted)
	if msg != decrypted {
		t.Errorf("Faild to decrypt:%s", decrypted)
	}

}
