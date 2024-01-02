package app

import "github.com/dpwgc/easierweb"

type MemberController struct{}

var memberService MemberService

// Add create
// [POST] http://localhost/api/v2/member
/*
body:
{
  "name": "hello",
  "mobile": "12345678",
  "weight": 10.24
}
*/
func (c *MemberController) Add(ctx *easierweb.Context, command MemberCommand) (*MemberDTO, error) {
	return memberService.Add(command)
}

// Del delete
// [DELETE] http://localhost/api/v2/member/1
func (c *MemberController) Del(ctx *easierweb.Context) {
	memberService.Del(ctx.Path.Int64("id"))
}

// Edit update
// [PUT] http://localhost/api/v2/member/1
/*
body:
{
  "name": "hello",
  "mobile": "87654321",
  "weight": 20.48
}
*/
func (c *MemberController) Edit(ctx *easierweb.Context, command MemberCommand) error {
	return memberService.Edit(ctx.Path.Int64("id"), command)
}

// Get select one
// [GET] http://localhost/api/v2/member/1
func (c *MemberController) Get(ctx *easierweb.Context) *MemberDTO {
	return memberService.Get(ctx.Path.Int64("id"))
}

// List select more
// [GET] http://localhost/api/v2/members?name=hello&mobile=12345678
func (c *MemberController) List(ctx *easierweb.Context, query MemberQuery) *[]MemberDTO {
	return memberService.List(query)
}
