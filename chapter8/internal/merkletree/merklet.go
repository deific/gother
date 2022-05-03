package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	transaction2 "gother/chapter8/internal/transaction"
	"gother/chapter8/internal/utils"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	Data      []byte
}

func CreateMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	tempNode := MerkleNode{}

	if left == nil && right == nil {
		tempNode.Data = data
	} else {
		catenateHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(catenateHash)
		tempNode.Data = hash[:]
	}

	tempNode.LeftNode = left
	tempNode.RightNode = right

	return &tempNode
}

func CreateMerkleTree(txs []*transaction2.Transaction) *MerkleTree {
	txsLen := len(txs)
	if txsLen%2 != 0 {
		txs = append(txs, txs[txsLen-1])
	}

	var nodePool []*MerkleNode
	for _, tx := range txs {
		nodePool = append(nodePool, CreateMerkleNode(nil, nil, tx.ID))
	}

	for len(nodePool) > 1 {
		var tempNodePool []*MerkleNode
		poolLen := len(nodePool)
		if poolLen%2 != 0 {
			tempNodePool = append(tempNodePool, nodePool[poolLen-1])
		}
		for i := 0; i < poolLen/2; i++ {
			tempNodePool = append(tempNodePool, CreateMerkleNode(nodePool[2*i], nodePool[2*i+1], nil))
		}
		nodePool = tempNodePool
	}

	merkleTree := MerkleTree{nodePool[0]}
	return &merkleTree
}

func (mn *MerkleNode) Find(data []byte, route []int, hashRoute [][]byte) (bool, []int, [][]byte) {
	findFlag := false

	if bytes.Equal(mn.Data, data) {
		findFlag = true
		return findFlag, route, hashRoute
	} else {
		if mn.LeftNode != nil {
			route_t := append(route, 0)
			hashRoute_t := append(hashRoute, mn.RightNode.Data)
			findFlag, route_t, hashRoute_t = mn.LeftNode.Find(data, route_t, hashRoute_t)
			if findFlag {
				return findFlag, route_t, hashRoute_t
			} else {
				if mn.RightNode.Data == nil {
					route_t = append(route, 1)
					hashRoute_t = append(hashRoute, mn.LeftNode.Data)
					findFlag, route_t, hashRoute_t = mn.Find(data, route_t, hashRoute_t)
					if findFlag {
						return findFlag, route_t, hashRoute_t
					} else {
						return findFlag, route, hashRoute
					}
				}
			}
		} else {
			return findFlag, route, hashRoute
		}
	}
	return findFlag, route, hashRoute
}

func (mn *MerkleTree) BackValidationRoute(txid []byte) ([]int, [][]byte, bool) {
	ok, route, hashRoute := mn.RootNode.Find(txid, []int{}, [][]byte{})
	return route, hashRoute, ok
}

func SimplePaymentValidation(txid, mtRootHash []byte, route []int, hashRoute [][]byte) bool {
	routeLen := len(route)
	var tempHash []byte
	tempHash = txid

	for i := routeLen - 1; i >= 0; i-- {
		if route[i] == 0 {
			catenateHash := append(tempHash, hashRoute[i]...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else if route[i] == 1 {
			catenateHash := append(hashRoute[i], tempHash...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else {
			utils.Handle(errors.New("error in validation route"))
		}
	}
	return bytes.Equal(tempHash, mtRootHash)
}
