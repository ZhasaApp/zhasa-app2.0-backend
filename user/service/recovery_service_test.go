package service

type TestRecoveryCodeSender struct {
	codesSent []OtpCode
}

func (t *TestRecoveryCodeSender) sendRecoveryCode(code OtpCode) error {
	t.codesSent = append(t.codesSent, code)
	return nil
}
