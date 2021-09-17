package installbase

import "testing"

func testDeploy(*StageContext) error {
	return nil
}

func TestDeployFunc(t *testing.T) {
	var fn InstallFunc = testDeploy
	fn.Deploy(nil)
}
