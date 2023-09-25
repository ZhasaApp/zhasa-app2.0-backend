package service

type TestRecoveryCodeSender struct {
	codesSent []int32
}

func (t *TestRecoveryCodeSender) sendRecoveryCode(code int32) error {
	t.codesSent = append(t.codesSent, code)
	return nil
}
