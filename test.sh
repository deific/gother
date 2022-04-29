rm -rf ./tmp/*

mkdir ./tmp/blocks
./main create -address Steven
./main info
./main balance -address Steven
./main send -from Steven -to One -amount 100
./main balance -address One
./main mine
./main info
./main balance -address Steven
./main balance -address One
./main send -from Steven -to Two -amount 100
./main send -from One -to Two -amount 30
./main mine
./main info
./main balance -address Steven
./main balance -address One
./main balance -address Two