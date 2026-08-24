package game

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewGame(t *testing.T) {
	testID := uuid.New()
	nilID := uuid.Nil
	tests := []struct {
		name string
		req  GameCreation
		want Game
	}{
		{
			name: "PvP",
			req: GameCreation{
				PlayerX: nilID,
				PlayerO: testID,
				Mode:    PvP,
			},
			want: Game{
				Field:   Field{},
				PlayerX: nilID,
				PlayerO: testID,
				Status:  WaitingForOpponent,
				Mode:    PvP,
				Version: 1,
			},
		}, {
			name: "PvE",
			req: GameCreation{
				PlayerX: nilID,
				PlayerO: testID,
				Mode:    PvE,
			},
			want: Game{
				Field:   Field{},
				PlayerX: nilID,
				PlayerO: testID,
				Status:  PlayerXTurn,
				Mode:    PvE,
				Version: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := NewGame(tt.req)

			if got.GameID == uuid.Nil {
				t.Error("GameID is nil UUID, want non-nil")
			}

			if got.Field != tt.want.Field {
				t.Errorf("Field = %v, want %v", got.Field, tt.want.Field)
			}

			if got.PlayerX != tt.want.PlayerX {
				t.Errorf("PlayerX = %q, want %q", got.PlayerX, tt.want.PlayerX)
			}

			if got.PlayerO != tt.want.PlayerO {
				t.Errorf("PlayerO = %q, want %q", got.PlayerO, tt.want.PlayerO)
			}

			if got.Status != tt.want.Status {
				t.Errorf("Status = %q, want %q", got.Status, tt.want.Status)
			}

			if got.Mode != tt.want.Mode {
				t.Errorf("Mode = %q, want %q", got.Mode, tt.want.Mode)
			}

			if got.Version != tt.want.Version {
				t.Errorf("Version = %d, want %d", got.Version, tt.want.Version)
			}
		})
	}
}
