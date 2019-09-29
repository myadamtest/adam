package grpcservice

import (
	"context"
	"github.com/myadamtest/adam/grpcservice/pb/adam"
)

type SocialRelationsServiceImpl struct{}

func (this *SocialRelationsServiceImpl) Subscribe(st adam.SocialRelationsService_SubscribeServer) error {
	panic("to implement")
	return nil
}

func (this *SocialRelationsServiceImpl) CreatePlayerCard(ctx context.Context, in *adam.CreatePlayerRequest) (*adam.CreatePlayerResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetPlayerCardInfo(ctx context.Context, in *adam.GetPlayerCardInfoRequest) (*adam.GetPlayerCardInfoResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) DeletePlayerCard(ctx context.Context, in *adam.DeletePlayerCardRequest) (*adam.DeletePlayerCardResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) ModifyPlayerCardInfo(ctx context.Context, in *adam.ModifyPlayerCardInfoRequest) (*adam.ModifyPlayerCardInfoResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) AddPlayerToRoster(ctx context.Context, in *adam.AddPlayerToRosterRequest) (*adam.AddPlayerToRosterResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) DeletePlayerFromRoster(ctx context.Context, in *adam.DeletePlayerFromRosterRequest) (*adam.DeletePlayerFromRosterResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetPlayerRoster(ctx context.Context, in *adam.GetPlayerRosterRequest) (*adam.GetPlayerRosterResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetPlayerWithTargetRoster(ctx context.Context, in *adam.GetPlayerWithTargetRosterRequest) (*adam.GetPlayerWithTargetRosterResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) ModifyPlayerRosterInfo(ctx context.Context, in *adam.ModifyPlayerRosterRequest) (*adam.ModifyPlayerRosterResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetPlayerLiked(ctx context.Context, in *adam.GetPlayerLikedRequest) (*adam.GetPlayerLikedResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) CreateOrganization(ctx context.Context, in *adam.CreateOrganizationRequest) (*adam.CreateOrganizationResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) ModifyOrgnization(ctx context.Context, in *adam.ModifyOrgnizationRequest) (*adam.ModifyOrgnizationResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) DeleteOrganization(ctx context.Context, in *adam.DeleteOrganizationRequest) (*adam.DeleteOrganizationResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetOrganizationInfo(ctx context.Context, in *adam.GetOrganizationInfoRequest) (*adam.GetOrganizationInfoResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetOrganizationMembers(ctx context.Context, in *adam.GetOrganizationMembersRequest) (*adam.GetOrganizationMembersResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) JoinOrganization(ctx context.Context, in *adam.JoinOrganizationRequest) (*adam.JoinOrganizationResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) InviteToJoinOrganization(ctx context.Context, in *adam.InviteToJoinOrganizationRequest) (*adam.InviteToJoinOrganizationResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) ModifyOrgnizationMember(ctx context.Context, in *adam.ModifyOrgnizationMemberRequest) (*adam.ModifyOrgnizationMemberResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) DeleteOrgnizationMember(ctx context.Context, in *adam.DeleteOrgnizationMemberRequest) (*adam.DeleteOrgnizationMemberResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetOrganizationApplyList(ctx context.Context, in *adam.GetOrganizationApplyListRequest) (*adam.GetOrganizationApplyListResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) AcceptOrganizationApply(ctx context.Context, in *adam.AcceptOrganizationApplyRequest) (*adam.AcceptOrganizationApplyResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetPlayerOrganizations(ctx context.Context, in *adam.GetPlayerOrganizationsRequest) (*adam.GetPlayerOrganizationsResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) EditOrgNotice(ctx context.Context, in *adam.EditOrgNoticeRequest) (*adam.EditOrgNoticeResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GainExp(ctx context.Context, in *adam.GainExpRequest) (*adam.GainExpResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) OpenOrgActivity(ctx context.Context, in *adam.OpenOrgActivityRequest) (*adam.OpenOrgActivityResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) ChangeCommunicateLanguage(ctx context.Context, in *adam.ChangeCommunicateLanguageRequest) (*adam.ChangeCommunicateLanguageResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) DisbandOrganization(ctx context.Context, in *adam.DisbandOrganizationRequest) (*adam.DisbandOrganizationResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) CreateSceneTeam(ctx context.Context, in *adam.CreateSceneTeamRequest) (*adam.CreateSceneTeamResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) AddLatelyPerson(ctx context.Context, in *adam.AddPlayerLatelyRequest) (*adam.AddPlayerLatelyResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetLatelyPerson(ctx context.Context, in *adam.GetPlayerLatelyRequest) (*adam.GetPlayerLatelyResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) AddItemReceiveLog(ctx context.Context, in *adam.AddItemReceiveLogRequest) (*adam.AddItemReceiveLogResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) GetItemReceiveLogList(ctx context.Context, in *adam.GetItemReceiveLogListRequest) (*adam.GetItemReceiveLogListResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) DeleteItemReceiveLog(ctx context.Context, in *adam.DeleteItemReceiveLogRequest) (*adam.DeleteItemReceiveLogResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) DeleteAllItemReceiveLog(ctx context.Context, in *adam.DeleteAllItemReceiveLogRequest) (*adam.DeleteAllItemReceiveLogResponse, error) {
	panic("to implement")
	return nil, nil
}

func (this *SocialRelationsServiceImpl) AcceptToJoinOrganization(ctx context.Context, in *adam.AcceptToJoinOrganizationRequest) (*adam.AcceptToJoinOrganizationResponse, error) {
	panic("to implement")
	return nil, nil
}
