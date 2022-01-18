package keystore_test

import (
	"embed" // Calls init function.
	"testing"

	"github.com/cpustejovsky/service3-video/foundation/keystore"
	"github.com/cpustejovsky/service3-video/foundation/testfmt"
)
//go:embed *.pem
var keyDocs embed.FS

func TestRead(t *testing.T) {
	t.Log("Given the need to parse a directory of private key files.")
	{
		testID := 0
		t.Log(testfmt.Info(testID, "When handling a directory of keyfile(s)."))
		{
			ks, err := keystore.NewFS(keyDocs)
			msg := "Should be able to construct key store"
			if err != nil {
				t.Fatal(testfmt.Error(testID, msg, err))
			}
			t.Logf(testfmt.Success(testID, msg))

			const keyID = "test"
			pk, err := ks.PrivateKey(keyID)
			msg = "Should be able to find key in store"
			if err != nil {
				t.Fatal(testfmt.Error(testID, msg, err))
			}
			t.Logf(testfmt.Success(testID, msg))

			msg = "Should be able to validate the key"
			if err := pk.Validate(); err != nil {
				t.Fatal(testfmt.Error(testID, msg, err))
			}
			t.Logf(testfmt.Success(testID, msg))
		}
	}
}
