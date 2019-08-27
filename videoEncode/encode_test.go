package videoencoder

import "testing"

func TestEncode(t *testing.T) {
	input := `C:\Users\nreal\Desktop\RecordRes\NetImages\RGB\`
	output := `C:\Users\nreal\Desktop\RecordRes\NetImages\output\`
	Do(input, output)
}
