/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package snapshot

import (
	"io"
	"os"

	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/internal/peer/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Client holds client side dependency for the snapshot commands
type Client struct {
	SnapshotClient pb.SnapshotClient
	Signer         common.Signer
	Writer         io.Writer
}

// NewClient creates a client instance
func NewClient(cryptoProvider bccsp.BCCSP) (*Client, error) {
	if err := validatePeerConnectionParameters(); err != nil {
		return nil, err
	}

	snapshotClient, err := common.GetSnapshotClient(peerAddress, tlsRootCertFile)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to retrieve snapshot client")
	}

	signer, err := common.GetDefaultSigner()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to retrieve default signer")
	}

	return &Client{
		Signer:         signer,
		SnapshotClient: snapshotClient,
		Writer:         os.Stdout,
	}, nil
}

func validatePeerConnectionParameters() error {
	switch viper.GetBool("peer.tls.enabled") {
	case true:
		if tlsRootCertFile == "" {
			return errors.New("the required parameter 'tlsRootCertFile' is empty. Rerun the command with --tlsRootCertFile flag")
		}
	case false:
		tlsRootCertFile = ""
	}

	return nil
}