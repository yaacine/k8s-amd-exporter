package fakekubelet

import (
	"context"
	"log/slog"
	"testing"

	"net"

	"github.com/openinnovationai/k8s-amd-exporter/internal/kubernetes"
	"github.com/openinnovationai/k8s-amd-exporter/internal/sdk/unittests/testlogs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"k8s.io/client-go/kubernetes/fake"
	podresourcesapi "k8s.io/kubelet/pkg/apis/podresources/v1alpha1"
)

type Service struct {
	podResources []*podresourcesapi.PodResources

	podresourcesapi.UnimplementedPodResourcesListerServer
}

func New(
	tb testing.TB,
	opts ...func(*Options),
) *kubernetes.Client {
	options := Options{
		podResources: make([]*podresourcesapi.PodResources, 0),
		logger:       testlogs.NewLogger(),
	}
	for _, opt := range opts {
		opt(&options)
	}
	svc := &Service{
		podResources: options.podResources,
	}

	return newFake(tb, svc, options, kubernetes.NewClient)
}

func (s *Service) RegisterOn(srv *grpc.Server) {
	podresourcesapi.RegisterPodResourcesListerServer(srv, s)
}

func (s *Service) List(ctx context.Context, req *podresourcesapi.ListPodResourcesRequest) (*podresourcesapi.ListPodResourcesResponse, error) {
	result := podresourcesapi.ListPodResourcesResponse{
		PodResources: s.podResources,
	}
	return &result, nil
}

type Options struct {
	podResources           []*podresourcesapi.PodResources
	logger                 *slog.Logger
	amdCustomResourceNames []string
	nodeName               string
	k8sClient              *fake.Clientset
}

func WithPodResources(podResources []*podresourcesapi.PodResources) func(*Options) {
	return func(options *Options) {
		options.podResources = podResources
	}
}

func WithAMDCustomResourceNames(names []string) func(*Options) {
	return func(options *Options) {
		options.amdCustomResourceNames = names
	}
}

func WithNodeName(name string) func(*Options) {
	return func(options *Options) {
		options.nodeName = name
	}
}

func WithClientSet(client *fake.Clientset) func(*Options) {
	return func(options *Options) {
		options.k8sClient = client
	}
}

type registerer interface {
	RegisterOn(*grpc.Server)
}

func newFake[T any](
	tb testing.TB,
	svc registerer,
	options Options,
	newServiceClient func(*kubernetes.Setup) *T,
) *T {
	tb.Helper()

	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()

	svc.RegisterOn(srv)
	go func() {
		if err := srv.Serve(lis); err != nil {
			tb.Error(err)
		}
	}()
	tb.Cleanup(srv.Stop)

	bufDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return lis.DialContext(ctx)
	}
	opts := []grpc.DialOption{
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient("passthrough://bufnet", opts...)
	if err != nil {
		tb.Fatal(err)
	}

	config := kubernetes.Setup{
		KubeletConn:         conn,
		Logger:              options.logger,
		CustomResourceNames: options.amdCustomResourceNames,
		NodeName:            options.nodeName,
		K8SClient:           options.k8sClient,
	}

	return newServiceClient(&config)
}
