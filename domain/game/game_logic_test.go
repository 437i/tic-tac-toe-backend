package game

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestFieldFull(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  bool
	}{
		{
			name:  "empty field",
			field: Field{},
			want:  false,
		}, {
			name: "1 cell filled",
			field: Field{
				{Empty, Empty, Empty},
				{Empty, X, Empty},
				{Empty, Empty, Empty},
			},
			want: false,
		}, {
			name: "1 cell empty",
			field: Field{
				{X, O, X},
				{O, Empty, O},
				{O, X, O},
			},
			want: false,
		}, {
			name: "field full with X",
			field: Field{
				{X, X, X},
				{X, X, X},
				{X, X, X},
			},
			want: true,
		}, {
			name: "field full with mixed X/O",
			field: Field{
				{X, O, X},
				{O, X, O},
				{O, X, O},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fieldFull(tt.field); got != tt.want {
				t.Errorf("fieldFull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWinner(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  XO
	}{
		{
			name: "no winner",
			field: Field{
				{X, O, X},
				{X, O, X},
				{O, X, O},
			},
			want: Empty,
		}, {
			name: "X wins top row",
			field: Field{
				{X, X, X},
				{O, O, X},
				{O, X, O},
			},
			want: X,
		}, {
			name: "X wins left col",
			field: Field{
				{X, O, X},
				{X, O, O},
				{X, X, O},
			},
			want: X,
		}, {
			name: "X wins bottom row",
			field: Field{
				{O, X, O},
				{X, O, O},
				{X, X, X},
			},
			want: X,
		}, {
			name: "X wins main diagonal",
			field: Field{
				{X, O, X},
				{O, X, O},
				{O, X, X},
			},
			want: X,
		}, {
			name: "O wins mid row",
			field: Field{
				{X, O, X},
				{O, O, O},
				{X, X, Empty},
			},
			want: O,
		}, {
			name: "O wins mid col",
			field: Field{
				{X, O, X},
				{O, O, X},
				{X, O, Empty},
			},
			want: O,
		}, {
			name: "O wins right col",
			field: Field{
				{X, X, O},
				{O, X, O},
				{X, O, O},
			},
			want: O,
		}, {
			name: "O wins secondary diagonal",
			field: Field{
				{X, O, O},
				{X, O, X},
				{O, X, X},
			},
			want: O,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWinner(tt.field); got != tt.want {
				t.Errorf("getWinner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateCell(t *testing.T) {
	type input struct {
		old, new, player XO
		counter          int
	}
	type want struct {
		counter int
		err     error
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "same cells, right player",
			input: input{
				old:     X,
				new:     X,
				player:  X,
				counter: 0,
			},
			want: want{
				counter: 0,
			},
		}, {
			name: "same cells, wrong player",
			input: input{
				old:     X,
				new:     X,
				player:  O,
				counter: 5,
			},
			want: want{
				counter: 5,
			},
		}, {
			name: "Empty -> X, right player",
			input: input{
				old:     Empty,
				new:     X,
				player:  X,
				counter: 0,
			},
			want: want{
				counter: 1,
			},
		}, {
			name: "Empty -> X, wrong player",
			input: input{
				old:     Empty,
				new:     X,
				player:  O,
				counter: 0,
			},
			want: want{
				counter: 0,
				err:     ErrWrongSign,
			},
		}, {
			name: "X -> O, right player",
			input: input{
				old:     X,
				new:     O,
				player:  O,
				counter: 5,
			},
			want: want{
				counter: 0,
				err:     ErrCheating,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter, err := validateCell(
				tt.input.old,
				tt.input.new,
				tt.input.player,
				tt.input.counter,
			)

			if !errors.Is(err, tt.want.err) {
				t.Errorf("err = %v, want %v", err, tt.want.err)
			}

			if counter != tt.want.counter {
				t.Errorf("counter = %d, want %d", counter, tt.want.counter)
			}
		})
	}
}

func TestCompareFieldWith(t *testing.T) {
	tests := []struct {
		name    string
		newGame Game
		oldGame Game
		player  XO
		wantErr error
	}{
		{
			name: "no changes",
			newGame: Game{
				Field: Field{
					{X, Empty, X},
					{O, Empty, O},
					{X, Empty, X},
				},
			},
			oldGame: Game{
				Field: Field{
					{X, Empty, X},
					{O, Empty, O},
					{X, Empty, X},
				},
			},
			player:  X,
			wantErr: ErrNoChanges,
		}, {
			name: "valid change",
			newGame: Game{
				Field: Field{
					{X, X, X},
					{O, Empty, O},
					{X, Empty, X},
				},
			},
			oldGame: Game{
				Field: Field{
					{X, Empty, X},
					{O, Empty, O},
					{X, Empty, X},
				},
			},
			player:  X,
			wantErr: nil,
		}, {
			name: "2 changes",
			newGame: Game{
				Field: Field{
					{X, X, X},
					{O, X, O},
					{X, Empty, X},
				},
			},
			oldGame: Game{
				Field: Field{
					{X, Empty, X},
					{O, Empty, O},
					{X, Empty, X},
				},
			},
			player:  X,
			wantErr: ErrManyChanges,
		}, {
			name: "cheating",
			newGame: Game{
				Field: Field{
					{O, Empty, X},
					{O, Empty, O},
					{X, Empty, X},
				},
			},
			oldGame: Game{
				Field: Field{
					{X, Empty, X},
					{O, Empty, O},
					{X, Empty, X},
				},
			},
			player:  X,
			wantErr: ErrCheating,
		}, {
			name: "wrong sign",
			newGame: Game{
				Field: Field{
					{X, O, X},
					{O, Empty, O},
					{X, Empty, X},
				},
			},
			oldGame: Game{
				Field: Field{
					{X, Empty, X},
					{O, Empty, O},
					{X, Empty, X},
				},
			},
			player:  X,
			wantErr: ErrWrongSign,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.newGame.compareFieldWith(tt.oldGame, tt.player); !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJoin(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	tests := []struct {
		name    string
		game    Game
		userID  uuid.UUID
		wantErr error
	}{
		{
			name: "valid join",
			game: Game{
				PlayerX: id1,
				PlayerO: uuid.Nil,
				Status:  WaitingForOpponent,
			},
			userID:  id2,
			wantErr: nil,
		}, {
			name: "player to join is PlayerX",
			game: Game{
				PlayerX: id1,
				PlayerO: uuid.Nil,
				Status:  WaitingForOpponent,
			},
			userID:  id1,
			wantErr: ErrGameUnavailableToJoin,
		}, {
			name: "player to join is PlayerO",
			game: Game{
				PlayerX: uuid.Nil,
				PlayerO: id1,
				Status:  WaitingForOpponent,
			},
			userID:  id1,
			wantErr: ErrGameUnavailableToJoin,
		}, {
			name: "game full",
			game: Game{
				PlayerX: uuid.Nil,
				PlayerO: id1,
				Status:  PlayerXTurn,
			},
			userID:  id2,
			wantErr: ErrGameUnavailableToJoin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.game.validateJoin(tt.userID); !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
