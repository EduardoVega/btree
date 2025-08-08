# BTree

## Roadmap

[ ] Implement ReplaceOrInsert operation.
[ ] Implement Delete operation.
[ ] Implement Get operation.
[ ] Implement Range operation.
[ ] Implement HasKey operation.
[ ] Implement Print operation.
[ ] Read about copy on write.

```
To test logic 
10, 20, 30, 40, 50, 60, 5, 15, 25, 35, 45, 55, 65, 75, 85, 95, 105

{44}
{4}
{28}
{3}
{15}
{30}
{48}
{17}
{38}
{23}
{21}
{20}
{22}
{9}
{27}
{18}
{2}
{44}
{5}
{43}
{13}
{34}
{39}









Index 2, len of nodes 3, keys of node [{17} {28}], key to insert {39} 
-1:0:{17} -1:0:{28} 
0:0:{3} 0:0:{8} 0:1:{21} 0:2:{38} 0:2:{44} 
0:0:{2} 0:1:{4} 0:1:{5} 0:2:{9} 0:2:{13} 0:2:{15} 1:3:{18} 1:3:{20} 1:4:{27} 1:4:{30} 2:5:{34} 2:5:{43} 
------------------
Index 1, len of nodes 1, keys of node [{38} {44}], key to insert {39} 
-1:0:{17} -1:0:{28} 
0:0:{3} 0:0:{8} 0:1:{21} 0:2:{38} 0:2:{44} 
0:0:{2} 0:1:{4} 0:1:{5} 0:2:{9} 0:2:{13} 0:2:{15} 1:3:{18} 1:3:{20} 1:4:{27} 1:4:{30} 2:5:{34} 2:5:{43} 
------------------
       0	               NaN ns/op	       0 B/op	       0 allocs/op
panic: runtime error: index out of range [1] with length 1



{29}
{22}
{25}
{37}
{4}
{43}
{2}
{39}
{5}
{48}
{41}
{34}

{26}
{33}
{17}
{6}
{30}
{9}
{0}
{49}
```