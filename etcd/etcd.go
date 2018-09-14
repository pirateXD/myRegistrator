package etcd

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/gliderlabs/registrator/bridge"
	"time"
)

func init() {
	bridge.Register(new(Factory), "etcd")
}

type Factory struct{}

func newClient(host []string) *clientv3.Client {
	// config
	cfg := clientv3.Config{
		Endpoints:            host,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    time.Second,
		DialKeepAliveTimeout: time.Second,
	}

	// create client
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return cli
}

func (f *Factory) New(uri *url.URL) bridge.RegistryAdapter {
	urls := make([]string, 0)
	if uri.Host != "" {
		urls = append(urls, "http://"+uri.Host)
	} else {
		urls = append(urls, "http://127.0.0.1:2379")
	}

	res, err := http.Get(urls[0] + "/version")
	if err != nil {
		log.Fatal("etcd: error retrieving version", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if match, _ := regexp.Match("0\\.4\\.*", body); match == true {
		log.Println("etcd: using v0 client")
		return &EtcdAdapter{client: newClient(urls), path: uri.Path}
	}

	return &EtcdAdapter{client2: newClient(urls), path: uri.Path}
}

type EtcdAdapter struct {
	client  *clientv3.Client
	client2 *clientv3.Client

	path string
}

func (r *EtcdAdapter) Ping() error {
	r.syncEtcdCluster()

	var err error
	if r.client != nil {
		_, err = r.client.MemberList(context.TODO())
	} else {
		_, err = r.client2.MemberList(context.TODO())
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *EtcdAdapter) syncEtcdCluster() {
	var result error
	if r.client != nil {
		result = r.client.Sync(context.TODO())
	} else {
		result = r.client2.Sync(context.TODO())
	}

	if nil != result {
		log.Println("etcd: sync cluster was unsuccessful")
	}
}

func (r *EtcdAdapter) Register(service *bridge.Service) error {
	r.syncEtcdCluster()

	path := r.path + "/" + service.Name + "/" + service.ID
	port := strconv.Itoa(service.Port)
	addr := net.JoinHostPort(service.IP, port)

	var err error
	if r.client != nil {
		var resp *clientv3.LeaseGrantResponse
		if resp, err = r.client.Grant(context.TODO(), int64(service.TTL)); err == nil {
			_, err = r.client.Put(context.TODO(), path, addr, clientv3.WithLease(resp.ID))
		}
	} else {
		var resp *clientv3.LeaseGrantResponse
		if resp, err = r.client2.Grant(context.TODO(), int64(service.TTL)); err == nil {
			_, err = r.client2.Put(context.TODO(), path, addr, clientv3.WithLease(resp.ID))
		}
	}

	if err != nil {
		log.Println("etcd: failed to register service:", err)
	}
	return err
}

func (r *EtcdAdapter) Deregister(service *bridge.Service) error {
	r.syncEtcdCluster()

	path := r.path + "/" + service.Name + "/" + service.ID

	var err error
	if r.client != nil {
		_, err = r.client.Delete(context.TODO(), path)
	} else {
		_, err = r.client2.Delete(context.TODO(), path)
	}

	if err != nil {
		log.Println("etcd: failed to deregister service:", err)
	}
	return err
}

func (r *EtcdAdapter) Refresh(service *bridge.Service) error {

	r.syncEtcdCluster()

	var leaseLeases *clientv3.LeaseLeasesResponse
	var err error
	var client *clientv3.Client
	if r.client != nil {
		client = r.client
	} else {
		client = r.client2
	}

	if leaseLeases, err = client.Leases(context.TODO()); nil != err {
		for _, lease := range leaseLeases.Leases {
			_, err = client.KeepAliveOnce(context.TODO(), lease.ID)
			if err != nil {
				log.Println("etcd: failed to refresh service:", err)
			}
		}
	}

	if err != nil {
		log.Println("etcd: failed to refresh service:", err)
	}
	return err
}

func (r *EtcdAdapter) Services() ([]*bridge.Service, error) {
	return []*bridge.Service{}, nil
}
