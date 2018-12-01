package openbsdvmm

import (
	"fmt"
	"log"
	"os"
)

type VmmArtifact struct {
	imageName string
	imageID   int
	imageSize int
}

func (*VmmArtifact) BuilderId() string {
	return BuilderID
}

// who needs this?
func (a *VmmArtifact) Files() []string {
	return append([]string{a.imageName}, "nuts\n")
}

func (a *VmmArtifact) Id() string {
	return fmt.Sprintf("%d", a.imageID)
}

func (a *VmmArtifact) String() string {
	return fmt.Sprintf("%s (Size %d)\n", a.imageName, a.imageSize)
}

func (a *VmmArtifact) State(name string) interface{} {
	return nil
}

func (a VmmArtifact) Destroy() error {
	log.Printf("Deleting %s", a.imageName)
	return os.Remove(a.imageName)
}
