package etcd

import (
	"context"
	"maden/pkg/shared"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdTransactionRepository struct {
	client EtcdClient
}

func NewEtcdTransactionRepository(client EtcdClient) Transactioner {
	return &EtcdTransactionRepository{client: client}
}

func (etr *EtcdTransactionRepository) PerformTransaction(ctx context.Context, key string, value string, resourceType shared.ResourceType) error {
	txnResp, err := etr.client.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, value)).
		Else(clientv3.OpGet(key)).
		Commit()

	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return &shared.ErrDuplicateResource{ID: key, ResourceType: resourceType}
	}
	return nil
}