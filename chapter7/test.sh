rm -rf ./tmp/test/*
mkdir ./tmp/test/blocks
mkdir ./tmp/test/wallets
mkdir ./tmp/test/ref_list
mkdir ./tmp/test/utxo
./gother createwallet
./gother walletslist
./gother createwallet -refname Steven
./gother walletinfo -refname Steven
./gother createwallet -refname One
./gother createwallet -refname Two
./gother createwallet
./gother walletslist
./gother createblockchain -refname Steven # 1000
./gother blockchaininfo
./gother balance -refname Steven # 1000
./gother sendbyrefname -from Steven -to One -amount 100
./gother balance -refname One
./gother mine
./gother blockchaininfo
./gother balance -refname Steven # Steven: 900 One:100
./gother balance -refname One    # Steven: 900 One:100
./gother sendbyrefname -from Steven -to Two -amount 100
./gother sendbyrefname -from One -to Two -amount 30
./gother mine # Steven: 800 Two:130 One:70
./gother blockchaininfo
./gother balance -refname Steven # Steven: 800 Two:130 One:70
./gother balance -refname One    # Steven: 800 Two:130 One:70
./gother balance -refname Two    # Steven: 800 Two:130 One:70
./gother sendbyrefname -from Two -to Steven -amount 90
./gother mine
./gother blockchaininfo
./gother balance -refname Steven # Steven: 890 Two:40 One:70
./gother balance -refname One # Steven: 890 Two:40 One:70
./gother balance -refname Two # Steven: 890 Two:40 One:70

./gother sendbyrefname -from Two -to One -amount 90
./gother mine

./gother balance2 -refname Steven # Steven: 890 Two:40 One:70
./gother balance2 -refname One # Steven: 890 Two:40 One:70
./gother balance2 -refname Two # Steven: 890 Two:40 One:70