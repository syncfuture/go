package ssecurity

type (
	IEncryptor interface {
		EncryptString(string) (string, error)
		DecryptString(string) (string, error)
		Encrypt([]byte) ([]byte, error)
		Decrypt([]byte) ([]byte, error)
	}
)
