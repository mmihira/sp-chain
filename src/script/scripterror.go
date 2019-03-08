package script

type PubKeyParseError struct {
	Msg string
}

func (p *PubKeyParseError) Error() string {
	return p.Msg
}

type SigParseError struct {
	Msg string
}

func (p *SigParseError) Error() string {
	return p.Msg
}

type SigValidationError struct {
	Msg string
}

func (p *SigValidationError) Error() string {
	return p.Msg
}

type InvalidType struct {
	Msg string
}

func (p *InvalidType) Error() string {
	return p.Msg
}
