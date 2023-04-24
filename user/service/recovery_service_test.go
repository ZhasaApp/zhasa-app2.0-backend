package service

type TestRecoveryCodeSender struct {
	codesSent []recoveryCode
}

func (t *TestRecoveryCodeSender) sendRecoveryCode(code recoveryCode) error {
	t.codesSent = append(t.codesSent, code)
	return nil
}
