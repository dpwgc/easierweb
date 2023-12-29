package app

import "github.com/dpwgc/easierweb"

var memberService MemberService

// AddMember create
func AddMember(ctx *easierweb.Context, command MemberCommand) (*MemberDTO, error) {
	return memberService.AddMember(command)
}

// DelMember delete
func DelMember(ctx *easierweb.Context) {
	memberService.DelMember(ctx.Path.Int64("id"))
}

// EditMember update
func EditMember(ctx *easierweb.Context, command MemberCommand) {
	memberService.EditMember(ctx.Path.Int64("id"), command)
}

// GetMember select one
func GetMember(ctx *easierweb.Context) *MemberDTO {
	return memberService.GetMember(ctx.Path.Int64("id"))
}

// ListMember select more
func ListMember(ctx *easierweb.Context, query MemberQuery) *[]MemberDTO {
	return memberService.ListMember(query)
}
