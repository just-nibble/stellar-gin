package pkg

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

type StellarClient struct {
	Net string
}

func (s *StellarClient) CreateKeyPair() (seed string, address string) {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	seed = pair.Seed()

	address = pair.Address()

	return seed, address
}

func (s *StellarClient) FundAccount(address string) (*string, error) {
	var status string

	if s.Net != "test" {
		status = "failed"
		err := errors.New("only available on testnet")
		if err != nil {
			return nil, err
		}
	}

	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	status = "success"

	return &status, nil

}

func (s *StellarClient) CheckBalance(address string) ([]*horizon.Balance, error) {
	request := horizonclient.AccountRequest{AccountID: address}

	balances := []*horizon.Balance{}

	switch s.Net {
	case "test":
		account, err := horizonclient.DefaultTestNetClient.AccountDetail(request)
		if err != nil {
			return nil, err
		}
		for _, balance := range account.Balances {
			balances = append(balances, &balance)
		}
	default:
		account, err := horizonclient.DefaultPublicNetClient.AccountDetail(request)
		if err != nil {
			return nil, err
		}
		for _, balance := range account.Balances {
			balances = append(balances, &balance)
		}
	}

	return balances, nil
}

func (s *StellarClient) BuildTransaction(dest string, amount string, key string) (*string, error) {
	client := horizonclient.DefaultTestNetClient
	passphrase := network.TestNetworkPassphrase

	if s.Net != "test" {
		client = horizonclient.DefaultPublicNetClient
		passphrase = network.PublicNetworkPassphrase
	}

	// Make sure destination account exists
	destAccountRequest := horizonclient.AccountRequest{AccountID: dest}

	_, err := client.AccountDetail(destAccountRequest)
	if err != nil {
		return nil, errors.New("destination account does not exist")
	}

	// Load the source account
	sourceKP := keypair.MustParseFull(key)
	sourceAccountRequest := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(sourceAccountRequest)
	if err != nil {
		return nil, errors.New("source account does not exist")
	}

	_, err = strconv.Atoi(amount)
	if err != nil {

		return nil, errors.New("invalid amount")
	}

	// Build transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewTimebounds(5, 20), // Use a real timeout in production!
			},
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: dest,
					Amount:      amount,
					Asset:       txnbuild.NativeAsset{},
				},
			},
		},
	)

	if err != nil {
		log.Println(err)
		return nil, errors.New("error building transaction")
	}

	// Sign the transaction to prove you are actually the person sending it.
	tx, err = tx.Sign(passphrase, sourceKP)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error signing transaction")
	}

	// And finally, send it off to Stellar!
	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error submitting transaction")
	}

	return &resp.Hash, nil

}

func (s *StellarClient) BuildTrustLine(key string, assetCode string, issuer string) (*string, error) {
	client := horizonclient.DefaultTestNetClient
	passphrase := network.TestNetworkPassphrase

	if s.Net != "test" {
		client = horizonclient.DefaultPublicNetClient
		passphrase = network.PublicNetworkPassphrase
	}

	// Load the source account
	sourceKP := keypair.MustParseFull(key)
	sourceAccountRequest := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(sourceAccountRequest)
	if err != nil {
		return nil, errors.New("source account does not exist")
	}

	// Build asset
	asset := txnbuild.CreditAsset{Code: assetCode, Issuer: issuer}

	// Build transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(), // Use a real timeout in production!
			},
			Operations: []txnbuild.Operation{
				&txnbuild.SetTrustLineFlags{
					Trustor:       sourceKP.Address(),
					Asset:         asset,
					SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized, txnbuild.TrustLineAuthorizedToMaintainLiabilities},
					ClearFlags:    []txnbuild.TrustLineFlag{txnbuild.TrustLineClawbackEnabled},
					SourceAccount: key,
				},
			},
		},
	)

	if err != nil {
		log.Println(err)
		return nil, errors.New("error building transaction")
	}

	// Sign the transaction to prove you are actually the person sending it.
	tx, err = tx.Sign(passphrase, sourceKP)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error signing transaction")
	}

	// And finally, send it off to Stellar!
	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error submitting transaction")
	}

	return &resp.Hash, nil
}
