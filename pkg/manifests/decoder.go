package manifests

import (
	tektonV1Beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	jsonSerializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var K8sObjDecoder runtime.Decoder
var CodecFactory serializer.CodecFactory
var TektonV1Beta1ObjEncoder runtime.Encoder

func init() {
	sch := runtime.NewScheme()
	if err := scheme.AddToScheme(sch); err != nil {
		panic(err)
	}
	if err := tektonV1Beta1.AddToScheme(sch); err != nil {
		panic(err)
	}
	CodecFactory = serializer.NewCodecFactory(sch)
	K8sObjDecoder = CodecFactory.UniversalDeserializer()
	serializer := jsonSerializer.NewYAMLSerializer(jsonSerializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	TektonV1Beta1ObjEncoder = CodecFactory.EncoderForVersion(serializer, tektonV1Beta1.SchemeGroupVersion)
}
