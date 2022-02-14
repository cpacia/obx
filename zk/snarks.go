package zk

import "bytes"

type PrivateParams struct {
	Inputs []struct {
		KeyHash         []byte
		Amount          uint64
		Nonce           []byte
		CommitmentIndex int
		MerkleProof     []byte
		Threshold       int
		Pubkeys         [][]byte
		Signatures      [][]byte
	}
	Outputs []struct {
		KeyHash []byte
		Amount  uint64
		Nonce   []byte
	}
}

type PublicParams struct {
	Anchor            []byte
	SigHash           []byte
	OutputCommitments [][]byte
	Nullifiers        [][]byte
	Fee               uint64
	Coinbase          uint64
}

func StandardCircuit(priv PrivateParams, pub PublicParams) bool {
	inVal := 0

	for i, in := range priv.Inputs {
		// First obtain the hash of the UTXO
		commitmentHash := Hash(Serialize(in.KeyHash, in.Amount, in.Nonce))

		// Then validate the merkle proof
		if !ValidateMerkleProof(commitmentHash, in.CommitmentIndex, in.MerkleProof, pub.Anchor) {
			return false
		}

		// Validate that the provided threshold and pubkeys hash to the keyHash
		calculatedKeyHash := CalculateAddress(in.Threshold, in.Pubkeys)
		if calculatedKeyHash != in.KeyHash {
			return false
		}

		// Validate the signature(s)
		if !ValidateMultiSignature(in.Pubkeys, in.Signatures, pub.SigHash) {
			return false
		}

		// Validate that the nullifier is calculated correctly.
		calculatedNullifier := Hash(in.CommitmentIndex, in.Threshold, in.Pubkeys, in.Nonce)
		if !bytes.Equal(calculatedNullifier, pub.Nullifiers[i]) {
			return false
		}

		// Total up the input amounts
		inVal += in.Amount

	}

	outVal := 0
	for i, out := range priv.Outputs {
		// Make sure the OutputCommitment provided in the PublicParams
		// actually matches the calculated output commitment. This prevents
		// someone from putting a different output hash containing a
		// different amount in the transactions.
		outputCommitment := Hash(Serialize(out.KeyHash, out.Amount, out.Nonce))
		if !bytes.Equal(outputCommitment, pub.OutputCommitments[i]) {
			return false
		}

		outVal += out.Amount
	}

	// Verify the transactions is not spending more than it is allowed to
	if outVal+pub.Fee > inVal+pub.Coinbase {
		return false
	}

	return true
}
