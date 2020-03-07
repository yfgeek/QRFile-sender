package core

import (
	"fmt"
	"log"
)

type SplitProtocol struct{
	id int
	content string
	Encoder
}

func (s *SplitProtocol) String() string{
	result, err := s.Encode([]byte(s.content))
	if err !=nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d|%s", s.id, result)
}



