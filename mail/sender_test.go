package mail

import (
	"testing"

	"github.com/Samudra-G/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if(testing.Short()) {
		t.Skip()
	}
	
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test Email"
	content := `
	<h1>Hello, test email</h1>
	<p>This is a test email sent from the Simple Bank application by <a href="https://samudra-portfolio.vercel.app/"> Samudra G</p>
	`
	to := []string{"samudramukhar@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(
		subject,
		content,
		to,
		nil,
		nil,
		attachFiles,
	)
	require.NoError(t, err, "Failed to send email")
}	