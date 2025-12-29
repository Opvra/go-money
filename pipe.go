package money

//go:generate go run ./internal/cmd/genpipe -in money.go -out pipe_gen.go

type Pipe struct {
	money Money
	err   error
}

func PipeOf(m Money) Pipe {
	return Pipe{money: m}
}

func (p Pipe) Result() (Money, error) {
	if p.err != nil {
		return Money{}, p.err
	}
	return p.money, nil
}
