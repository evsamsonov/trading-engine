package tinkoff

import (
	"context"
	"testing"

	"bou.ke/monkey"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"go.uber.org/zap"

	"github.com/evsamsonov/trengin"
)

const (
	float64EqualityThreshold = 1e-6
)

func TestTinkoff_OpenPosition(t *testing.T) {
	type testWant struct {
		positionType       trengin.PositionType
		orderDirection     investapi.OrderDirection
		stopOrderDirection investapi.StopOrderDirection
		openPrice          *investapi.MoneyValue
		stopLoss           *investapi.Quotation
		takeProfit         *investapi.Quotation
		stopLossID         string
		takeProfitID       string
	}

	tests := []struct {
		name               string
		openPositionAction trengin.OpenPositionAction
		want               testWant
	}{
		{
			name: "long with stop loss and take profit",
			openPositionAction: trengin.OpenPositionAction{
				Type:             trengin.Long,
				StopLossIndent:   11.5,
				TakeProfitIndent: 20.1,
			},
			want: testWant{
				orderDirection:     investapi.OrderDirection_ORDER_DIRECTION_BUY,
				stopOrderDirection: investapi.StopOrderDirection_STOP_ORDER_DIRECTION_SELL,
				positionType:       trengin.Long,
				openPrice: &investapi.MoneyValue{
					Units: 148,
					Nano:  0.2 * 10e8,
				},
				stopLoss: &investapi.Quotation{
					Units: 136,
					Nano:  0.7 * 10e8,
				},
				takeProfit: &investapi.Quotation{
					Units: 168,
					Nano:  0.3 * 10e8,
				},
				stopLossID:   "123",
				takeProfitID: "321",
			},
		},
		{
			name: "short with stop loss and take profit",
			openPositionAction: trengin.OpenPositionAction{
				Type:             trengin.Short,
				StopLossIndent:   11.5,
				TakeProfitIndent: 20.1,
			},
			want: testWant{
				orderDirection:     investapi.OrderDirection_ORDER_DIRECTION_SELL,
				stopOrderDirection: investapi.StopOrderDirection_STOP_ORDER_DIRECTION_BUY,
				positionType:       trengin.Short,
				openPrice: &investapi.MoneyValue{
					Units: 148,
					Nano:  0.2 * 10e8,
				},
				stopLoss: &investapi.Quotation{
					Units: 159,
					Nano:  0.7 * 10e8,
				},
				takeProfit: &investapi.Quotation{
					Units: 128,
					Nano:  0.1 * 10e8,
				},
				stopLossID:   "123",
				takeProfitID: "321",
			},
		},
		{
			name: "without stop loss and take profit",
			openPositionAction: trengin.OpenPositionAction{
				Type:             trengin.Long,
				StopLossIndent:   0.0,
				TakeProfitIndent: 0.0,
			},
			want: testWant{
				orderDirection:     investapi.OrderDirection_ORDER_DIRECTION_BUY,
				stopOrderDirection: investapi.StopOrderDirection_STOP_ORDER_DIRECTION_SELL,
				positionType:       trengin.Long,
				openPrice: &investapi.MoneyValue{
					Units: 148,
					Nano:  0.2 * 10e8,
				},
				stopLoss:     &investapi.Quotation{},
				takeProfit:   &investapi.Quotation{},
				stopLossID:   "",
				takeProfitID: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ordersServiceClient := &mockOrdersServiceClient{}
			stopOrdersServiceClient := &mockStopOrdersServiceClient{}

			monkey.Patch(uuid.New, func() uuid.UUID {
				return uuid.MustParse("8942e9ae-e4e1-11ec-8fea-0242ac120002")
			})

			tinkoff := &Tinkoff{
				accountID:       "123",
				orderClient:     ordersServiceClient,
				stopOrderClient: stopOrdersServiceClient,
				instrumentFIGI:  "FUTSBRF06220",
				tradedQuantity:  2,
				instrument: &investapi.Instrument{
					MinPriceIncrement: &investapi.Quotation{
						Units: 0,
						Nano:  0.1 * 10e8,
					},
				},
				currentPosition: &currentPosition{},
				logger:          zap.NewNop(),
			}

			ordersServiceClient.On("PostOrder", mock.Anything, &investapi.PostOrderRequest{
				Figi:      "FUTSBRF06220",
				Quantity:  2,
				Direction: tt.want.orderDirection,
				AccountId: "123",
				OrderType: investapi.OrderType_ORDER_TYPE_MARKET,
				OrderId:   "8942e9ae-e4e1-11ec-8fea-0242ac120002",
			}).Return(&investapi.PostOrderResponse{
				ExecutionReportStatus: investapi.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_FILL,
				ExecutedOrderPrice:    tt.want.openPrice,
			}, nil)

			if tt.openPositionAction.StopLossIndent != 0 {
				stopOrdersServiceClient.On("PostStopOrder", mock.Anything, &investapi.PostStopOrderRequest{
					Figi:           "FUTSBRF06220",
					Quantity:       2,
					StopPrice:      tt.want.stopLoss,
					Direction:      tt.want.stopOrderDirection,
					AccountId:      "123",
					ExpirationType: investapi.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
					StopOrderType:  investapi.StopOrderType_STOP_ORDER_TYPE_STOP_LOSS,
				}).Return(&investapi.PostStopOrderResponse{
					StopOrderId: "123",
				}, nil).Once()
			}

			if tt.openPositionAction.TakeProfitIndent != 0 {
				stopOrdersServiceClient.On("PostStopOrder", mock.Anything, &investapi.PostStopOrderRequest{
					Figi:           "FUTSBRF06220",
					Quantity:       2,
					StopPrice:      tt.want.takeProfit,
					Direction:      tt.want.stopOrderDirection,
					AccountId:      "123",
					ExpirationType: investapi.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
					StopOrderType:  investapi.StopOrderType_STOP_ORDER_TYPE_TAKE_PROFIT,
				}).Return(&investapi.PostStopOrderResponse{
					StopOrderId: "321",
				}, nil).Once()
			}

			position, _, err := tinkoff.OpenPosition(context.Background(), tt.openPositionAction)
			assert.NoError(t, err)

			assert.Equal(t, tt.want.positionType, position.Type)
			assert.InEpsilon(t, NewMoneyValue(tt.want.openPrice).ToFloat(), position.OpenPrice, float64EqualityThreshold)

			wantStopLoss := NewMoneyValue(tt.want.stopLoss).ToFloat()
			if wantStopLoss != 0 {
				assert.InEpsilon(t, wantStopLoss, position.StopLoss, float64EqualityThreshold)
			} else {
				assert.Equal(t, 0., position.StopLoss)
			}

			wantTakeProfit := NewMoneyValue(tt.want.takeProfit).ToFloat()
			if wantTakeProfit != 0 {
				assert.InEpsilon(t, wantTakeProfit, position.TakeProfit, float64EqualityThreshold)
			} else {
				assert.Equal(t, 0., position.TakeProfit)
			}

			assert.Equal(t, tt.want.stopLossID, tinkoff.currentPosition.StopLossID())
			assert.Equal(t, tt.want.takeProfitID, tinkoff.currentPosition.TakeProfitID())
		})
	}
}

func TestTinkoff_ChangeConditionalOrder_noOpenPosition(t *testing.T) {
	tinkoff := &Tinkoff{
		currentPosition: &currentPosition{},
	}
	_, err := tinkoff.ChangeConditionalOrder(context.Background(), trengin.ChangeConditionalOrderAction{})
	assert.Errorf(t, err, "no open position")
}

func TestTinkoff_ChangeConditionalOrder(t *testing.T) {
	type testWant struct {
		stopLoss           *investapi.Quotation
		takeProfit         *investapi.Quotation
		stopOrderDirection investapi.StopOrderDirection
		stopLossID         string
		takeProfitID       string
	}

	tests := []struct {
		name                       string
		changeConditionOrderAction trengin.ChangeConditionalOrderAction
		positionType               trengin.PositionType
		want                       testWant
	}{
		{
			name: "stop loss and take profit equal zero",
			changeConditionOrderAction: trengin.ChangeConditionalOrderAction{
				PositionID: trengin.PositionID{},
				StopLoss:   0,
				TakeProfit: 0,
			},
		},
		{
			name: "long position, stop loss and take profit set are given",
			changeConditionOrderAction: trengin.ChangeConditionalOrderAction{
				PositionID: trengin.PositionID{},
				StopLoss:   123.43,
				TakeProfit: 156.31,
			},
			positionType: trengin.Long,
			want: testWant{
				stopLoss: &investapi.Quotation{
					Units: 123,
					Nano:  0.43 * 10e8,
				},
				takeProfit: &investapi.Quotation{
					Units: 156,
					Nano:  0.31 * 10e8,
				},
				stopOrderDirection: investapi.StopOrderDirection_STOP_ORDER_DIRECTION_SELL,
				stopLossID:         "2",
				takeProfitID:       "4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ordersServiceClient := &mockOrdersServiceClient{}
			stopOrdersServiceClient := &mockStopOrdersServiceClient{}

			tinkoff := &Tinkoff{
				accountID:       "123",
				orderClient:     ordersServiceClient,
				stopOrderClient: stopOrdersServiceClient,
				instrumentFIGI:  "FUTSBRF06220",
				tradedQuantity:  2,
				instrument: &investapi.Instrument{
					MinPriceIncrement: &investapi.Quotation{
						Units: 0,
						Nano:  0.01 * 10e8,
					},
				},
				currentPosition: &currentPosition{
					position: &trengin.Position{
						Type: tt.positionType,
					},
					stopLossID:   "1",
					takeProfitID: "3",
				},
				logger: zap.NewNop(),
			}

			if tt.changeConditionOrderAction.StopLoss != 0 {
				stopOrdersServiceClient.On("CancelStopOrder", mock.Anything, &investapi.CancelStopOrderRequest{
					AccountId:   "123",
					StopOrderId: "1",
				}).Return(&investapi.CancelStopOrderResponse{}, nil).Once()

				stopOrdersServiceClient.On("PostStopOrder", mock.Anything, &investapi.PostStopOrderRequest{
					Figi:           "FUTSBRF06220",
					Quantity:       2,
					StopPrice:      tt.want.stopLoss,
					Direction:      tt.want.stopOrderDirection,
					AccountId:      "123",
					ExpirationType: investapi.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
					StopOrderType:  investapi.StopOrderType_STOP_ORDER_TYPE_STOP_LOSS,
				}).Return(&investapi.PostStopOrderResponse{
					StopOrderId: "2",
				}, nil).Once()
			}

			if tt.changeConditionOrderAction.TakeProfit != 0 {
				stopOrdersServiceClient.On("CancelStopOrder", mock.Anything, &investapi.CancelStopOrderRequest{
					AccountId:   "123",
					StopOrderId: "3",
				}).Return(&investapi.PostStopOrderResponse{
					StopOrderId: "4",
				}, nil).Once()

				stopOrdersServiceClient.On("PostStopOrder", mock.Anything, &investapi.PostStopOrderRequest{
					Figi:           "FUTSBRF06220",
					Quantity:       2,
					StopPrice:      tt.want.takeProfit,
					Direction:      tt.want.stopOrderDirection,
					AccountId:      "123",
					ExpirationType: investapi.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
					StopOrderType:  investapi.StopOrderType_STOP_ORDER_TYPE_STOP_LOSS,
				}).Return(&investapi.PostStopOrderResponse{
					StopOrderId: "123",
				}, nil).Once()
			}

			position, err := tinkoff.ChangeConditionalOrder(context.Background(), trengin.ChangeConditionalOrderAction{
				PositionID: trengin.PositionID{},
				StopLoss:   tt.changeConditionOrderAction.StopLoss,
				TakeProfit: tt.changeConditionOrderAction.TakeProfit,
			})
			assert.NoError(t, err)
			_ = position
			assert.Equal(t, tt.want.stopLossID, tinkoff.currentPosition.StopLossID())
			assert.Equal(t, tt.want.takeProfitID, tinkoff.currentPosition.TakeProfitID())

		})
	}
}

func TestTinkoff_stopLossPriceByOpen(t *testing.T) {
	tests := []struct {
		name      string
		openPrice *investapi.MoneyValue
		action    trengin.OpenPositionAction
		want      *investapi.Quotation
	}{
		{
			name: "long nano is zero",
			openPrice: &investapi.MoneyValue{
				Units: 123,
				Nano:  0,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Long,
				StopLossIndent: 5,
			},
			want: &investapi.Quotation{
				Units: 118,
				Nano:  0,
			},
		},
		{
			name: "long nano without overflow",
			openPrice: &investapi.MoneyValue{
				Units: 123,
				Nano:  430000000,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Long,
				StopLossIndent: 50.5,
			},
			want: &investapi.Quotation{
				Units: 72,
				Nano:  930000000,
			},
		},
		{
			name: "long nano with overflow",
			openPrice: &investapi.MoneyValue{
				Units: 123,
				Nano:  530000000,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Long,
				StopLossIndent: 50.556,
			},
			want: &investapi.Quotation{
				Units: 72,
				Nano:  974000000,
			},
		},
		{
			name: "short nano is zero",
			openPrice: &investapi.MoneyValue{
				Units: 123,
				Nano:  0,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Short,
				StopLossIndent: 5,
			},
			want: &investapi.Quotation{
				Units: 128,
				Nano:  0,
			},
		},
		{
			name: "short nano without overflow",
			openPrice: &investapi.MoneyValue{
				Units: 123,
				Nano:  430000000,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Short,
				StopLossIndent: 50.4,
			},
			want: &investapi.Quotation{
				Units: 173,
				Nano:  830000000,
			},
		},
		{
			name: "short nano with overflow",
			openPrice: &investapi.MoneyValue{
				Units: 123,
				Nano:  530000000,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Short,
				StopLossIndent: 50.556,
			},
			want: &investapi.Quotation{
				Units: 174,
				Nano:  86000000,
			},
		},
	}

	tinkoff := Tinkoff{
		instrument: &investapi.Instrument{
			MinPriceIncrement: &investapi.Quotation{
				Nano: 1000000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openPrice := NewMoneyValue(tt.openPrice)
			quotation := tinkoff.stopLossPriceByOpen(openPrice, tt.action)
			assert.Equal(t, tt.want, quotation)
		})
	}
}