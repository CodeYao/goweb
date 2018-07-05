package databaseop

import (
	"testing"
)

func Test_insertData_1(t *testing.T) {
	args := []string{"'001'", "'123456'", "'99'"}
	insertData("account", args)
}

func Test_queryDataById(t *testing.T) {
	println("***")
	type args struct {
		tableName string
		id        string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "account", args: args{"account", "4"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryDataById(tt.args.tableName, tt.args.id)
		})
	}
}
