package payments

import "testing"

func TestResolveOrderAmounts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		planCode        string
		requestedAmount int64
		isAdmin         bool
		wantPlanCode    string
		wantBaseAmount  int64
		wantAmount      int64
		wantTipAmount   int64
		wantErr         error
	}{
		{
			name:            "public pro ignores lower requested amount",
			planCode:        PlanPro,
			requestedAmount: 1000,
			wantPlanCode:    PlanPro,
			wantBaseAmount:  ProAmountKopecks,
			wantAmount:      ProAmountKopecks,
			wantTipAmount:   0,
		},
		{
			name:            "public donation below minimum price is rejected",
			planCode:        PlanDonation,
			requestedAmount: 1000,
			wantErr:         ErrInvalidAmount,
		},
		{
			name:            "admin can lower pro amount",
			planCode:        PlanPro,
			requestedAmount: 1000,
			isAdmin:         true,
			wantPlanCode:    PlanPro,
			wantBaseAmount:  1000,
			wantAmount:      1000,
			wantTipAmount:   0,
		},
		{
			name:            "admin donation below public minimum is normalized to pro",
			planCode:        PlanDonation,
			requestedAmount: 1000,
			isAdmin:         true,
			wantPlanCode:    PlanPro,
			wantBaseAmount:  1000,
			wantAmount:      1000,
			wantTipAmount:   0,
		},
		{
			name:            "admin donation above base keeps donation split",
			planCode:        PlanDonation,
			requestedAmount: 50000,
			isAdmin:         true,
			wantPlanCode:    PlanDonation,
			wantBaseAmount:  ProAmountKopecks,
			wantAmount:      50000,
			wantTipAmount:   10100,
		},
		{
			name:            "admin override below one ruble is rejected",
			planCode:        PlanPro,
			requestedAmount: 99,
			isAdmin:         true,
			wantErr:         ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotPlanCode, gotBaseAmount, gotAmount, gotTipAmount, err := resolveOrderAmounts(tt.planCode, tt.requestedAmount, tt.isAdmin)
			if err != tt.wantErr {
				t.Fatalf("resolveOrderAmounts() error = %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if gotPlanCode != tt.wantPlanCode {
				t.Fatalf("resolveOrderAmounts() planCode = %q, want %q", gotPlanCode, tt.wantPlanCode)
			}
			if gotBaseAmount != tt.wantBaseAmount {
				t.Fatalf("resolveOrderAmounts() baseAmount = %d, want %d", gotBaseAmount, tt.wantBaseAmount)
			}
			if gotAmount != tt.wantAmount {
				t.Fatalf("resolveOrderAmounts() amount = %d, want %d", gotAmount, tt.wantAmount)
			}
			if gotTipAmount != tt.wantTipAmount {
				t.Fatalf("resolveOrderAmounts() tipAmount = %d, want %d", gotTipAmount, tt.wantTipAmount)
			}
		})
	}
}
