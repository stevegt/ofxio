package ofxio

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/stevegt/goadapt"
)

func TestIO(t *testing.T) {

	ifh, err := os.Open("testdata/in.ofx")
	Tassert(t, err == nil, "%s", err)
	doc, err := Import(ifh)
	Tassert(t, err == nil, "%s", err)
	_ = doc

	ofh, err := os.Create("testdata/out.ofx")
	Tassert(t, err == nil, "%s", err)
	err = doc.Export(ofh)
	Tassert(t, err == nil, "%s", err)
	_ = doc

	rbuf, err := ioutil.ReadFile("testdata/ref.ofx")
	Tassert(t, err == nil, "%s", err)

	obuf, err := ioutil.ReadFile("testdata/out.ofx")
	Tassert(t, err == nil, "%s", err)

	Tassert(t, bytes.Equal(rbuf, obuf), "output differs from ref")

}
