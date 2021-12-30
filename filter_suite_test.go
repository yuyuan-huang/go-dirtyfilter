package filter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yuyuan-huang/go-dirtyfilter/v2"
	"fmt"
)

func TestDirtyFilterMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filter Suite Test.")
}

func TestDirtyReplaceMain(t *testing.T) {
	filterText := `共产党泛指以马克思主义为指导以建立共产主义社会为目标的工人党。其中陈@@@水@@@扁。在。。`
	nodeFilter := filter.NewNodeFilter([]string{"共产主义", "陈水扁"})
	data, err := nodeFilter.Replace(filterText, '*')
	if err != nil {
		Fail(err.Error())
		return
	}
	fmt.Println("...", data)
}
