package install

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
	metafuzzer "k8s.io/apimachinery/pkg/apis/meta/fuzzer"

	restaurantfuzzer "github.com/zeroisme/pizza-apiserver/pkg/apis/restaurant/fuzzer"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForAPIGroup(t, Install, fuzzer.MergeFuzzerFuncs(metafuzzer.Funcs, restaurantfuzzer.Funcs))
}
