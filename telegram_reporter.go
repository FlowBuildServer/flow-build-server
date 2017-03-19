package main

import(
    "github.com/levigross/grequests"
    "fmt"
)

type TelegramReporter struct {
    Telegram *Telegram
}

func (self *TelegramReporter) ReportPending(report *Report) (error) {
    return self.sendReport("pending", report)
}

func (self *TelegramReporter) ReportSuccess(report *Report) (error) {
    return self.sendReport("success", report)
}

func (self *TelegramReporter) ReportError(report *Report) (error) {
    return self.sendReport("error", report)
}

func (self *TelegramReporter) sendReport(status string, report *Report) (error) {
    //message := status
    //if report != nil {
    //    message = report.Message
    //}

    _, err := grequests.Post(
        fmt.Sprintf(
            "https://api.telegram.org/bot%v/sendMessage?chat_id=%v&text=Build status for PR %v: %v",
            self.Telegram.Token,
            self.Telegram.ChatId,
            *report.PullRequest.Number,
            status,
        ),
        nil,
    )

    return err;
}