package app

import "github.com/dpwgc/easierweb"

var memberService MemberService

// AddMember create
// [POST] http://localhost/api/v2/member
/*
body:
{
  "name": "hello",
  "mobile": "12345678",
  "weight": 10.24
}
*/
func AddMember(ctx *easierweb.Context, command MemberCommand) (*MemberDTO, error) {
	return memberService.Add(command)
}

// DelMember delete
// [DELETE] http://localhost/api/v2/member/1
func DelMember(ctx *easierweb.Context) {
	memberService.Del(ctx.Path.Int64("id"))
}

// EditMember update
// [PUT] http://localhost/api/v2/member/1
/*
body:
{
  "name": "hello",
  "mobile": "87654321",
  "weight": 20.48
}
*/
func EditMember(ctx *easierweb.Context, command MemberCommand) {
	memberService.Edit(ctx.Path.Int64("id"), command)
}

// GetMember select one
// [GET] http://localhost/api/v2/member/1
func GetMember(ctx *easierweb.Context) *MemberDTO {
	return memberService.Get(ctx.Path.Int64("id"))
}

// ListMember select more
// [GET] http://localhost/api/v2/members?name=hello&mobile=12345678
func ListMember(ctx *easierweb.Context, query MemberQuery) *[]MemberDTO {
	return memberService.List(query)
}
