package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendOtp(toAddress string, otp string) error {
	signUpMessage := fmt.Sprintf(
		`
		<!DOCTYPE html>
		<html>

		<head>
			<meta charset="UTF-8">
			<title>Your FitNora OTP</title>
		</head>

		<body style="margin:0; padding:0; background:#f5f7fa; font-family:Arial, Helvetica, sans-serif;">
			<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px 0;">
				<tr>
					<td align="center">
						<table width="600" cellpadding="0" cellspacing="0"
							style="background:#ffffff; border-radius:10px; padding:30px; box-shadow:0 4px 14px rgba(0,0,0,0.08);">

							<!-- Header -->
							<tr>
								<td align="center" style="padding-bottom:20px;">
									<h2 style="margin:0; font-size:24px; color:#111827;">FitNora Verification Code</h2>
								</td>
							</tr>

							<!-- Message -->
							<tr>
								<td style="font-size:15px; color:#4b5563; line-height:1.6; padding-bottom:22px;">
									Hello,<br />
									Use the One-Time Password (OTP) below to complete your signup or verification.
									This code is valid for <strong>5 minutes</strong>.
								</td>
							</tr>

							<!-- OTP Box -->
							<tr>
								<td align="center" style="padding:25px 0;">
									<div style="
							font-size:32px;
							letter-spacing:10px;
							font-weight:bold;
							background:#f0f4ff;
							color:#1d4ed8;
							padding:18px 30px;
							border-radius:8px;
							display:inline-block;
						">
										%s
									</div>
								</td>
							</tr>

							<!-- Footer text -->
							<tr>
								<td style="font-size:13px; color:#6b7280; line-height:1.6; padding-top:20px;">
									If you did not request this code, please ignore this email.
									Someone may have entered your email by mistake.
								</td>
							</tr>

							<!-- Bottom spacing -->
							<tr>
								<td style="padding-top:25px; text-align:center; font-size:12px; color:#9ca3af;">
									&copy; FitNora • All rights reserved.
								</td>
							</tr>

						</table>
					</td>
				</tr>
			</table>
		</body>

		</html>
		`, otp)

	return sendEmail(toAddress, signUpMessage)
}

func sendEmail(toAddress string, htmlMessage string) error {
	// Email Credentials
	fromAddress := os.Getenv("SMTP_MAIL")
	appPassword := os.Getenv("SMTP_PASSWORD")

	// Gmail SMTP Config
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// MIME headers -> required for HTML email
	msg := "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "From: FitNora <" + fromAddress + ">\r\n"
	msg += "To: " + toAddress + "\r\n"
	msg += "Subject: Your OTP Code\r\n\r\n"
	msg += htmlMessage

	// Auth
	auth := smtp.PlainAuth("", fromAddress, appPassword, smtpHost)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		fromAddress,
		[]string{toAddress},
		[]byte(msg),
	)

	if err != nil {
		fmt.Println("SMTP Error:", err)
		return err
	}

	log.Println("Email sent successfully to", toAddress)
	return nil
}
