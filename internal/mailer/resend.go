package mailer

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"

	_ "github.com/joho/godotenv/autoload"
)

var resendKey = os.Getenv("RESEND_API_KEY")

func SendEmailRequest() {
	client := resend.NewClient(resendKey)
	params := &resend.SendEmailRequest{
		From:    "onboarding@mobo-ticketing.tech",
		To:      []string{"antonymbeka@gmail.com"},
		Html:    "<strong>This is just a test</strong>",
		Subject: "Mobo Ticketing Test Email",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(sent)
}
