package bitcoind

import "testing"

func TestValidation(t *testing.T) {
	config := appConfig{}
	error := validateFlags(&config)
	if error == nil {
		t.Error("Didn't complain about missing cookie file location")
	}
	config = appConfig{rpcCookieFile: "dummy"}
	error = validateFlags(&config)
	if error == nil {
		t.Error("Didn't complain about missing rpc host")
	}
}
