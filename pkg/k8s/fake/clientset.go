package fake

import (
	"context"

	tektonv1beta1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tektonv1beta1fakeclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1beta1/fake"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/meta/testrestmapper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apimachinery/pkg/watch"
	fakediscovery "k8s.io/client-go/discovery/fake"
	fakek8sclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/testing"

	"github.com/ocurity/dracon/pkg/k8s"
)

// NewSchemeAndCodecs returns a new scheme populated with the types defined in
// clientSetSchemes.
func NewSchemeAndCodecs() (*runtime.Scheme, *serializer.CodecFactory, error) {
	scheme := runtime.NewScheme()

	// register core V1 K8s APIs
	utilruntime.Must(fakek8sclient.AddToScheme(scheme))

	// register the Tekton V1Beta1 APIs
	if err := tektonv1beta1api.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	codecs := serializer.NewCodecFactory(scheme)
	return scheme, &codecs, nil
}

var _ k8s.ClientInterface = (*DraconClientSet)(nil)

// ApplyHookType is the signature of the function that will be called to mock
// an Apply invocation
type ApplyHookType func(clientset ClientsetSubset, obj runtime.Object, namespace string, forceConflicts bool) error

// ClientsetSubset is just the fake clientset struct
type ClientsetSubset struct {
	*fakek8sclient.Clientset
	*tektonv1beta1fakeclient.FakeTektonV1beta1
}

// DraconClientSet is a mock implementation of the `k8s.Clientset`
type DraconClientSet struct {
	ClientsetSubset
	discovery      *fakediscovery.FakeDiscovery
	ApplyHook      ApplyHookType
	MetaRESTMapper meta.RESTMapper
}

// NewFakeTypedClient returns a mock K8s client that implements the
// `k8s.ClientInterface` and a `meta.RESTMapper` implementation that can return
// a correct response for all known types offered by the `k8s.ClientInterface`.
func NewFakeTypedClient(objects ...runtime.Object) (DraconClientSet, error) {
	return NewFakeTypedClientWithApplyHook(
		func(_ ClientsetSubset, _ runtime.Object, _ string, _ bool) error { return nil },
		objects...,
	)
}

// NewFakeTypedClientWithApplyHook returns a mock client that implements the
// `k8s.ClientInterface` and a `meta.RESTMapper` implementation that can return
// a correct response for all known types offered by the `k8s.ClientInterface`.
func NewFakeTypedClientWithApplyHook(applyHook ApplyHookType, objects ...runtime.Object) (DraconClientSet, error) {
	scheme, codecs, err := NewSchemeAndCodecs()
	if err != nil {
		return DraconClientSet{}, err
	}

	objectTracker := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := objectTracker.Add(obj); err != nil {
			panic(err)
		}
	}

	fakeCoreK8sClient := &fakek8sclient.Clientset{}
	fakeCoreK8sClient.AddReactor("*", "*", testing.ObjectReaction(objectTracker))
	fakeCoreK8sClient.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := objectTracker.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})

	return DraconClientSet{
		discovery: &fakediscovery.FakeDiscovery{
			Fake: &fakeCoreK8sClient.Fake,
			FakedServerVersion: &version.Info{
				Major: "1",
				Minor: "28",
			},
		},
		ClientsetSubset: ClientsetSubset{
			Clientset: fakeCoreK8sClient,
			FakeTektonV1beta1: &tektonv1beta1fakeclient.FakeTektonV1beta1{
				Fake: &fakeCoreK8sClient.Fake,
			},
		},
		ApplyHook: applyHook,
		MetaRESTMapper: restmapper.NewShortcutExpander(
			testrestmapper.TestOnlyStaticRESTMapper(scheme), fakeCoreK8sClient.Discovery(),
		),
	}, nil
}

// Apply mocks the `kubectl apply`
func (f DraconClientSet) Apply(_ context.Context, obj runtime.Object, namespace string, forceConflicts bool) error {
	return f.ApplyHook(f.ClientsetSubset, obj, namespace, forceConflicts)
}

// RESTMapper returns an instance implementing the `meta.RESTMapper` interface
func (f DraconClientSet) RESTMapper() meta.RESTMapper {
	return f.MetaRESTMapper
}
