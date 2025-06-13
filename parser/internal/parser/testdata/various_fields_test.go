package main_test

import "testing"

func TestVariousFields(t *testing.T) {
	tests := []struct {
		description string
		title       string
		scenario    string
		testName    string
	}{
		{description: "説明フィールド"},
		{title: "タイトルフィールド"},
		{scenario: "シナリオフィールド"},
		{testName: "テスト名フィールド"},
	}
}
