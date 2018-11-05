package autoEntity

import (
	"fmt"
	"strings"
	"testing"
)

func TestUnderLine2Camel(t *testing.T) {
	cases := make(map[string][]string)
	cases["hello_world"] = []string{"helloWorld", "HelloWorld"}
	cases["es_order_detail"] = []string{"esOrderDetail", "EsOrderDetail"}
	cases["this_is_wield_w33_"] = []string{"thisIsWieldW33", "ThisIsWieldW33"}
	for key, value := range cases {
		if res := underLine2Camel(key, false); res != value[0] {
			t.Errorf("%s should change to %s but result is %s \n", key, value[0], res)
		}
		if res := underLine2Camel(key, true); res != value[1] {
			t.Errorf("%s should change to %s but result is %s \n", key, value[1], res)
		}
	}

}

func TestGenerateSeq(t *testing.T) {
	cases := "lyy_equipment, hello"
	names := strings.Split(cases, ",")
	result := GenerateSeq(names)
	fmt.Println(result)
}
