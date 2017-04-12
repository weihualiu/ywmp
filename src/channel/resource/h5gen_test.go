
package resource

import "testing"

func TestGenH5Desc(t *testing.T) {
	err := GenH5Desc("/home/liuweihua/workspace/ywmp/tests/bmosoa0104.zip")
	if err != nil {
		t.Errorf(`GenH5Desc() is failed!`, err)
	}
}