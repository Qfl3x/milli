package rope

import (
	"fmt"
)

const PARTITIONSIZE int = 5

var emptyString string = ""

type Node struct {
	left *Node
	right *Node
	weight int
	value *string
}

type Rope struct {
	root Node
}

func mergeHead(stack *[]Node) {
	stackVal := *stack
	left := stackVal[1]
	right := stackVal[2]
	stackVal = stackVal[3:]
	weight := nodeWeight(left)
	newNode := Node{left: &left, right: &right, weight: weight}
	stackVal = append([]Node{newNode}, stackVal...)
}

func NewRope(s string)(Rope)  {
	root := Node{weight: len(s), value: &s}
	return Rope{root}
}


func recursiveLeaves(node Node) ([]Node) {
	if node.left == nil {
		return []Node{node}	
	}
	return append(recursiveLeaves(*node.left), recursiveLeaves(*node.right)...)
}

func recursiveWeights(node Node, w *int) {
	if node.left == nil && node.right == nil {
		*w += node.weight
		return
	}
	if node.left != nil {
		recursiveWeights(*node.left, w)
	}
	if node.right != nil {
		recursiveWeights(*node.right, w)
	}
}

func nodeWeight(node Node) (int) {
	w := node.weight
	if node.right != nil {
		recursiveWeights(*node.right, &w)
	}
	return w
}


func recursiveReWeigh(node *Node) (int)  {
	if node.left == nil && node.right == nil {
		node.weight = len(*node.value)
		return node.weight
	}
	w := 0
	if node.left != nil {
		if node.left.value != nil && len(*node.left.value) == 0 {
			node.left = nil
			node.weight = w
		} else {
			w += recursiveReWeigh(node.left)
			node.weight = w
		}
	}
	if node.right != nil {
		if node.right.value != nil && len(*node.right.value) == 0 {
			node.right = nil
		} else {
			w += recursiveReWeigh(node.right)
		}
	}
	return w
}

func (r *Rope) reWeigh()  {
	recursiveReWeigh(&r.root)
}

func (r *Rope) Length() (int) {
	return nodeWeight(r.root)
}


func traverseRune(i int, n Node) (rune) {
	if i < n.weight && n.left != nil {
		return traverseRune(i, *n.left)
	}
	if n.right != nil {
		return traverseRune(i - n.weight, *n.right)
	}
	s := *n.value
	return rune(s[i])
}

func (r *Rope) Getindex(i int) (error, rune) {
	if i >= r.Length() {
		return fmt.Errorf("Index out of range"), rune('e')
	}
	return nil, traverseRune(i, r.root)
}

func (r *Rope) GetRange(lo int, hi int) (error, string) {
	if hi >= r.Length() || lo < 0 {
		return fmt.Errorf("Index out of range"), ""
	}
	s := ""
	for i := lo; i < hi + 1; i++ {
		_, rune := r.Getindex(i)
		s += string(rune)
	}
	return nil, s
}

func (r *Rope) String() (error, string) {
	s := ""
	for i := 0; i < r.Length(); i++ {
		err, c := r.Getindex(i)
		if err != nil {
			return err, ""
		}
		s += string(c)
	}
	return nil, s
}

func (r *Rope) Concat(s string)  {
	sRope := NewRope(s)
	oldRoot := r.root
	newRoot := Node{left: &oldRoot, right: &sRope.root, weight: r.Length() }
	*r = Rope{newRoot}
}

func traverseNode(i int, n *Node) (*Node, int) {
	if i < n.weight && n.left != nil {
		return traverseNode(i, n.left)
	}
	if n.right != nil {
		return traverseNode(i - n.weight, n.right)
	}
	return n, i
}


func (r *Rope) Insert(i int, s string)  {
	if i >= r.Length() {
		r.Concat(s)
		return 
	}
	ancestorNode, relInd := traverseNode(i, &r.root)
	leftString := *ancestorNode.value
	if relInd < len(*ancestorNode.value){
		v := *ancestorNode.value
		s = s +  v[relInd:] 
		leftString = v[:relInd]
	}
	newLeft := Node{weight: len(leftString), value: &leftString}
	newRight := Node{weight: len(s), value: &s}
	ancestorNode.value = nil
	ancestorNode.left = &newLeft
	ancestorNode.weight = newLeft.weight
	ancestorNode.right = &newRight
	r.reWeigh()
}

func (r *Rope) nodesToDelete(ind int, l int) ([]*Node, int, int) {
	firstNode, loRelInd := traverseNode(ind, &r.root)
	nodes := []*Node{firstNode}
	for i := ind + 1; i < ind + l; i++ {
		node, _ := traverseNode(i, &r.root)
		if node == nodes[len(nodes) -1] {
			continue
		} else {
			nodes = append(nodes, node)
		}
	}
	lastNode, hiRelInd := traverseNode(ind + l, &r.root)
	if lastNode != nodes[len(nodes) -1] {
		nodes = append(nodes, lastNode)
	}
	return nodes, loRelInd, hiRelInd 
}

func (r *Rope) deleteOne(i int)  {
	ancestorNode, relInd := traverseNode(i, &r.root)
	ancestorString := *ancestorNode.value
	if relInd == 0 {
		*ancestorNode.value = ancestorString[1:]
		ancestorNode.weight -= 1
	} else if relInd == len(ancestorString) - 1 {
		*ancestorNode.value = ancestorString[:relInd]
		ancestorNode.weight -= 1
	} else {
		leftString := ancestorString[:relInd]
		rightString := ancestorString[relInd+1:]
		newLeft := Node{weight: len(leftString), value: &leftString}
		newRight := Node{weight: len(rightString), value: &rightString}
		ancestorNode.value = nil
		ancestorNode.left = &newLeft
		ancestorNode.weight = newLeft.weight
		ancestorNode.right = &newRight
	}
	r.reWeigh()
	return 
}

func deleteFromNode(node *Node, loRelInd int, hiRelInd int)  {
	fmt.Println(loRelInd, hiRelInd, *node.value)
	if loRelInd == 0 {
		if hiRelInd == len(*node.value) - 1 {
			node.value = &emptyString
			node.weight = 0
		} else {
			nodeString := *node.value
			newString := nodeString[hiRelInd + 1:]
			node.weight = len(newString)
			node.value = &newString
		}
	} else {
		if hiRelInd == len(*node.value) - 1 {
			nodeString := *node.value
			newString := nodeString[:loRelInd]
			node.weight = len(newString)
			node.value = &newString
		} else {
			nodeString := *node.value
			leftString := nodeString[:loRelInd]
			rightString := nodeString[hiRelInd + 1:]
			newLeft := Node{weight: len(leftString), value: &leftString}
			newRight := Node{weight: len(rightString), value: &rightString}
			node.value = nil
			node.left = &newLeft
			node.weight = newLeft.weight
			node.right = &newRight
		}
	}
}

func (r *Rope) deleteMany(ind int, l int)  {
	nodes, loRelInd, hiRelInd := r.nodesToDelete(ind, l)
	// First Node Operations
	firstNode := nodes[0]
	if len(nodes) == 1 {
		deleteFromNode(firstNode, loRelInd, hiRelInd)
	} else {
		// First Node
		deleteFromNode(firstNode, loRelInd, firstNode.weight - 1)
		// Intermediate Nodes
		for i := 1; i < len(nodes) - 1; i++ {
			node := nodes[i]
			node.value = &emptyString
			node.weight = 0
		}
		// Last Node
		lastNode := nodes[len(nodes) - 1]
		deleteFromNode(lastNode, 0, hiRelInd)
	}
	r.reWeigh()
	return
}

func (r *Rope) Delete(i int, l int) (error) {
	if i < 0 || i + l >= r.Length() {
		return fmt.Errorf("Index out of range")
	}
	if l == 0 {
		r.deleteOne(i)
	} else {
		r.deleteMany(i, l)
	}
	return nil
}

func main()  {
	s := "abdefg"
	r := NewRope(s)
	r.deleteMany(0, 1)
	fmt.Println(r.String())
}
