package etcd

import (
	"testing"

	"github.com/gliderlabs/registrator/bridge"
	"net/url"
)

func Test_newClient(t *testing.T) {

	etcdFactory := new(Factory)
	uri, _ := url.Parse("etcd://172.26.147.228:2379")
	adapter := etcdFactory.New(uri)

	if err := adapter.Ping(); err != nil {
		t.Error("ping error", err)
	}

	service := &bridge.Service{
		ID:     "testId",
		Name:   "testName",
		Port:   0,
		IP:     "127.0.0.1",
		Tags:   []string{"tags"},
		Attrs:  map[string]string{"attr1": "v1"},
		TTL:    30,
		Origin: bridge.ServicePort{},
	}
	if regError := adapter.Register(service); nil != regError {
		t.Error("register error")
	}

}
