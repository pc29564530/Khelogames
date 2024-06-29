package sports

import (
	"github.com/gin-gonic/gin"
)

type SportsService interface {
	Otp(ctx *gin.Context)
	CreateSignupFunc(ctx *gin.Context)
	CreateUserFunc(ctx *gin.Context)
	CreateLoginFunc(ctx *gin.Context)
	DeleteSessionFunc(ctx *gin.Context)
	RenewAccessTokenFunc(ctx *gin.Context)
	GetUsersFunc(ctx *gin.Context)
	GetProfileFunc(ctx *gin.Context)
	GetTournamentMatchFunc(ctx *gin.Context)
	GetClubFunc(ctx *gin.Context)
	GetFootballMatchesFunc(ctx *gin.Context)
	AddJoinCommunityFunc(ctx *gin.Context)
	GetUserByCommunityFunc(ctx *gin.Context)
	GetCommunityByUserFunc(ctx *gin.Context)
	ListUsersFunc(ctx *gin.Context)
	CreateCommunitesFunc(ctx *gin.Context)
	GetCommunityFunc(ctx *gin.Context)
	GetAllCommunitiesFunc(ctx *gin.Context)
	GetCommunityByCommunityNameFunc(ctx *gin.Context)
	CheckLikeByUserFunc(ctx *gin.Context)
	CreateThreadFunc(ctx *gin.Context)
	GetThreadFunc(ctx *gin.Context)
	UpdateThreadLikeFunc(ctx *gin.Context)
	GetAllThreadsFunc(ctx *gin.Context)
	GetAllThreadByCommunityFunc(ctx *gin.Context)
	GetCommunitiesMemberFunc(ctx *gin.Context)
	CreateFollowingFunc(ctx *gin.Context)
	GetAllFollowerFunc(ctx *gin.Context)
	GetAllFollowingFunc(ctx *gin.Context)
	CreateCommentFunc(ctx *gin.Context)
	GetAllCommentFunc(ctx *gin.Context)
	GetCommentByUserFunc(ctx *gin.Context)
	DeleteFollowingFunc(ctx *gin.Context)
	CreateLikeFunc(ctx *gin.Context)
	CountLikeFunc(ctx *gin.Context)
	CreateProfileFunc(ctx *gin.Context)
	UpdateProfileFunc(ctx *gin.Context)
	UpdateAvatarUrlFunc(ctx *gin.Context)
	UpdateCoverUrlFunc(ctx *gin.Context)
	UpdateFullNameFunc(ctx *gin.Context)
	UpdateBioFunc(ctx *gin.Context)
	GetThreadByUserFunc(ctx *gin.Context)
	GetMessageByReceiverFunc(ctx *gin.Context)
	UpdateClubSportFunc(ctx *gin.Context)
	AddClubMemberFunc(ctx *gin.Context)
	CreateTournamentFunc(ctx *gin.Context)
	GetPlayerProfileFunc(ctx *gin.Context)
	AddPlayerProfileFunc(ctx *gin.Context)
	GetUserByMessageSendFunc(ctx *gin.Context)
	GetAllPlayerProfileFunc(ctx *gin.Context)
	UpdatePlayerProfileAvatarUrlFunc(ctx *gin.Context)
	AddGroupTeamFunc(ctx *gin.Context)
	CreateTournamentOrganizationFunc(ctx *gin.Context)
	GetTournamentOrganizationFunc(ctx *gin.Context)
	CreateUploadMediaFunc(ctx *gin.Context)
	CreateMessageMediaFunc(ctx *gin.Context)
	CreateCommunityMessageFunc(ctx *gin.Context)
	GetCommunityMessageFunc(ctx *gin.Context)
	GetCommunityByMessageFunc(ctx *gin.Context)
	CreateOrganizerFunc(ctx *gin.Context)
	GetOrganizerFunc(ctx *gin.Context)
	CreateClubFunc(ctx *gin.Context)
	CreateTournamentMatchFunc(ctx *gin.Context)
	GetTeamsByGroupFunc(ctx *gin.Context)
	GetTeamsFunc(ctx *gin.Context)
	GetTournamentsBySportFunc(ctx *gin.Context)
	GetTournamentFunc(ctx *gin.Context)
	AddFootballMatchScoreFunc(ctx *gin.Context)
	GetFootballMatchScoreFunc(ctx *gin.Context)
	UpdateFootballMatchScoreFunc(ctx *gin.Context)
	AddFootballGoalByPlayerFunc(ctx *gin.Context)
	GetClubsFunc(ctx *gin.Context)
	GetClubMemberFunc(ctx *gin.Context)
	GetAllTournamentMatchFunc(ctx *gin.Context)
	AddCricketMatchScoreFunc(ctx *gin.Context)
	GetCricketMatchScoreFunc(ctx *gin.Context)
	UpdateCricketMatchRunsScoreFunc(ctx *gin.Context)
	UpdateCricketMatchWicketFunc(ctx *gin.Context)
	UpdateCricketMatchExtrasFunc(ctx *gin.Context)
	UpdateCricketMatchInningsFunc(ctx *gin.Context)
	AddCricketMatchTossFunc(ctx *gin.Context)
	GetCricketMatchTossFunc(ctx *gin.Context)
	AddCricketTeamPlayerScoreFunc(ctx *gin.Context)
	GetCricketTeamPlayerScoreFunc(ctx *gin.Context)
	GetCricketPlayerScoreFunc(ctx *gin.Context)
	UpdateCricketMatchScoreBattingFunc(ctx *gin.Context)
	UpdateCricketMatchScoreBowlingFunc(ctx *gin.Context)
	GetMatchByClubNameFunc(ctx *gin.Context)
	UpdateTournamentDateFunc(ctx *gin.Context)
	CreateTournamentStandingFunc(ctx *gin.Context)
	CreateTournamentGroupFunc(ctx *gin.Context)
	GetTournamentGroupFunc(ctx *gin.Context)
	GetTournamentGroupsFunc(ctx *gin.Context)
	GetTournamentStandingFunc(ctx *gin.Context)
	GetClubsBySportFunc(ctx *gin.Context)
	GetTournamentsByClubFunc(ctx *gin.Context)
	AddTeamFunc(ctx *gin.Context)
	GetMatchFunc(ctx *gin.Context)
	GetTournamentByLevelFunc(ctx *gin.Context)
	GetFootballTournamentMatchesFunc(ctx *gin.Context)
}

type SportsServiceImp struct{}

func (s *SportsServiceImp) Otp(ctx *gin.Context)                                {}
func (s *SportsServiceImp) CreateSignupFunc(ctx *gin.Context)                   {}
func (s *SportsServiceImp) CreateUserFunc(ctx *gin.Context)                     {}
func (s *SportsServiceImp) CreateLoginFunc(ctx *gin.Context)                    {}
func (s *SportsServiceImp) DeleteSessionFunc(ctx *gin.Context)                  {}
func (s *SportsServiceImp) RenewAccessTokenFunc(ctx *gin.Context)               {}
func (s *SportsServiceImp) GetUsersFunc(ctx *gin.Context)                       {}
func (s *SportsServiceImp) GetProfileFunc(ctx *gin.Context)                     {}
func (s *SportsServiceImp) GetTournamentMatchFunc(ctx *gin.Context)             {}
func (s *SportsServiceImp) GetClubFunc(ctx *gin.Context)                        {}
func (s *SportsServiceImp) GetFootballMatchesFunc(ctx *gin.Context)             {}
func (s *SportsServiceImp) AddJoinCommunityFunc(ctx *gin.Context)               {}
func (s *SportsServiceImp) GetUserByCommunityFunc(ctx *gin.Context)             {}
func (s *SportsServiceImp) GetCommunityByUserFunc(ctx *gin.Context)             {}
func (s *SportsServiceImp) ListUsersFunc(ctx *gin.Context)                      {}
func (s *SportsServiceImp) CreateCommunitesFunc(ctx *gin.Context)               {}
func (s *SportsServiceImp) GetCommunityFunc(ctx *gin.Context)                   {}
func (s *SportsServiceImp) GetAllCommunitiesFunc(ctx *gin.Context)              {}
func (s *SportsServiceImp) GetCommunityByCommunityNameFunc(ctx *gin.Context)    {}
func (s *SportsServiceImp) CheckLikeByUserFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) CreateThreadFunc(ctx *gin.Context)                   {}
func (s *SportsServiceImp) GetThreadFunc(ctx *gin.Context)                      {}
func (s *SportsServiceImp) UpdateThreadLikeFunc(ctx *gin.Context)               {}
func (s *SportsServiceImp) GetAllThreadsFunc(ctx *gin.Context)                  {}
func (s *SportsServiceImp) GetAllThreadByCommunityFunc(ctx *gin.Context)        {}
func (s *SportsServiceImp) GetCommunitiesMemberFunc(ctx *gin.Context)           {}
func (s *SportsServiceImp) CreateFollowingFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) GetAllFollowerFunc(ctx *gin.Context)                 {}
func (s *SportsServiceImp) GetAllFollowingFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) CreateCommentFunc(ctx *gin.Context)                  {}
func (s *SportsServiceImp) GetAllCommentFunc(ctx *gin.Context)                  {}
func (s *SportsServiceImp) GetCommentByUserFunc(ctx *gin.Context)               {}
func (s *SportsServiceImp) DeleteFollowingFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) CreateLikeFunc(ctx *gin.Context)                     {}
func (s *SportsServiceImp) CountLikeFunc(ctx *gin.Context)                      {}
func (s *SportsServiceImp) CreateProfileFunc(ctx *gin.Context)                  {}
func (s *SportsServiceImp) UpdateProfileFunc(ctx *gin.Context)                  {}
func (s *SportsServiceImp) UpdateAvatarUrlFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) UpdateCoverUrlFunc(ctx *gin.Context)                 {}
func (s *SportsServiceImp) UpdateFullNameFunc(ctx *gin.Context)                 {}
func (s *SportsServiceImp) UpdateBioFunc(ctx *gin.Context)                      {}
func (s *SportsServiceImp) GetThreadByUserFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) GetMessageByReceiverFunc(ctx *gin.Context)           {}
func (s *SportsServiceImp) UpdateClubSportFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) AddClubMemberFunc(ctx *gin.Context)                  {}
func (s *SportsServiceImp) CreateTournamentFunc(ctx *gin.Context)               {}
func (s *SportsServiceImp) GetPlayerProfileFunc(ctx *gin.Context)               {}
func (s *SportsServiceImp) AddPlayerProfileFunc(ctx *gin.Context)               {}
func (s *SportsServiceImp) GetUserByMessageSendFunc(ctx *gin.Context)           {}
func (s *SportsServiceImp) GetAllPlayerProfileFunc(ctx *gin.Context)            {}
func (s *SportsServiceImp) UpdatePlayerProfileAvatarUrlFunc(ctx *gin.Context)   {}
func (s *SportsServiceImp) AddGroupTeamFunc(ctx *gin.Context)                   {}
func (s *SportsServiceImp) CreateTournamentOrganizationFunc(ctx *gin.Context)   {}
func (s *SportsServiceImp) GetTournamentOrganizationFunc(ctx *gin.Context)      {}
func (s *SportsServiceImp) CreateUploadMediaFunc(ctx *gin.Context)              {}
func (s *SportsServiceImp) CreateMessageMediaFunc(ctx *gin.Context)             {}
func (s *SportsServiceImp) CreateCommunityMessageFunc(ctx *gin.Context)         {}
func (s *SportsServiceImp) GetCommunityMessageFunc(ctx *gin.Context)            {}
func (s *SportsServiceImp) GetCommunityByMessageFunc(ctx *gin.Context)          {}
func (s *SportsServiceImp) CreateOrganizerFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) GetOrganizerFunc(ctx *gin.Context)                   {}
func (s *SportsServiceImp) CreateClubFunc(ctx *gin.Context)                     {}
func (s *SportsServiceImp) CreateTournamentMatchFunc(ctx *gin.Context)          {}
func (s *SportsServiceImp) GetTeamsByGroupFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) GetTeamsFunc(ctx *gin.Context)                       {}
func (s *SportsServiceImp) GetTournamentsBySportFunc(ctx *gin.Context)          {}
func (s *SportsServiceImp) GetTournamentFunc(ctx *gin.Context)                  {}
func (s *SportsServiceImp) AddFootballMatchScoreFunc(ctx *gin.Context)          {}
func (s *SportsServiceImp) GetFootballMatchScoreFunc(ctx *gin.Context)          {}
func (s *SportsServiceImp) UpdateFootballMatchScoreFunc(ctx *gin.Context)       {}
func (s *SportsServiceImp) AddFootballGoalByPlayerFunc(ctx *gin.Context)        {}
func (s *SportsServiceImp) GetClubsFunc(ctx *gin.Context)                       {}
func (s *SportsServiceImp) GetClubMemberFunc(ctx *gin.Context)                  {}
func (s *SportsServiceImp) GetAllTournamentMatchFunc(ctx *gin.Context)          {}
func (s *SportsServiceImp) AddCricketMatchScoreFunc(ctx *gin.Context)           {}
func (s *SportsServiceImp) GetCricketMatchScoreFunc(ctx *gin.Context)           {}
func (s *SportsServiceImp) UpdateCricketMatchRunsScoreFunc(ctx *gin.Context)    {}
func (s *SportsServiceImp) UpdateCricketMatchWicketFunc(ctx *gin.Context)       {}
func (s *SportsServiceImp) UpdateCricketMatchExtrasFunc(ctx *gin.Context)       {}
func (s *SportsServiceImp) UpdateCricketMatchInningsFunc(ctx *gin.Context)      {}
func (s *SportsServiceImp) AddCricketMatchTossFunc(ctx *gin.Context)            {}
func (s *SportsServiceImp) GetCricketMatchTossFunc(ctx *gin.Context)            {}
func (s *SportsServiceImp) AddCricketTeamPlayerScoreFunc(ctx *gin.Context)      {}
func (s *SportsServiceImp) GetCricketTeamPlayerScoreFunc(ctx *gin.Context)      {}
func (s *SportsServiceImp) GetCricketPlayerScoreFunc(ctx *gin.Context)          {}
func (s *SportsServiceImp) UpdateCricketMatchScoreBattingFunc(ctx *gin.Context) {}
func (s *SportsServiceImp) UpdateCricketMatchScoreBowlingFunc(ctx *gin.Context) {}
func (s *SportsServiceImp) GetMatchByClubNameFunc(ctx *gin.Context)             {}
func (s *SportsServiceImp) UpdateTournamentDateFunc(ctx *gin.Context)           {}
func (s *SportsServiceImp) CreateTournamentStandingFunc(ctx *gin.Context)       {}
func (s *SportsServiceImp) CreateTournamentGroupFunc(ctx *gin.Context)          {}
func (s *SportsServiceImp) GetTournamentGroupFunc(ctx *gin.Context)             {}
func (s *SportsServiceImp) GetTournamentGroupsFunc(ctx *gin.Context)            {}
func (s *SportsServiceImp) GetTournamentStandingFunc(ctx *gin.Context)          {}
func (s *SportsServiceImp) GetClubsBySportFunc(ctx *gin.Context)                {}
func (s *SportsServiceImp) GetTournamentsByClubFunc(ctx *gin.Context)           {}
func (s *SportsServiceImp) AddTeamFunc(ctx *gin.Context)                        {}
func (s *SportsServiceImp) GetMatchFunc(ctx *gin.Context)                       {}
func (s *SportsServiceImp) GetTournamentByLevelFunc(ctx *gin.Context)           {}
func (s *SportsServiceImp) GetFootballTournamentMatchesFunc(ctx *gin.Context)   {}

// Implement all the methods of the SportsService interface...

// NewSportsService returns a new SportsService implementation
func NewSportsService(ctx *gin.Context) SportsService {

	return &SportsServiceImp{}
}
