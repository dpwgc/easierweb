package app

import (
	"fmt"
	"time"
)

type MemberService struct{}

// AddMember create
func (s *MemberService) AddMember(command MemberCommand) (*MemberDTO, error) {
	fmt.Printf("[AddMember] name: %s, mobile: %s, weight: %v \n", command.Name, command.Mobile, command.Weight)
	time.Sleep(100 * time.Millisecond)
	return &MemberDTO{
		ID:     1,
		Name:   command.Name,
		Weight: command.Weight,
	}, nil
}

// DelMember delete
func (s *MemberService) DelMember(id int64) {
	fmt.Printf("[DelMember] id: %v \n", id)
	time.Sleep(100 * time.Millisecond)
}

// EditMember update
func (s *MemberService) EditMember(id int64, command MemberCommand) {
	fmt.Printf("[EditMember] id: %v, name: %s, mobile: %s, weight: %v \n", id, command.Name, command.Mobile, command.Weight)
	time.Sleep(100 * time.Millisecond)
}

// GetMember select one
func (s *MemberService) GetMember(id int64) *MemberDTO {
	fmt.Printf("[GetMember] id: %v \n", id)
	time.Sleep(100 * time.Millisecond)
	return &MemberDTO{
		ID:     1,
		Name:   "hello",
		Weight: 10.24,
	}
}

// ListMember select more
func (s *MemberService) ListMember(query MemberQuery) *[]MemberDTO {
	fmt.Printf("[ListMember] name: %s, mobile: %v \n", query.Name, query.Mobile)
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
