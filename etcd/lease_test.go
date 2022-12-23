package main

import (
	"context"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestLeaseGrant(t *testing.T) {
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		t.Error(err)
	}

	_, err = cli.Put(context.TODO(), "foo", "bar", clientv3.WithLease(resp.ID))
	if err != nil {
		t.Error(err)
	}
}

func TestLeaseRevoke(t *testing.T) {
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		t.Error(err)
	}

	_, err = cli.Put(context.TODO(), "foo", "bar", clientv3.WithLease(resp.ID))
	if err != nil {
		t.Error(err)
	}

	_, err = cli.Revoke(context.TODO(), resp.ID)
	if err != nil {
		t.Error(err)
	}

	gresp, err := cli.Get(context.TODO(), "foo")
	if err != nil {
		t.Error(err)
	}
	t.Log("number of keys:", len(gresp.Kvs))
}

func TestLeaseKeepAlive(t *testing.T) {
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		t.Error(err)
	}

	_, err = cli.Put(context.TODO(), "foo", "bar", clientv3.WithLease(resp.ID))
	if err != nil {
		t.Error(err)
	}

	ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		t.Error(kaerr)
	}

	ka := <-ch
	if ka != nil {
		t.Log("ttl:", ka.TTL)
	} else {
		t.Log("Unexpected NULL")
	}

	time.Sleep(time.Second * 10)
}

func TestLeaseKeepAliveOnce(t *testing.T) {
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		t.Error(err)
	}

	_, err = cli.Put(context.TODO(), "foo", "bar", clientv3.WithLease(resp.ID))
	if err != nil {
		t.Error(err)
	}

	ka, kaerr := cli.KeepAliveOnce(context.TODO(), resp.ID)
	if kaerr != nil {
		t.Error(err)
	}

	t.Log("ttl:", ka.TTL)
}
