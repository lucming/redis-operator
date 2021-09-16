package service

import (
	redisfailoverv1 "github.com/spotahome/redis-operator/api/redisfailover/v1"
	redisfailoverclientset "github.com/spotahome/redis-operator/client/k8s/clientset/versioned"
	"github.com/spotahome/redis-operator/log"
	"github.com/spotahome/redis-operator/service/k8s"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"testing"
)

func TestRedisFailoverKubeClient_EnsureSlaveService(t *testing.T) {
	type fields struct {
		K8SService k8s.Services
		logger     log.Logger
	}
	type args struct {
		rf        *redisfailoverv1.RedisFailover
		labels    map[string]string
		ownerRefs []metav1.OwnerReference
	}

	kubehome := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubehome)
	stdclient, _ := kubernetes.NewForConfig(config)
	customclient, _ := redisfailoverclientset.NewForConfig(config)
	aeClientset, _ := apiextensionsclientset.NewForConfig(config)
	svc := k8s.New(stdclient, customclient, aeClientset, log.Dummy)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add slave service",
			fields: fields{
				K8SService: svc,
				logger:     log.Dummy,
			},
			args: args{
				rf: &redisfailoverv1.RedisFailover{
					TypeMeta: metav1.TypeMeta{
						Kind:       "RedisFailover",
						APIVersion: "databases.spotahome.com/v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "redisfailover",
						Namespace: "redis",
					},
					Spec: redisfailoverv1.RedisFailoverSpec{
						Redis: redisfailoverv1.RedisSettings{
							Replicas: 2,
						},
						Sentinel: redisfailoverv1.SentinelSettings{
							Replicas: 1,
						},
						SlaveExported: true,
					},
				},
				labels:    nil,
				ownerRefs: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisFailoverKubeClient{
				K8SService: tt.fields.K8SService,
				logger:     tt.fields.logger,
			}
			if err := r.EnsureSlaveService(tt.args.rf, tt.args.labels, tt.args.ownerRefs); (err != nil) != tt.wantErr {
				t.Errorf("EnsureSlaveService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
