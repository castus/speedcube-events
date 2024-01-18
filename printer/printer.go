package printer

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(obj any) {
	j, _ := json.MarshalIndent(obj, "", "    ")
	fmt.Printf("%s\n", j)
}
