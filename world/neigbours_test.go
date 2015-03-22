package world

import (
	"CaribbeanWarServer/commonStructs"
	"encoding/json"
	"testing"
)

func TestGetDifference(t *testing.T) {
	test1 := commonStructs.NearestUsers{
		commonStructs.NearestUser{ID: 0},
		commonStructs.NearestUser{ID: 1},
	}
	test2 := commonStructs.NearestUsers{
		commonStructs.NearestUser{ID: 1},
	}
	chanel := make(chan commonStructs.NearestUsers)

	go getDifference(&test1, &test2, &chanel)
	difference := <-chanel
	if len(difference) != 1 || (len(difference) == 1 && difference[0].ID != 0) {
		marshal, _ := json.Marshal(difference)
		t.Error("Incorrect! Len = ", len(difference), string(marshal))
	}

	go getDifference(&test2, &test1, &chanel)
	difference = <-chanel
	if len(difference) != 0 {
		marshal, _ := json.Marshal(difference)
		t.Error("Incorrect! Len = ", len(difference), string(marshal))
	}

	go getDifference(&test1, &test1, &chanel)
	difference = <-chanel
	if len(difference) != 0 {
		marshal, _ := json.Marshal(difference)
		t.Error("Incorrect! Len = ", len(difference), string(marshal))
	}

	go getDifference(&test1, &commonStructs.NearestUsers{}, &chanel)
	difference = <-chanel
	if len(difference) != len(test1) {
		marshal, _ := json.Marshal(difference)
		t.Error("Incorrect! Len = ", len(difference), string(marshal))
	}

	go getDifference(&commonStructs.NearestUsers{}, &test1, &chanel)
	difference = <-chanel
	if len(difference) != 0 {
		marshal, _ := json.Marshal(difference)
		t.Error("Incorrect! Len = ", len(difference), string(marshal))
	}

	close(chanel)
}
