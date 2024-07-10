package main

import (
	"encoding/json"
	"log"
	"os"
)

type Poem struct {
	Title  string
	Author string
	Text   string
	Length int
}

var m = Poem{
	"Test Poem",
	"Somebody",
	"abcdefg",
	7,
}

var m2 Poem

var err error

func main() {
	m_bytes, err := json.Marshal(m)
  println(m.Title)
  println(m.Author)
  println(m.Text)
  println(m.Length)
	println(string(m_bytes))
	println(err)

	err = os.WriteFile("file.json", m_bytes, 0666)
	if err != nil {
		log.Fatal(err)
	}

  file_bytes, err := os.ReadFile("file.json")
	if err != nil {
		log.Fatal(err)
	}
  println(string(file_bytes))
  println(err)

  json.Unmarshal(file_bytes, &m2)
  println(m2.Title)
  println(m2.Author)
  println(m2.Text)
  println(m2.Length)

}
