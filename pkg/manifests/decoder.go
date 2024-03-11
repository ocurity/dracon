package manifests

import (
	"context"
	"fmt"

	tektonV1Beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	jsonSerializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var K8sObjDecoder runtime.Decoder
var CodecFactory serializer.CodecFactory
var TektonV1Beta1ObjEncoder runtime.Encoder
var BatchV1ObjEncoder runtime.Encoder

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
	BatchV1ObjEncoder = CodecFactory.EncoderForVersion(serializer, batchv1.SchemeGroupVersion)
}

func LoadK8sManifest(ctx context.Context, root, pathOrURI, targetFile string) (runtime.Object, *schema.GroupVersionKind, error) {
	loader, err := NewLoader(root, pathOrURI, targetFile)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: could not resolve path or URI: %w", pathOrURI, err)
	}

	manifestBytes, err := loader.Load(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: could not load path or URI: %w", loader.Path(), err)
	}

	obj, gKV, err := K8sObjDecoder.Decode(manifestBytes, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: could not decode file into a K8s object:3 %w", loader.Path(), err)
	}

	return obj, gKV, nil
}
