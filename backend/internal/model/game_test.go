package model

import (
	"testing"

	"github.com/HackYourCareer/SmartKickers/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestResetScore(t *testing.T) {
	game := &game{score: GameScore{3, 1}, scoreChannel: make(chan GameScore, 32)}

	game.ResetScore()
	resultScore := <-game.scoreChannel

	if resultScore.BlueScore != 0 || resultScore.WhiteScore != 0 {
		t.Errorf("Score did not reset. Goals white: %v, Goals blue: %v", game.score.WhiteScore, game.score.BlueScore)
	}
}

func TestAddGoal(t *testing.T) {
	game := &game{score: GameScore{3, 1}, scoreChannel: make(chan GameScore, 32)}

	type args struct {
		name               string
		teamID             int
		expectedBlueScore  int
		expectedWhiteScore int
		expectedError      string
	}
	tests := []args{
		{
			name:               "should increment team white score by one",
			teamID:             config.TeamWhite,
			expectedBlueScore:  0,
			expectedWhiteScore: 1, expectedError: "",
		},
		{
			name:               "should increment team blue score by one",
			teamID:             config.TeamBlue,
			expectedBlueScore:  1,
			expectedWhiteScore: 0, expectedError: "",
		},
		{
			name:   "should cause an error when invalid team ID",
			teamID: -1, expectedBlueScore: 0,
			expectedWhiteScore: 0,
			expectedError:      "bad team ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			game.score.WhiteScore = 0
			game.score.BlueScore = 0
			err := game.AddGoal(tt.teamID)
			if err == nil {
				resultScore := <-game.scoreChannel

				assert.Equal(t, resultScore.BlueScore, tt.expectedBlueScore, "blue team score changes incorrectly")
				assert.Equal(t, resultScore.WhiteScore, tt.expectedWhiteScore, "white team score changes incorrectly")
			}

			if tt.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestGameSubGoal(t *testing.T) {
	game := &game{score: GameScore{1, 1}, scoreChannel: make(chan GameScore, 32)}

	type args struct {
		name               string
		teamID             int
		expectedBlueScore  int
		expectedWhiteScore int
		expectedError      string
	}
	tests := []args{
		{
			name:               "should decrement team white score by one",
			teamID:             config.TeamWhite,
			expectedBlueScore:  1,
			expectedWhiteScore: 0,
			expectedError:      "",
		},
		{
			name:               "should decrement team blue score by one",
			teamID:             config.TeamBlue,
			expectedBlueScore:  0,
			expectedWhiteScore: 1,
			expectedError:      "",
		},
		{
			name:               "should cause an error when invalid team ID",
			teamID:             -1,
			expectedBlueScore:  1,
			expectedWhiteScore: 1,
			expectedError:      "bad team ID",
		},
		{
			name:               "shoud not decrement the score",
			teamID:             config.TeamBlue,
			expectedBlueScore:  0,
			expectedWhiteScore: 0,
			expectedError:      "",
		},
	}

	for id, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game.score.WhiteScore = 1
			game.score.BlueScore = 1

			if id == 3 {
				game.score.WhiteScore = 0
				game.score.BlueScore = 0
			}

			err := game.SubGoal(tt.teamID)
			if err == nil {
				resultScore := <-game.scoreChannel

				assert.Equal(t, resultScore.BlueScore, tt.expectedBlueScore, "blue team score changes incorrectly")
				assert.Equal(t, resultScore.WhiteScore, tt.expectedWhiteScore, "white team score changes incorrectly")
			}
			if tt.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func Test_game_UpdateManualGoals(t *testing.T) {
	game := &game{score: GameScore{0, 0}, manualGoals: ManualGoals{0, 0, 0, 0}}

	type args struct {
		name                string
		teamID              int
		action              string
		expectedManualGoals ManualGoals
		expectedError       string
	}
	tests := []args{
		// No error cases
		{name: "should increment blue manual goals by 1", teamID: config.TeamBlue, action: "add", expectedManualGoals: ManualGoals{1, 0, 0, 0}, expectedError: ""},
		{name: "should decrement blue manual goals by 1", teamID: config.TeamBlue, action: "sub", expectedManualGoals: ManualGoals{0, 1, 0, 0}, expectedError: ""},
		{name: "should increment white manual goals by 1", teamID: config.TeamWhite, action: "add", expectedManualGoals: ManualGoals{0, 0, 1, 0}, expectedError: ""},
		{name: "should decrement white manual goals by 1", teamID: config.TeamWhite, action: "sub", expectedManualGoals: ManualGoals{0, 0, 0, 1}, expectedError: ""},

		// Error cases
		{name: "should return bad team ID error", teamID: 0, action: "add", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad team ID"},
		{name: "should return bad team ID error", teamID: 0, action: "sub", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad team ID"},
		{name: "should return bad action error", teamID: config.TeamBlue, action: "addd", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad action type"},
		{name: "should return bad action error", teamID: config.TeamWhite, action: "addd", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad action type"},
		{name: "should return bad action error", teamID: -10, action: "addd", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad action type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game.manualGoals = ManualGoals{0, 0, 0, 0}
			err := game.UpdateManualGoals(tt.teamID, tt.action)

			if tt.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}

			assert.Equal(t, game.manualGoals.AddedBlue, tt.expectedManualGoals.AddedBlue, "Manual goals for blue team added incorrectly")
			assert.Equal(t, game.manualGoals.AddedWhite, tt.expectedManualGoals.AddedWhite, "Manual goals for white team added incorrectly")
			assert.Equal(t, game.manualGoals.SubtractedBlue, tt.expectedManualGoals.SubtractedBlue, "Manual goals for blue team subtracted incorrectly")
			assert.Equal(t, game.manualGoals.SubtractedWhite, tt.expectedManualGoals.SubtractedWhite, "Manual goals for white team subtracted incorrectly")
		})
	}
}

func TestUpdateShotsData(t *testing.T) {
	game := &game{}

	type args struct {
		name               string
		shot               Shot
		expectedCountWhite int
		expectedCountBlue  int
		expectedError      string
	}
	tests := []args{
		{
			name: "should increment team white shot count by one",
			shot: Shot{
				Speed: 15,
				Team:  1,
			},
			expectedCountWhite: 1,
			expectedCountBlue:  0,
			expectedError:      "",
		},
		{
			name: "should increment team blue shot count by one",
			shot: Shot{
				Speed: 15,
				Team:  2,
			},
			expectedCountWhite: 0,
			expectedCountBlue:  1,
			expectedError:      "",
		},
		{
			name: "should cause an error when invalid team ID",
			shot: Shot{
				Speed: 15,
				Team:  3,
			},
			expectedCountWhite: 0,
			expectedCountBlue:  0,
			expectedError:      "incorrect team ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game.shotsData.WhiteCount = 0
			game.shotsData.BlueCount = 0

			err := game.UpdateShotsData(tt.shot)

			if tt.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}

			assert.Equal(t, tt.expectedCountWhite, game.shotsData.WhiteCount)
			assert.Equal(t, tt.expectedCountBlue, game.shotsData.BlueCount)
		})
	}
}

func TestSaveFastestGoal(t *testing.T) {

	game := &game{
		shotsData: ShotsData{
			Fastest: Shot{
				Speed: 18.98,
				Team:  1,
			},
		},
	}

	type args struct {
		name            string
		shotMsgIn       Shot
		expectedFastest Shot
	}

	tests := []args{
		{
			name: "Should save new fastest of team 1",
			shotMsgIn: Shot{
				Speed: 21.45,
				Team:  1,
			},
			expectedFastest: Shot{
				Speed: 21.45,
				Team:  1,
			},
		},
		{
			name: "Should save new fastest of team 2",
			shotMsgIn: Shot{
				Speed: 55.5555,
				Team:  2,
			},
			expectedFastest: Shot{
				Speed: 55.5555,
				Team:  2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game.saveFastestShot(tt.shotMsgIn)

			assert.Equal(t, tt.expectedFastest, game.shotsData.Fastest)
		})
	}
}
