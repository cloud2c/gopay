package alipay

import (
	"errors"
	"testing"

	"github.com/cloud2c/gopay"
)

func TestGatewayURLSelectsEnvironment(t *testing.T) {
	if GatewayURL(true) != ProductionGatewayURL {
		t.Fatal("isProd=true must return the production gateway")
	}
	if GatewayURL(false) != SandboxGatewayURL {
		t.Fatal("isProd=false must return the sandbox gateway")
	}
	if GatewayURL(false) == GatewayURL(true) {
		t.Fatal("sandbox and production gateways must differ")
	}
}

func TestBizErrSatisfiesBusinessErr(t *testing.T) {
	err := &BizErr{Code: "40004", Msg: "Business Failed"}
	if !errors.Is(err, gopay.BusinessErr) {
		t.Fatal("BizErr must satisfy errors.Is(err, gopay.BusinessErr)")
	}
}
