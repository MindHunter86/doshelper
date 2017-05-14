package main


type testStruct struct {
	count int
}

func (self * testStruct) increase() int { return self.count + 1 }
