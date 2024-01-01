package app

import (
	"fmt"
	"time"
)

type MemberService struct{}

// Add create
func (s *MemberService) Add(command MemberCommand) (*MemberDTO, error) {
	fmt.Printf("[Add] name: %s, mobile: %s, weight: %v \n", command.Name, command.Mobile, command.Weight)
	time.Sleep(100 * time.Millisecond)
	return &MemberDTO{
		ID:     1,
		Name:   command.Name,
		Weight: command.Weight,
	}, nil
}

// Del delete
func (s *MemberService) Del(id int64) {
	fmt.Printf("[Del] id: %v \n", id)
	time.Sleep(100 * time.Millisecond)
}

// Edit update
func (s *MemberService) Edit(id int64, command MemberCommand) {
	fmt.Printf("[Edit] id: %v, name: %s, mobile: %s, weight: %v \n", id, command.Name, command.Mobile, command.Weight)
	time.Sleep(100 * time.Millisecond)
}

// Get select one
func (s *MemberService) Get(id int64) *MemberDTO {
	fmt.Printf("[Get] id: %v \n", id)
	time.Sleep(100 * time.Millisecond)
	return &MemberDTO{
		ID:     1,
		Name:   "hello",
		Weight: 10.24,
	}
}

// List select more
func (s *MemberService) List(query MemberQuery) *[]MemberDTO {
	fmt.Printf("[List] name: %s, mobile: %v \n", query.Name, query.Mobile)
	time.Sleep(100 * time.Millisecond)
	var list []MemberDTO
	list = append(list, MemberDTO{
		ID:     1,
		Name:   "hello",
		Weight: 10.24,
	}, MemberDTO{
		ID:     2,
		Name:   "world",
		Weight: 20.48,
	})
	return &list
}
