package custominitializer_test

import (
	"testing"
	"time"

	"k8s.io/apiserver/pkg/admission"

	"github.com/zeroisme/pizza-apiserver/pkg/admission/custominitializer"
	"github.com/zeroisme/pizza-apiserver/pkg/generated/clientset/versioned/fake"
	informers "github.com/zeroisme/pizza-apiserver/pkg/generated/informers/externalversions"
)

// TestWantsInternalRestaurantInformerFactory ensures that the informer factory is injected
// when the WantsRestaurantInformerFactory interface is implemented by a plugin.
func TestWantsInternalRestaurantInformerFactory(t *testing.T) {
	cs := &fake.Clientset{}
	sf := informers.NewSharedInformerFactory(cs, time.Duration(1)*time.Second)
	target := custominitializer.New(sf)

	wantRestaurantInformerFactory := &wantRestaurantInformerFactory{}
	target.Initialize(wantRestaurantInformerFactory)
	if wantRestaurantInformerFactory.sf != sf {
		t.Errorf("expected informer factory to be initialized")
	}
}

// wantRestaurantInformerFactory is a test stub that fulfills the WantsRestaurantInformerFactory interface
type wantRestaurantInformerFactory struct {
	sf informers.SharedInformerFactory
}

func (self *wantRestaurantInformerFactory) SetRestaurantInformerFactory(sf informers.SharedInformerFactory) {
	self.sf = sf
}
func (self *wantRestaurantInformerFactory) Admit(a admission.Attributes) error { return nil }
func (self *wantRestaurantInformerFactory) Handles(o admission.Operation) bool { return false }
func (self *wantRestaurantInformerFactory) ValidateInitialization() error      { return nil }

var _ admission.Interface = &wantRestaurantInformerFactory{}
var _ custominitializer.WantsRestaurantInformerFactory = &wantRestaurantInformerFactory{}
