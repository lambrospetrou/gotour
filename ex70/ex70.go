package main

import (
 	"fmt"
    "code.google.com/p/go-tour/tree"
)

// recursive function for Walk
func WalkRec(t *tree.Tree, ch chan int) {
 	 if t == nil {
     	return    
    }
    WalkRec(t.Left, ch)
    ch <- t.Value
    WalkRec(t.Right, ch)  
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
    WalkRec(t, ch)
    close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
    ch1 := make(chan int)
    ch2 := make(chan int)
    go Walk(t1, ch1)
    go Walk(t2, ch2)
    // exhaust all t1 nodes
    for nv1 := range ch1{
        nv2, fin2 := <-ch2
        if !fin2 {
         	return false   
        }
        if nv1 != nv2 {
         	return false   
        }
    }
    // make sure there are no more nodes in t2
    _, fin2 := <-ch2
    if fin2 {
     	return false   
    }
    return true
}

func main() {

    ch := make(chan int)
    // test Walk function - question 3
    go Walk(tree.New(1), ch)
    for x:= range ch{
     	fmt.Println(x)   
    }
    // test the Same function - question 4
    fmt.Println(Same(tree.New(1), tree.New(1)))
    fmt.Println(Same(tree.New(1), tree.New(2)))
}
