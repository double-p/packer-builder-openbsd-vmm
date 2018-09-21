package openbsd-vmm

import "fmt"

type Artifact struct {
	imageName string
	imageID   int
}

func (Artifact) BuilderId() string {
	return BuilderID
}

func (Artifact) Files() []string {
	return nil
}

func (a Artifact) Id() string {
	return fmt.Sprintf("%d", a.imageID)
}

func (a Artifact) String() string {
	return fmt.Sprintf("%d (%s)", a.imageID, a.imageName)
}

func (Artifact) State(name string) interface{} {
	return nil
}

func (a Artifact) Destroy() error {
	return nil
}
