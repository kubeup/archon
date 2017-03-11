package pbmoney

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMath(t *testing.T) {
	money, err := NewMoneySimpleUSD(12, 0).
		Plus(NewMoneySimpleUSD(16, 0)).
		Minus(NewMoneySimpleUSD(8, 0)).
		TimesInt(2).
		DivInt(4).
		TimesFloat(2.0).
		Result()
	require.NoError(t, err)
	require.Equal(t, NewMoneySimpleUSD(20, 0), money)
	money, err = NewMoneySimpleUSD(12, 0).
		Plus(NewMoneySimpleUSD(16, 0)).
		Minus(NewMoneySimpleUSD(8, 0)).
		TimesInt(-2).
		DivInt(4).
		TimesFloat(2.0).
		PlusInt(500000).
		Result()
	require.NoError(t, err)
	require.Equal(t, NewMoneySimpleUSD(-19, 50), money)
}

func TestToFromGoogleMoney(t *testing.T) {
	testToFromGoogleMoney(t, NewMoneySimpleUSD(123456, 78))
	testToFromGoogleMoney(t, NewMoneySimpleUSD(-123456, 78))
}

func testToFromGoogleMoney(t *testing.T, money *Money) {
	money2, err := GoogleMoneyToMoney(money.ToGoogleMoney())
	fmt.Println(money.SimpleString())
	fmt.Println(money2.SimpleString())
	require.NoError(t, err)
	require.Equal(t, money, money2)
}
