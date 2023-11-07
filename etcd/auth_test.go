package main

import (
	"context"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestAuth(t *testing.T) {
	if _, err := cli.RoleAdd(context.TODO(), "root"); err != nil {
		t.Error(err)
		return
	}
	if _, err := cli.UserAdd(context.TODO(), "root", "123"); err != nil {
		t.Error(err)
		return
	}
	if _, err := cli.UserGrantRole(context.TODO(), "root", "root"); err != nil {
		t.Error(err)
		return
	}

	if _, err := cli.RoleAdd(context.TODO(), "r"); err != nil {
		t.Error(err)
		return
	}

	if _, err := cli.RoleGrantPermission(
		context.TODO(),
		"r",   // role name
		"foo", // key
		"zoo", // range end
		clientv3.PermissionType(clientv3.PermReadWrite),
	); err != nil {
		t.Error(err)
		return
	}

	if _, err := cli.UserAdd(context.TODO(), "u", "123"); err != nil {
		t.Error(err)
		return
	}
	if _, err := cli.UserGrantRole(context.TODO(), "u", "r"); err != nil {
		t.Error(err)
		return
	}
	if _, err := cli.AuthEnable(context.TODO()); err != nil {
		t.Error(err)
		return
	}

	cliAuth, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		Username:    "u",
		Password:    "123",
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer cliAuth.Close()

	if _, err := cliAuth.Put(context.TODO(), "foo1", "bar"); err != nil {
		t.Error(err)
		return
	}

	_, err = cliAuth.Txn(context.TODO()).
		If(clientv3.Compare(clientv3.Value("zoo1"), ">", "abc")).
		Then(clientv3.OpPut("zoo1", "XYZ")).
		Else(clientv3.OpPut("zoo1", "ABC")).
		Commit()
	t.Log(err)

	rootCli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		Username:    "root",
		Password:    "123",
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer rootCli.Close()

	resp, err := rootCli.RoleGet(context.TODO(), "r")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("user u permission: key %q, reange end %q\n", resp.Perm[0].Key, resp.Perm[0].RangeEnd)

	if _, err := rootCli.AuthDisable(context.TODO()); err != nil {
		t.Error(err)
	}

}
