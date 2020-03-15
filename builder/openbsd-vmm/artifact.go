package openbsdvmm

import (
	"fmt"
	"log"
)

type VmmArtifact struct {
	imageName []string
	imageDir  string
}

func (*VmmArtifact) BuilderId() string {
	return BuilderID
}

// who needs this?
func (a *VmmArtifact) Files() []string {
	return a.imageName
}

func (a *VmmArtifact) Id() string {
	return "VMM"
}

func (a *VmmArtifact) String() string {
	return fmt.Sprintf("Image: %s\n", a.imageName)
}

func (a *VmmArtifact) State(name string) interface{} {
	return nil
}

func (a VmmArtifact) Destroy() error {
	log.Printf("Deleting %s", a.imageName)
	return nil
}
