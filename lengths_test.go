package main

import "testing"

func TestGetLengths(t *testing.T) {
  downloadPoems("poems.json")
  poems, _ := readPoemDB("poems.json")
  splitPoems(poems, "poem_folder")
  _, err := getLengths("poem_folder")
  if err != nil {
    t.Fail()
  }
}
