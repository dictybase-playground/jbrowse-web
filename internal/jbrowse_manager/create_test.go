package jbrowse_manager

import (
	"context"
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	if err := Create(context.Background()); err != nil {
		fmt.Println(err)
	}
}
