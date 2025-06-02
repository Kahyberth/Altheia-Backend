package mail

import "net/smtp"

// Variable to store the send function
var sendMailFunc = smtp.SendMail
