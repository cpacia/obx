OBX Protocol Documentation
==========================

### Staking
A node that wishes to validate transactions must first stake coins on the network. He does this by creating and broadcasting a `StakeTransaction`.
The `StakeTransaction` proves that the originator of the transaction owns some utxo in the utxo without revealing which specific utxo it is. It also
reveals the utxo amount publicly which is used to weight the stake. Finally, the utxo's nullifier is also made public.

When a `StakeTransaction` is confirmed in a block, the validator's ID is added to the set of validating nodes, weighted by the stake, and nodes store the associated nullifier.
Since all spends must include a nullifier in the transaction, validating nodes can determine when stake has been spent and remove the validator from the validator set.

### Block Creation
After a block is mined all nodes will pseudorandomly and deterministically select a subset of six validators from the weighted validator set, using the next block's height to seed the selection function.
The first validator will have five seconds of local time to create and broadcast a block. If he fails to broadcast a block the next validator from the subset will have five seconds to do so, and so on.
If all six validators fail to create a block then any validator will be allowed to create one.

To avoid having every validator broadcast a block and flood the network after the 30 second timeout has been reached, the software will have each node attempt to broadcast a block only
if no block has arrived within a random time interval.

### Block Finalization
Once blocks have been created all nodes will poll the weighted set of validators, using the avalanche consensus mechanism, to come to agreement on the validity of the block and resolve
conflicts between conflicting blocks. 

### Coin Generation
Some UTXOs will be created out of the gate and seeded in the genesis block. From there stakers can earn interest on their staked coins by periodically broadcasting a `CoinbaseTransaction`
in which they pay themselves with newly generated coins. These `CoinbaseTransaction`s will only be considered valid by other validating nodes if the validator has met a threshold
for uptime and latency and a minimum time has passed since the previous coinbase.

### Bootstrapping
Bitcoin uses total accumulated work to distinguish between conflicting forks. If you download every block you can find, regardless of the fork, then you can figure out which one is the 'valid' fork by looking
at total work.

Avalanche does not have a notion of accumulated work. While you can use the avalanche consensus mechansim to poll
stakers and ask them if a particular fork is 'valid', the set of stakers themselves are determined by the choice of fork. This means that bootstrapping is inherently
going to be a semi-trusted process. 

The software is going to give the user the option of adding trusted bootstrap nodes. The node will poll them to get their 
best block ID. If there is a conflict it will pause execution and ask the user to manually select a chain via
the command line. 

Absent trusted bootstrap node the software will poll random nodes on the network for their best blocks. If they all return the
same value then it will boostrap to that chain. If there is a conflict it will pause execution and ask the user to decide.


### Transaction Formats:
```protobuf
syntax = "proto3";

message Output {
    bytes commitment      = 1;
    bytes ephemeralPubkey = 2;
    bytes ciphertext      = 3;
}

message StandardTransaction {
    repeated Output outputs =   1;
    uint64 fee                = 2;
    repeated bytes nullifiers = 3;
    bytes anchor              = 4;
    bytes proof               = 5;
}

message CoinbaseTransaction {
    bytes validatorID       = 1;
    uint64 newCoins         = 2;
    repeated Output outputs = 3;
    bytes signature         = 4;
    bytes proof             = 5;
}

message StakeTransaction {
    bytes validatorID = 1;
    uint64 amount     = 2;
    bytes nullifier   = 3;
    bytes signature   = 4;
    bytes proof       = 5;
}
```

### Block Format
```protobuf
syntax = "proto3";

message BlockHeader {
    uint64 version      = 1;
    bytes parent        = 2;
    int64 timestamp     = 3;
    bytes merkleRoot    = 4;
    bytes validatorRoot = 5;
    bytes anchor        = 6;
    bytes producerID    = 7;
    bytes signature     = 8;
}

message Block {
    BlockHeader header                = 1;
    repeated Transaction transactions = 2;
    
    message Transaction {
        oneof data {
            StandardTransaction standardTransaction = 1;
            CoinbaseTransaction coinbaseTransaction = 2;
            StakeTransaction    stakeTransaction    = 3;
        }
    } 
}
```

### Address Format:
```
keyHash = Sha256(<threshold><pubkeys...>)
address = bech32Encode(keyHash || viewPubkey)
```

### Output Commitment Format:
```
commitment = Sha256(keyHash || amount || nonce)
```

### Nullifier Format
```
nullifier = Sha256(commitmentIndex || keyHashPreimage || nonce)
```

### Output Ciphertext Format
```
sharedSecret = ECDH(ephemeralPrivkey, viewPubkey)
sharedSecret = ECDH(ephemeralPubkey, viewPrivkey)
cipherText = Encrypt(commitmentPreimage, sharedSecret)
```

### Circuits (in Go):
```go
type PrivateParams struct {
        Inputs []struct {
            KeyHash []byte
            Amount uint64
            Nonce []byte
            CommitmentIndex int
            MerkleProof []byte
            Threshold int
            Pubkeys [][]byte
            Signatures [][]byte
        }
        Outputs []struct {
            KeyHash []byte
            Amount uint64
            Nonce []byte
        }
}

type PublicParams struct {
        Anchor []byte
        SigHash []byte
        OutputCommitments [][]byte
        Nullifiers [][]byte
        Fee uint64
        Coinbase uint64
}

// The standard circuit proves, in zero knowledge, that:
// - The inputs exist in the UTXO set
// - That the spender is authorized to spend the inputs
// - That the transactions is not spending more than its allowed to
// - That the nullifier is calculated correctly.
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
	    if outVal + pub.Fee > inVal + pub.Coinbase {
		        return false
	    }
 
        return true
}

type PrivateParams struct {
        KeyHash []byte
        Amount uint65
        Nonce []byte
        CommitmentIndex int
        MerkleProof []byte
        Threshold int
        Pubkeys [][]byte
        Signatures [][]byte
}

type PublicParams struct {
        Anchor []byte
        SigHash []byte
        Amount uint64
        Nullifier []byte
}

// The StakingCircuit proves, in zero knowledge, that:
// - Some UTXO exists in the UTXO
// - The creator of the proof is the owner of the UTXO
// - The public nullifier is calculated correctly
// - The public amount is the correct amount associated with the UTXO
func StakingCircuit(priv PrivateParams, pub PublicParams) bool {
        // First obtain the hash of the UTXO
        commitmentHash := Hash(Serialize(priv.KeyHash, priv.Amount, priv.Nonce))
        
        // Then validate the merkle proof
        if !ValidateMerkleProof(commitmentHash, priv.CommitmentIndex, priv.MerkleProof, pub.Anchor) {
                return false
        }
        
        // Validate that the provided threshold and pubkeys hash to the keyHash 
        calculatedKeyHash := CalculateAddress(priv.Threshold, priv.Pubkeys)
        if calculatedKeyHash != priv.KeyHash {
                return false
        }
        
        // Validate the signature(s) 
        if !ValidateMultiSignature(priv.Pubkeys, priv.Signatures, pub.SigHash) {
                return false
        }
        
        // Validate that the nullifier is calculated correctly.
        calculatedNullifier := Hash(priv.CommitmentIndex, priv.Threshold, priv.Pubkeys, priv.Nonce)
        if !bytes.Equal(calculatedNullifier, pub.Nullifier) {
                return false
        }
	
	    // Validate that the public amount matches the private amount
        if priv.Amount != pub.Amount {
            return false
        }

        return true
}
```
