package models

import "testing"

func TestPongPaddle(t *testing.T) {

    tests := []struct {
        expectedX float32
        expectedY float32
    }{
        {20, 20},
        {30, 30},
    }

    paddles := []PongPaddle{
        {20, 20},
        {30, 30},
    }

    for i, tt := range tests {
        if tt.expectedX != paddles[i].X || tt.expectedY != paddles[i].Y {
            t.Fatal("Expected paddle X or Y not equal!")
        }
    }

}

func TestPongPaddleDimensions(t *testing.T) {
    paddles := []PongPaddle{
        {20, 20},
        {30, 30},
    }

    for _, p := range paddles {
        if p.Height() != 70 {
            t.Fatalf("p.Height() not equal to 70. Got: %f", p.Height())
        }

        if p.Width() != 18 {
            t.Fatalf("p.Width() not equal to 18. Got: %f", p.Width())
        }
    }
}
