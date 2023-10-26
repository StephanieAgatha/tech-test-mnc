package manager

import "mnc-test/repository"

type RepoManager interface {
	UserCredRepo() repository.UserCredential
	MerchantRepo() repository.MerchantRepository
	TransactionRepo() repository.TransactionRepository
	TransferRepo() repository.TransferRepository
}

type repoManager struct {
	im InfraManager
}

func (r *repoManager) UserCredRepo() repository.UserCredential {
	return repository.NewUserCredentials(r.im.Connect())
}

func (r *repoManager) MerchantRepo() repository.MerchantRepository {
	//TODO implement me
	return repository.NewMerchantRepository(r.im.Connect())
}

func (r *repoManager) TransactionRepo() repository.TransactionRepository {
	//TODO implement me
	return repository.NewTransactionRepository(r.im.Connect())
}

func (r *repoManager) TransferRepo() repository.TransferRepository {
	//TODO implement me
	return repository.NewTransferRepository(r.im.Connect())
}

func NewRepoManager(im InfraManager) RepoManager {
	return &repoManager{
		im: im,
	}
}
