# LIRS2
Repository for final project :  
Simulation and Analysis of the Cache Replacement Policy Algorithm LIRS2 (Low Inter-reference Recency Set 2)

## Paper
https://dl.acm.org/doi/10.1145/3456727.3463772

## How to run
1. Get module
   ```
   go get github.com/petar/GoLLRB
   go get github.com/secnot/orderedmap
   go get github.com/tidwall/btree
   ```
3. Set the cache size in the ***cachelist*** file. e.g:
   ```
   1000
   20000
   500000
   ```
4. Run program
   ```
   go run ./main.go [algoritma(LRU|LIRS|LIRS2)] [path-to-dataset] ./cachelist
   ```
