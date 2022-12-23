package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var endpoints = []string{"127.0.0.1:2379"}

const (
	dialTimeout    = time.Second * 10
	requestTimeout = time.Second * 100
)

var cli *clientv3.Client

func init() {
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		panic(err)
	}
}

func TestKVPut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.Put(ctx, "sample_key", "sample_value")
	cancel()
	if err != nil {
		t.Error(err)
	}
}

func TestKVPutErrorHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.Put(ctx, "", "sample_value")
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			t.Errorf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			t.Errorf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			t.Logf("client-side error: %v\n", err)
		default:
			t.Errorf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}
}

func TestKVGet(t *testing.T) {
	_, err := cli.Put(context.TODO(), "foo", "sample_value")
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "foo")
	cancel()
	if err != nil {
		t.Error(err)
	}
	for _, ev := range resp.Kvs {
		t.Logf("%s : %s\n", ev.Key, ev.Value)
	}
}

func TestKVGetWithRev(t *testing.T) {
	presp, err := cli.Put(context.TODO(), "foo", "bar1")
	if err != nil {
		t.Error(err)
	}

	_, err = cli.Put(context.TODO(), "foo", "bar2")
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "foo", clientv3.WithRev(presp.Header.Revision))
	cancel()
	if err != nil {
		t.Error(err)
	}
	for _, ev := range resp.Kvs {
		t.Logf("%s : %s\n", ev.Key, ev.Value)
	}
}

func TestKVGetSortedPrefix(t *testing.T) {
	for i := range make([]int, 3) {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		_, err := cli.Put(ctx, fmt.Sprintf("key_%d", i), "value")
		cancel()
		if err != nil {
			t.Error(err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "key", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()
	if err != nil {
		t.Error(err)
	}
	for _, ev := range resp.Kvs {
		t.Logf("%s : %s\n", ev.Key, ev.Value)
	}
}

func TestKVDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	gresp, err := cli.Get(ctx, "key", clientv3.WithPrefix())
	if err != nil {
		t.Error(err)
	}

	dresp, err := cli.Delete(ctx, "key", clientv3.WithPrefix())
	if err != nil {
		t.Error(err)
	}

	t.Log("Deleted all keys: ", int64(len(gresp.Kvs)) == dresp.Deleted)
}

func TestKVCompact(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "foo")
	cancel()
	if err != nil {
		t.Error(err)
	}
	compRev := resp.Header.Revision

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Compact(ctx, compRev)
	cancel()
	if err != nil {
		t.Error(err)
	}
}

func TestKVTxn(t *testing.T) {
	kvc := clientv3.NewKV(cli)

	_, err := kvc.Put(context.TODO(), "key", "xyz")
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = kvc.Txn(ctx).
		If(clientv3.Compare(clientv3.Value("key"), ">", "abc")).
		Then(clientv3.OpPut("key", "XYZ")).
		Else(clientv3.OpPut("key", "ABC")).
		Commit()
	cancel()
	if err != nil {
		t.Error(err)
	}

	gresp, err := kvc.Get(context.TODO(), "key")
	if err != nil {
		t.Error(err)
	}
	for _, ev := range gresp.Kvs {
		t.Logf("%s : %s\n", ev.Key, ev.Value)
	}
}

func TestKVDo(t *testing.T) {
	ops := []clientv3.Op{
		clientv3.OpPut("put-key", "123"),
		clientv3.OpGet("put-key"),
		clientv3.OpPut("put-key", "456"),
	}

	for _, op := range ops {
		if _, err := cli.Do(context.TODO(), op); err != nil {
			t.Error(err)
		}
	}
}
