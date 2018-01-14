package base

import "sync"

// SomeMap is a typed map that contains a reference to test data held in a Some struct.
var CopyLock sync.Mutex

type SomeMap map[string]*Some
