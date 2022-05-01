rm -rf ./tmp/*
mkdir ./tmp/blocks
mkdir ./tmp/wallets
mkdir ./tmp/ref_list
./main createwallet
./main walletslist
./main createwallet -refname Steven
./main walletinfo -refname Steven
./main createwallet -refname One
./main createwallet -refname Two
./main createwallet
./main walletslist
./main createblockchain -refname Steven
./main blockchaininfo
./main balance -refname Steven
./main sendbyrefname -from Steven -to One -amount 100
./main balance -refname One
./main mine
./main blockchaininfo
./main balance -refname Steven
./main balance -refname One
./main sendbyrefname -from Steven -to Two -amount 100
./main sendbyrefname -from One -to Two -amount 30
./main mine
./main blockchaininfo
./main balance -refname Steven
./main balance -refname One
./main balance -refname Two
./main sendbyrefname -from Two -to Steven -amount 90
./main sendbyrefname -from Two -to One -amount 90
./main mine
./main blockchaininfo
./main balance -refname Steven
./main balance -refname One
./main balance -refname Two