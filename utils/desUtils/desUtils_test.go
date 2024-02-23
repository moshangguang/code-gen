package desUtils

import "testing"

func TestAes(t *testing.T) {
	srcData := "hello world !"
	//测试加密
	encData, err := ECBEncrypt([]byte(srcData), DefaultKey)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	//测试解密
	decData, err := ECBDecrypt(encData, DefaultKey)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(string(decData))
}
func TestDefaultKey(t *testing.T) {
	t.Log(len(DefaultKey))
}
