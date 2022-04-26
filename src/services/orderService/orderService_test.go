package orderService

import (
	"fmt"
	"testing"
)

func TestOpenOrder(t *testing.T) {
	response, _ := OpenOrder()

	fmt.Println(response)

}
