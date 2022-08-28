package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Register struct {
	EtcdAddrs   []string
	DialTimeout int

	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo Server
	srvTTL  int64
	cli     *clientv3.Client
	logger  *logrus.Logger
}

func NewRegister(etcdAddrs []string, logger *logrus.Logger) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// 初始化自己的etcd配置
func (r *Register) Register(srcInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error

	if strings.Split(srcInfo.Addr, ":")[0] == "" {
		return nil, errors.New("Invalid ip address")
	}

	// 初始化
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvTTL = ttl
	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})
	go r.keepAlive()
	return r.closeCh, nil
}

// 创建etcd自带的配置
func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	leaseResp, err := r.cli.Grant(ctx, r.srvTTL)
	if err != nil {
		return err
	}

	r.leasesID = leaseResp.ID

	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leasesID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unRegister(); err != nil {
				fmt.Println("unRegister failed error", err)
			}
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				fmt.Println("revoke failed")
			}
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					fmt.Println("register failed")
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					fmt.Println("register failed")
				}
			}
		}
	}
}

func (r *Register) unRegister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegisterPath(r.srvInfo))
	return err
}
