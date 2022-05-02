rm -rf ./tmp/*
mkdir ./tmp/blocks
mkdir ./tmp/wallets
mkdir ./tmp/ref_list
./gother createwallet
./gother walletslist
./gother createwallet -refname Steven
./gother walletinfo -refname Steven
./gother createwallet -refname One
./gother createwallet -refname Two
./gother createwallet
./gother walletslist
./gother createblockchain -refname Steven
./gother blockchaininfo
./gother balance -refname Steven
./gother sendbyrefname -from Steven -to One -amount 100
./gother balance -refname One
./gother mine
./gother blockchaininfo
./gother balance -refname Steven
./gother balance -refname One
./gother sendbyrefname -from Steven -to Two -amount 100
./gother sendbyrefname -from One -to Two -amount 30
./gother mine
./gother blockchaininfo
./gother balance -refname Steven
./gother balance -refname One
./gother balance -refname Two
./gother sendbyrefname -from Two -to Steven -amount 90
./gother sendbyrefname -from Two -to One -amount 90
./gother mine
./gother blockchaininfo
./gother balance -refname Steven
./gother balance -refname One
./gother balance -refname Two