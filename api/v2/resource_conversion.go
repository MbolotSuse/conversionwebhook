package v2

import (
	v1 "github.com/Mbolotsuse/conversionwebhook/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts a v2 resource to the corresponding hub resource
func (src *Foo) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.Foo)
	dst.Spec.InitialField = src.Spec.InitialField
	// note: if we didn't want AddedField == RemovedField, we could implement another strategy
	dst.Spec.RemovedField = src.Spec.AddedField
	dst.ObjectMeta = src.ObjectMeta
	return nil
}

// ConvertFrom converts a hub resource to the v2 resource
func (dst *Foo) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.Foo)
	dst.Spec.InitialField = src.Spec.InitialField
	dst.Spec.AddedField = src.Spec.RemovedField
	dst.ObjectMeta = src.ObjectMeta
	return nil
}
